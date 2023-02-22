package service

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/sbt-svc/internal/config"
)

type service struct {
	log *logan.Entry
}

func (s *service) run(cfg config.Config) error {
	s.log.Info("Service started")
	err := Start(cfg)
	if err != nil {
		cfg.Log().Error(err)
	}

	return nil
}

func newService(cfg config.Config) *service {
	return &service{
		log: cfg.Log(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}
}
