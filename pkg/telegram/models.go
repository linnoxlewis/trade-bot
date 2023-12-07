package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type UpdatesReplyResponse struct {
	Ok     bool          `json:"ok"`
	Result []ReplyUpdate `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type ReplyUpdate struct {
	ID      int            `json:"update_id"`
	Message *CallbackQuery `json:"callback_query"`
}

type CallbackQuery struct {
	ID      string `json:"id"`
	From    From   `json:"from"`
	Message struct {
		MessageID int    `json:"message_id"`
		Chat      Chat   `json:"chat"`
		Date      int    `json:"date"`
		Text      string `json:"text"`
	} `json:"message"`
	ChatInstance string `json:"chat_instance"`
	Data         string `json:"data"`
}

type IncomingMessage struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
	Language string `json:"language_code"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
