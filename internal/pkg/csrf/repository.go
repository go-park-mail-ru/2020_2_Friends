package csrf

import "time"

type Repository interface {
	Add(token string, session string, expires time.Duration) error
	Get(token string) (string, error)
}
