package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/service/helpers"
	"gitlab.com/tokend/course-certificates/sbt-bot/internal/service/requests"
	"net/http"
)

func ClaimCertificate(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewClaimCertificateRequest(r)
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to parse date:"))
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	bot := helpers.Bot(r)
	bot.Info = bot.NewInfo(req.Name, req.Date, req.Address, req.Telegram)
	helpers.Log(r).Info(bot.Info)
	err = bot.SendToAdmin()
	if err != nil {
		helpers.Log(r).Error(errors.Wrap(err, "failed to send messages:"))
		ape.RenderErr(w, problems.InternalError())
		return
	}
	w.WriteHeader(http.StatusOK)
}
