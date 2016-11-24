package main

import (
    "log"
    "com/tldr/bot/impl"
    "telegram/api"
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

    token := string(bytes)
    fmt.Println("Starting to work with bot token: " + token)

    join := sync.WaitGroup{}
    join.Add(2)

    messages := make(chan *api.Message)

    // 1 - poller
    poller := impl.Poller{Token: token, MsgChannel: messages}
    poller.Start(&join)

    // 2 - response publisher
    summarizer := impl.Summarizer{MsgChannel: messages}
    summarizer.Start(&join)

    join.Wait()
}
