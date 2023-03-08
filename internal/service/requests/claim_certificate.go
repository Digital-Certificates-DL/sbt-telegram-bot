package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"regexp"
)

type ClaimCertificateDate struct {
	Address    string `json:"address"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	CourseName string `json:"course"`
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
	ethAddressRegex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return validation.Errors{
		"/address":  validation.Validate(&r.Address, validation.Required, validation.Match(ethAddressRegex)),
		"/name":     validation.Validate(&r.Name, validation.Required),
		"/date":     validation.Validate(&r.Date, validation.Required),
		"/telegram": validation.Validate(&r.Telegram, validation.Required),
	}.Filter()
}
