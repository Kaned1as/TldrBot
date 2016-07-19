package main

import (
//    "telegram/api"
    "log"
    "io/ioutil"
    "fmt"
    "net/http"
    "time"
    "sync"
    "encoding/json"
)

const API_ENDPOINT string = "https://api.telegram.org/bot"
const GET_UPDATES_PATH = "/getUpdates"

const POLLING_INTERVAL = time.Millisecond * 300

// Poller for bot updates
// repeatedly retrieves updates from the telegram server and updates its read position
// Note: runs in its own thread
type Poller struct {
    client http.Client
    token string
    lastId uint64
}

// starts poller
func (this Poller) start(waiter *sync.WaitGroup) {
    defer waiter.Done()

    this.client = http.Client{Timeout: time.Second * 5}
    ticker := time.NewTicker(POLLING_INTERVAL)
    go this.doWork(ticker)
}

func (this Poller) doWork(tick *time.Ticker) {
    for t := range tick.C {
        // call to server
        resp, err := this.client.Get(API_ENDPOINT + this.token + GET_UPDATES_PATH)
        if err != nil {
            fmt.Println("Something is wrong with Telegram servers ..." + err.Error())
            continue
        }

        // parse response
        decoder := json.NewDecoder(resp.Body)
        //var updates api.GetUpdatesResponse
        var updates  map[string]interface{}
        parseErr := decoder.Decode(&updates)
        if parseErr != nil {
            fmt.Println("Error parsing response ..." + parseErr.Error())
            continue
        }

        fmt.Println(t)
        fmt.Println(updates)
    }
}

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

    // 1 - poller
    poller := Poller{token: token}
    poller.start(&join)

    // 2 - response publisher
    // TODO

    join.Wait()
}
