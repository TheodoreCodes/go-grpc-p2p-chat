package events

type EmptyMessage struct{}

type SendMessageEvent struct {
	Message string
}

type SentMessageEvent struct {
	Message string
}

type ReceiveMessageEvent struct {
	Message string
}

type DisplayModalEvent struct {
	Message string
}
