package api

type Chat struct {
    // Unique identifier for this chat. This number may be greater than 32 bits
    // and some programming languages may have difficulty/silent defects in interpreting it.
    // But it smaller than 52 bits, so a signed 64 bit integer or double-precision float type
    // are safe for storing this identifier.
    Id int64
    // Type of chat, can be either “private”, “group”, “supergroup” or “channel”
    Type string 
    // Optional. Title, for channels and group chats
    Title string 
    // Optional. Username, for private chats, supergroups and channels if available
    Username string 
    // Optional. First name of the other party in a private chat
    First_name string 
    // Optional. Last name of the other party in a private chat
    Last_name string 

}

type ChatMember struct {
    // Information about the user
    User User 
    // The member's status in the chat. Can be “creator”, “administrator”, “member”, “left” or “kicked”
    Status string 
}