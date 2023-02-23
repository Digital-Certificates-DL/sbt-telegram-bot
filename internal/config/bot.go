package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type BotConfiger interface {
	BotConfig() *BotConfig
}

type BotConfig struct {
	Token string `fig:"token"`
}

func NewBotConfiger(getter kv.Getter) BotConfiger {
	return &tokenConfig{
		getter: getter,
	}
}

type tokenConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *tokenConfig) BotConfig() *BotConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "bot")
		config := BotConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &config
	}).(*BotConfig)
}
