package helpers

import (
	"context"
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/bot"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	botCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxBot(entry *bot.Bot) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, botCtxKey, entry)
	}
}

func Bot(r *http.Request) *bot.Bot {
	return r.Context().Value(botCtxKey).(*bot.Bot)
}
