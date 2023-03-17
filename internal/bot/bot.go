package bot

import (
	"context"
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"strings"
	"sync"
)

const timeoutInSec = 60

type Bot struct {
	Info     UserInfo
	token    string
	Admins   []int64
	Bot      *tgbot.BotAPI
	Logger   *logan.Entry
	Template string
}

type UserInfo struct {
	Name     string
	Address  string
	Course   string
	Date     string
	Telegram string
}

func NewBotInit(token string, logger *logan.Entry, template string) (*Bot, error) {
	botAPI, err := tgbot.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	return &Bot{
		token:    token,
		Bot:      botAPI,
		Logger:   logger,
		Template: template,
	}, nil

}

func (b *Bot) NewInfo(name, date, address, course, tg string) UserInfo {
	userInfo := UserInfo{
		Date:     date,
		Name:     name,
		Address:  address,
		Course:   course,
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
			if strings.Contains(err.Error(), "bot was blocked by the user") {
				isAdmin := b.checkAdmin(b.Admins, chatID)
				if isAdmin {
					b.Admins = b.deleteAdmin(b.Admins, chatID)
					b.Logger.Debug("Admin was deleted: ", chatID)
				}
				return nil
			}
			return errors.Wrap(err, fmt.Sprint("failed to send date to: ", chatID))
		}
		b.Logger.Debug("Admins", b.Admins)
	}
	return nil
}

func (b *Bot) PrepareLastMessage() string {
	return fmt.Sprintf(b.Template, b.Info.Name, b.Info.Date, b.Info.Course, b.Info.Address, b.Info.Telegram)
}

func (b *Bot) Start(wg *sync.WaitGroup) error {
	defer wg.Done()
	u := tgbot.NewUpdate(0)
	u.Timeout = timeoutInSec
	ctx, cancel := context.WithCancel(context.Background())
	updates := b.Bot.GetUpdatesChan(u)
	err := b.receiveUpdates(ctx, updates)
	if err != nil {
		return errors.Wrap(err, "failed to receive updates")
	}
	cancel()
	return nil
}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbot.UpdatesChannel) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case update, ok := <-updates:
			if !ok {
				ctx.Done()
				return errors.New("failed to read from chan")
			}
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) handleUpdate(update tgbot.Update) {
	switch {
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
	b.Logger.Info(fmt.Sprintf("%s wrote %s", user.FirstName, text))
	if strings.HasPrefix(text, "/") {
		b.handleCommand(message.Chat.ID, text)
	}
	return
}

func (b *Bot) handleCommand(chatID int64, command string) {
	switch command {
	case "/admin":
		isAdmin := b.checkAdmin(b.Admins, chatID)
		if !isAdmin {
			b.Admins = append(b.Admins, chatID)
		}
		b.Logger.Debug(b.Admins)
		break
	case "/exit":
		isAdmin := b.checkAdmin(b.Admins, chatID)
		if isAdmin {
			b.Admins = b.deleteAdmin(b.Admins, chatID)
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

func (b *Bot) deleteAdmin(sl []int64, name int64) []int64 {
	res := make([]int64, 0)
	for _, value := range sl {
		if value != name {
			res = append(res, value)
		}
	}
	return res
}
