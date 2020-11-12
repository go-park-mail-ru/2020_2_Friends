package session

import "github.com/friends/internal/pkg/models"

type Repository interface {
	Create(session models.Session) error
	Check(sessionName string) (userID string, err error)
	Delete(sessionName string) error
}
