package database

type Contact struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Message struct {
	Sender    uint8  `json:"sender"` // 0 = self, 1 = interlocutor
	Content   string `json:"content"`
	Timestamp int64
}
