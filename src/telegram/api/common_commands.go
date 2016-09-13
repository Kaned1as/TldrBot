package api

type SendMessage struct {
    ChatId                  int         `json:"chat_id"`
    Text                    string      `json:"text"`
    ParseMode               string      `json:"parse_mode,omitempty"`
    DisableWebPagePreview   *bool       `json:"disable_web_page_preview"`
    DisableNotification     *bool       `json:"disable_notification"`
    ReplyToMessageId        *int        `json:"reply_to_message_id"`
    //ReplyMarkup *
}
