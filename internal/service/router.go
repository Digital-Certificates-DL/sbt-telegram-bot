package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/service/handlers"
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/service/helpers"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxBot(s.bot),
		),
	)
	r.Route("/integrations/sbt-bot", func(r chi.Router) {
		r.Post("/", handlers.ClaimCertificate)
	})
	return r
}
