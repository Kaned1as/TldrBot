package api

type InlineQuery struct {
    // Unique identifier for this query
    Id string 
    // Sender
    From User 
    // Optional. Sender location, only for bots that request user location
    Location Location 
    // Text of the query (up to 512 characters)
    Query string 
    // Offset of the results to be returned, can be controlled by the bot
    Offset string 
}

type ChosenInlineResult struct {
    // The unique identifier for the result that was chosen
    Result_id string 
    // The user that chose the result
    From User 
    // Optional. Sender location, only for bots that require user location
    Location *Location
    // Optional. Identifier of the sent inline message.
    // Available only if there is an inline keyboard attached to the message.
    // Will be also received in callback queries and can be used to edit the message.
    Inline_message_id string 
    // The query that was used to obtain the result
    Query string 
}

// TODO: add other objects