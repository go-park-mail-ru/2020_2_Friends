package session

type Usecase interface {
	Create(userID string) (string, error)
	Check(sessionName string) (userID string, err error)
	Delete(sessionName string) error
}
