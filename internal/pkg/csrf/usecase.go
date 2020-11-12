package csrf

type Usecase interface {
	Add(session string) (string, error)
	Check(token string, session string) bool
}
