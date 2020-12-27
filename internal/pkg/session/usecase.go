package session

//go:generate mockgen -destination=./usecase_mock.go -package=session github.com/friends/internal/pkg/session Usecase
type Usecase interface {
	Create(userID string) (string, error)
	Check(sessionName string) (userID string, err error)
	Delete(sessionName string) error
	SetCSRFToken(userID string) (token string, err error)
	GetTokenFromUser(userID string) (token string, err error)
}
