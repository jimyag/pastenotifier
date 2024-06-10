package internal

type Handler interface {
	Handle(content string) (title string, message string, err error)
}
