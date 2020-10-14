package session

type Usecase interface {
	Create(userID string) (string, error)
}
