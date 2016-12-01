package main

import (
    "log"
    "com/tldr/bot/impl"
    "io/ioutil"
    "fmt"
    "sync"
)

func main() {
    bytes, err := ioutil.ReadFile("bot-token.txt")
    if (err != nil) {
        log.Fatal("Error opening bot token file!" + err.Error())
        return
    }

    impl.BOT_TOKEN = string(bytes)
    fmt.Println("Starting to work with bot token: " + impl.BOT_TOKEN)

    join := sync.WaitGroup{}
    join.Add(1)

    poller := impl.Poller{}
    poller.Start(&join)

    join.Wait()
}
