package service

import (
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/bot"
	"net"
	"net/http"
	"sync"

	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/config"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	cfg      config.Config
	bot      *bot.Bot
}

func (s *service) run() error {
	s.log.Info("Service started")

	botAPI, err := bot.NewBotInit(s.cfg.BotConfig().Token, s.log)
	if err != nil {
		return errors.Wrap(err, "failed to init bot")
	}
	s.bot = botAPI
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	errHTTPChan := make(chan error)
	wg := new(sync.WaitGroup)
	go StartServer(s.listener, r, errHTTPChan, wg)
	go s.bot.Start(wg)
	for {
		err := <-errHTTPChan
		switch err {
		case nil:
			continue
		default:
			return err
		}
	}
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		cfg:      cfg,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}

func StartServer(l net.Listener, handler http.Handler, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	srv := &http.Server{Handler: handler}
	errChan <- srv.Serve(l)
}
