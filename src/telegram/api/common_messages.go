package api

type Message struct {
    // Unique message identifier
    Message_id int 
    // Optional. Sender, can be empty for messages sent to channels
    From User 
    // Date the message was sent in Unix time
    Date int 
    // Conversation the message belongs to
    Chat Chat 
    // Optional. For forwarded messages, sender of the original message
    Forward_from User 
    // Optional. For messages forwarded from a channel, information about the original channel
    Forward_from_chat Chat 
    // Optional. For forwarded messages, date the original message was sent in Unix time
    Forward_date int 
    // Optional. For replies, the original message. Note that the Message object in this field
    // will not contain further reply_to_message fields even if it itself is a reply.
    //Reply_to_message Message
    // Optional. Date the message was last edited in Unix time
    Edit_date int 
    // Optional. For text messages, the actual UTF-8 text of the message, 0-4096 characters.
    Text string 
    // Optional. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text
    Entities []MessageEntity 
    // Optional. Message is an audio file, information about the file
    Audio Audio 
    // Optional. Message is a general file, information about the file
    Document Document 
    // Optional. Message is a photo, available sizes of the photo
    Photo []PhotoSize 
    // Optional. Message is a sticker, information about the sticker
    Sticker Sticker 
    // Optional. Message is a video, information about the video
    Video Video 
    // Optional. Message is a voice message, information about the file
    Voice Voice 
    // Optional. Caption for the document, photo or video, 0-200 characters
    Caption string 
    // Optional. Message is a shared contact, information about the contact
    Contact Contact 
    // Optional. Message is a shared location, information about the location
    Location Location 
    // Optional. Message is a venue, information about the venue
    Venue Venue 
    // Optional. A new member was added to the group, information about them (this member may be the bot itself)
    New_chat_member User 
    // Optional. A member was removed from the group, information about them (this member may be the bot itself)
    Left_chat_member User 
    // Optional. A chat title was changed to this value
    New_chat_title string 
    // Optional. A chat photo was change to this value
    New_chat_photo []PhotoSize 
    // Optional. Service message: the chat photo was deleted
    Delete_chat_photo bool 
    // Optional. Service message: the group has been created
    Group_chat_created bool 
    // Optional. Service message: the supergroup has been created.
    // This field can‘t be received in a message coming through updates,
    // because bot can’t be a member of a supergroup when it is created.
    // It can only be found in reply_to_message if someone replies to a very first message
    // in a directly created supergroup.
    Supergroup_chat_created bool 
    // Optional. Service message: the channel has been created.
    // This field can‘t be received in a message coming through updates, because bot can’t be a member of a channel
    // when it is created. It can only be found in reply_to_message
    // if someone replies to a very first message in a channel.
    Channel_chat_created bool 
    // Optional. The group has been migrated to a supergroup with the specified identifier.
    // This number may be greater than 32 bits and some programming languages may have difficulty/silent defects
    // in interpreting it. But it smaller than 52 bits, so a signed 64 bit integer or double-precision float type
    // are safe for storing this identifier.
    Migrate_to_chat_id uint64
    // Optional. The supergroup has been migrated from a group with the specified identifier.
    // This number may be greater than 32 bits and some programming languages may have difficulty/silent defects
    // in interpreting it. But it smaller than 52 bits, so a signed 64 bit integer or double-precision float type
    // are safe for storing this identifier.
    Migrate_from_chat_id uint64
    // Optional. Specified message was pinned.
    // Note that the Message object in this field will not contain further reply_to_message
    // fields even if it is itself a reply.
    //Pinned_message Message
}

type MessageEntity struct {
    // Type of the entity. Can be mention (@username), hashtag, bot_command, url, email, bold (bold text), italic (italic text), code (monowidth string), pre (monowidth block), text_link (for clickable text URLs), text_mention (for users without usernames)
    Type string 
    // Offset in UTF-16 code units to the start of the entity
    Offset int 
    // Length of the entity in UTF-16 code units
    Length int 
    // Optional. For “text_link” only, url that will be opened after user taps on the text
    Url string 
    // Optional. For “text_mention” only, the mentioned user
    User User 
}

type PhotoSize struct {
    // Unique identifier for this file
    File_id string 
    // Photo width
    Width int 
    // Photo height
    Height int 
    // Optional. File size
    File_size int 
}

type Audio struct {
    // Unique identifier for this file
    File_id string 
    // Duration of the audio in seconds as defined by sender
    Duration int 
    // Optional. Performer of the audio as defined by sender or by audio tags
    Performer string 
    // Optional. Title of the audio as defined by sender or by audio tags
    Title string 
    // Optional. MIME type of the file as defined by sender
    Mime_type string 
    // Optional. File size
    File_size int 
}

type Voice struct {
    // Unique identifier for this file
    File_id string 
    // Duration of the audio in seconds as defined by sender
    Duration int 
    // Optional. MIME type of the file as defined by sender
    Mime_type string 
    // Optional. File size
    File_size int 
}

type Video struct {
    // Unique identifier for this file
    File_id string 
    // Video width as defined by sender
    Width int 
    // Video height as defined by sender
    Height int 
    // Duration of the video in seconds as defined by sender
    Duration int 
    // Optional. Video thumbnail
    Thumb PhotoSize 
    // Optional. Mime type of a file as defined by sender
    Mime_type string 
    // Optional. File size
    File_size int 
}

type Sticker struct {
    // Unique identifier for this file
    File_id string 
    // Sticker width
    Width int 
    // Sticker height
    Height int 
    // Optional. Sticker thumbnail in .webp or .jpg format
    Thumb PhotoSize 
    // Optional. Emoji associated with the sticker
    Emoji string 
    // Optional. File size
    File_size int 
}

type Document struct {
    // Unique file identifier
    File_id string 
    // Optional. Document thumbnail as defined by sender
    Thumb PhotoSize 
    // Optional. Original filename as defined by sender
    File_name string 
    // Optional. MIME type of the file as defined by sender
    Mime_type string 
    // Optional. File size
    File_size int 
}

type Contact struct {
    // Contact's phone number
    Phone_number string 
    // Contact's first name
    First_name string 
    // Optional. Contact's last name
    Last_name string 
    // Optional. Contact's user identifier in Telegram
    User_id int 
}

type Location struct {
    // Longitude as defined by sender
    Longitude float64 
    // Latitude as defined by sender
    Latitude float64 
}

type Venue struct {
    // Venue location
    Location Location 
    // Name of the venue
    Title string 
    // Address of the venue
    Address string 
    // Optional. Foursquare identifier of the venue
    Foursquare_id string 
}

type UserProfilePhotos struct {
    // Total number of profile pictures the target user has
    Total_count int 
    // Requested profile pictures (in up to 4 sizes each)
    Photos [][]PhotoSize 
}

type File struct {
    // Unique identifier for this file
    File_id string 
    // Optional. File size, if known
    File_size string 
    // Optional. File path. Use https://api.telegram.org/file/bot<token>/<file_path> to get the file.
    File_path string 
}