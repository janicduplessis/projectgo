package usecases

type Logger interface {
	Log(message string)
	Error(err error)
}
