package messaging

type Message struct {
	Filename  string
	Filenames []string
	Contents  string
}

type Action func() ChanResult

type ChanResult struct {
	Message
	Error error
}
