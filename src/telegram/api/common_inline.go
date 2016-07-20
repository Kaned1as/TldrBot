package api

type InlineQuery struct {
    // Unique identifier for this query
    id string 
    // Sender
    from User 
    // Optional. Sender location, only for bots that request user location
    location Location 
    // Text of the query (up to 512 characters)
    query string 
    // Offset of the results to be returned, can be controlled by the bot
    offset string 
}

type ChosenInlineResult struct {
    // The unique identifier for the result that was chosen
    result_id string 
    // The user that chose the result
    from User 
    // Optional. Sender location, only for bots that require user location
    location Location 
    // Optional. Identifier of the sent inline message.
    // Available only if there is an inline keyboard attached to the message.
    // Will be also received in callback queries and can be used to edit the message.
    inline_message_id string 
    // The query that was used to obtain the result
    query string 
}

// TODO: add other objects