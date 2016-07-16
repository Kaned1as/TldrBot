package main

import (
//    "telegram/api"
    "log"
    "io/ioutil"
    "fmt"
)

const API_ENDPOINT string = "https://api.telegram.org/bot"

func main() {
    bytes, err := ioutil.ReadFile("bot-token.txt")
    if (err != nil) {
        log.Fatal("Error opening bot token file!" + err.Error())
        return
    }

    token := string(bytes[:])
    fmt.Println(token)
}
