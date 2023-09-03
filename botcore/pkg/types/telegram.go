package types

type SendMessage struct {
	Method string `json:"method"`
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}
