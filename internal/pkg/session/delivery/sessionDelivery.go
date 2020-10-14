package delivery

import "github.com/friends/internal/pkg/session"

type SessionDelivery struct {
	sessionUsecase session.Usecase
}

func NewSessionDelivery(usecase session.Usecase) SessionDelivery {
	return SessionDelivery{
		sessionUsecase: usecase,
	}
}

func (sd SessionDelivery) Create(userID string) (string, error) {
	return sd.sessionUsecase.Create(userID)
}
