package ports

type Logger interface {
	Log(text string)
	Exception(msg string, err error)
}
