package api

type User struct {
    id int `desc:"Unique identifier for this user or bot"`
    first_name string `desc:"User‘s or bot’s first name"`
    last_name string `desc:"Optional. User‘s or bot’s last name"`
    username string `desc:"Optional. User‘s or bot’s username"`
}
