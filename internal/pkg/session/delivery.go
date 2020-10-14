package session

type Delivery interface {
	Create(userID string) (string, error)
}
