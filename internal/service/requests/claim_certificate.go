package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type ClaimCertificateDate struct {
	Address    string `json:"address"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	CourseName string `json:"course_name"`
	Telegram   string `json:"telegram"`
}

func NewClaimCertificateRequest(r *http.Request) (ClaimCertificateDate, error) {
	request := ClaimCertificateDate{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return ClaimCertificateDate{}, errors.Wrap(err, "failed to decode request")
	}
	return request, request.validation()
}

func (r ClaimCertificateDate) validation() error {
	return validation.Errors{
		"/address":  validation.Validate(&r.Address, validation.Required),
		"/name":     validation.Validate(&r.Name, validation.Required),
		"/telegram": validation.Validate(&r.Telegram, validation.Required),
	}.Filter()
}
