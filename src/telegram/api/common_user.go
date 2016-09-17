package api

type User struct {
    // Unique identifier for this user or bot
    Id int64
    // User‘s or bot’s first name
    First_name string 
    // Optional. User‘s or bot’s last name
    Last_name string 
    // Optional. User‘s or bot’s username
    Username string 
}
