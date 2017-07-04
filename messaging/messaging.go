package messaging

type Message struct {
	Filename  string
	Filenames []string
	Contents  string
}

type Action func() Result

type Result struct {
	Message
	Error error
}

type Payload struct {
	Message Message
	Action  Action
	Done    chan Result
}
