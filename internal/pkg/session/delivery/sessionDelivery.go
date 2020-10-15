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

func (sd SessionDelivery) Check(sessionName string) (userID string, err error) {
	userID, err = sd.sessionUsecase.Check(sessionName)

	return userID, err
}

func (sd SessionDelivery) Delete(sessionName string) error {
	return sd.sessionUsecase.Delete(sessionName)
}
