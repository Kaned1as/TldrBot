package api

type GetUpdatesResponse struct {
    Ok bool
    result []Update
}

type Update struct {
    update_id int
    message Message
    edited_message Message
    inline_query InlineQuery
    chosen_inline_result ChosenInlineResult
    callback_query CallbackQuery
}

type Message struct {

}

type InlineQuery struct {

}

type ChosenInlineResult struct {

}

type CallbackQuery struct {

}
