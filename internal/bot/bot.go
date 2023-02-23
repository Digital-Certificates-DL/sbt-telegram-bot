package bot

import (
	"bufio"
	"context"
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"log"
	"os"
	"strings"
	"sync"
)

const template = "student: %s\nexam date: %s\ncourse: %s\naddress: %s\ntelegram: %s"

type Bot struct {
	Info   UserInfo
	Mutex  *sync.Mutex
	token  string
	Admins []int64
	Bot    *tgbot.BotAPI
	Logger *logan.Entry
}

type UserInfo struct {
	Name     string
	Address  string
	Course   string
	Date     string
	Telegram string
}

func NewBotInit(token string, logger *logan.Entry) (*Bot, error) {
	botAPI, err := tgbot.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	return &Bot{
		Mutex:  new(sync.Mutex),
		token:  token,
		Bot:    botAPI,
		Logger: logger,
	}, nil

}

func (b *Bot) NewInfo(name, date, address, tg string) UserInfo {
	userInfo := UserInfo{
		Date:     date,
		Name:     name,
		Address:  address,
		Telegram: tg,
	}
	return userInfo
}

func (b *Bot) SendToAdmin() error {
	for _, chatID := range b.Admins {
		msg := tgbot.NewMessage(chatID, b.PrepareLastMessage())
		msg.ParseMode = tgbot.ModeMarkdownV2
		_, err := b.Bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, fmt.Sprint("failed to send date to: ", chatID))
		}
	}
	return nil
}

func (b *Bot) PrepareLastMessage() string {
	return fmt.Sprintf(template, b.Info.Name, b.Info.Date, b.Info.Course, b.Info.Address, b.Info.Telegram)
}

func (b *Bot) Start(wg *sync.WaitGroup) error {
	defer wg.Done()
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	updates := b.Bot.GetUpdatesChan(u)
	go b.receiveUpdates(ctx, updates)
	b.Logger.Info("Start listening for updates. Press enter to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

	return nil
}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbot.UpdatesChannel) {

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) handleUpdate(update tgbot.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		b.handleMessage(update.Message)
		break

	}
}

func (b *Bot) handleMessage(message *tgbot.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("%s wrote %s", user.FirstName, text)

	if strings.HasPrefix(text, "/") {
		b.handleCommand(message.Chat.ID, text)

	}
	return
}

func (b *Bot) handleCommand(chatId int64, command string) {

	switch command {
	case "/start":
		isAdmin := b.checkAdmin(b.Admins, chatId)
		if !isAdmin {
			b.Admins = append(b.Admins, chatId)
		}
		b.Logger.Debug(b.Admins)
		break
	}
}

func (b *Bot) checkAdmin(sl []int64, name int64) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}
