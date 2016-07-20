package api

type ReplyKeyboardMarkup struct {
    // Array of button rows, each represented by an Array of KeyboardButton objects
    keyboard [][]KeyboardButton 
    // Optional. Requests clients to resize the keyboard vertically for optimal fit
    // (e.g., make the keyboard smaller if there are just two rows of buttons).
    // Defaults to false, in which case the custom keyboard is always of the same height as the app's standard keyboard.
    resize_keyboard bool 
    // Optional. Requests clients to hide the keyboard as soon as it's been used.
    // The keyboard will still be available, but clients will automatically display the usual
    // letter-keyboard in the chat – the user can press a special button in the input field
    // to see the custom keyboard again. Defaults to false.
    one_time_keyboard bool 
    // Optional. Use this parameter if you want to show the keyboard to specific users only.
    // Targets:
    //      1) users that are @mentioned in the text of the Message object;
    //      2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
    selective bool 
}

type KeyboardButton struct {
    // Text of the button.
    // If none of the optional fields are used, it will be sent to the bot as a message when the button is pressed
    text string 
    // Optional. If True, the user's phone number will be sent as a contact when the button is pressed.
    // Available in private chats only
    request_contact bool 
    // Optional. If True, the user's current location will be sent when the button is pressed.
    // Available in private chats only
    request_location bool 
}

type ReplyKeyboardHide struct {
    // Requests clients to hide the custom keyboard
    hide_keyboard bool 
    // Optional. Use this parameter if you want to hide keyboard for specific users only.
    // Targets:
    //      1) users that are @mentioned in the text of the Message object;
    //      2) if the bot's message is a reply (has reply_to_message_id), sender of the original message
    selective bool 
}

type InlineKeyboardMarkup struct {
    // Array of button rows, each represented by an Array of InlineKeyboardButton objects
    inline_keyboard [][]InlineKeyboardButton 
}

type InlineKeyboardButton struct {
    // Label text on the button
    text string 
    // Optional. HTTP url to be opened when button is pressed
    url string 
    // Optional. Data to be sent in a callback query to the bot when button is pressed, 1-64 bytes
    callback_data string 
    // Optional. If set, pressing the button will prompt the user to select one of their chats,
    // open that chat and insert the bot‘s username and the specified inline query in the input field.
    // Can be empty, in which case just the bot’s username will be inserted.
    switch_inline_query string 
}