package api

// Request for updates for this bot
// This is main API method for polling
// Note that update fields are case-sensitive so denote them correctly
type GetUpdatesRequest struct {
    // last update_id known to server
    Offset int64    `json:"offset"`
    // optional, from 1 to 100
    Limit int       `json:"limit,omitempty"`
    // timeout for polling (usually 0)
    Timeout int     `json:"timeout,omitempty"`
}

type GetUpdatesResponse struct {
    Ok bool
    Result []Update
}

type Update struct {
    // The update‘s unique identifier. Update identifiers start from a certain positive number and increase sequentially.
    // This ID becomes especially handy if you’re using Webhooks,
    // since it allows you to ignore repeated updates or to restore the correct update sequence,
    // should they get out of order.
    Update_id int64
    // Optional. New incoming message of any kind — text, photo, sticker, etc.
    Message *Message
    // Optional. New version of a message that is known to the bot and was edited
    Edited_message *Message
    // Optional. New incoming inline query
    Inline_query *InlineQuery
    // Optional. The result of an inline query that was chosen by a user and sent to their chat partner.
    Chosen_inline_result *ChosenInlineResult
    // Optional. New incoming callback query
    Callback_query *CallbackQuery
}

type CallbackQuery struct {
    // Unique identifier for this query
    Id string 
    // Sender
    From User 
    // Optional. Message with the callback button that originated the query.
    // Note that message content and message date will not be available if the message is too old
    Message Message 
    // Optional. Identifier of the message sent via the bot in inline mode, that originated the query
    Inline_message_id string 
    // Data associated with the callback button. Be aware that a bad client can send arbitrary data in this field
    Data string 
}

type ForceReply struct {
    // Shows reply interface to the user, as if they manually selected the bot‘s message and tapped ’Reply'
    Force_reply bool 
    // Optional. Use this parameter if you want to force reply from specific users only.
    // Targets:
    //      1) users that are @mentioned in the text of the Message object;
    //      2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
    Selective bool 
}