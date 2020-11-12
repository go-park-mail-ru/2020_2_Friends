package delivery

import (
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/csrf"
)

type CSRFDelivery struct {
	csrfUsecase csrf.Usecase
}

func New(csrfUsecase csrf.Usecase) CSRFDelivery {
	return CSRFDelivery{
		csrfUsecase: csrfUsecase,
	}
}

func (c CSRFDelivery) SetCSRF(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(configs.SessionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session := cookie.Value

	token, err := c.csrfUsecase.Add(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
	w.Header().Set("X-CSRF-Token", token)
}
