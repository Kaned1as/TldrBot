package impl

import (
    "fmt"
    "encoding/json"
    "telegram/api"
    "sync"
    "net/http"
    "time"
    "bytes"
)

const API_ENDPOINT string = "https://api.telegram.org/bot"
const GET_UPDATES_PATH = "/getUpdates"
const POLLING_INTERVAL = time.Millisecond * 300

// Poller for bot updates
// repeatedly retrieves updates from the telegram server and updates its read position
// Note: runs in its own thread
type Poller struct {
    client http.Client
    Token string
    lastId int64
}

// starts poller
func (this *Poller) Start(waiter *sync.WaitGroup) {
    defer waiter.Done()

    this.client = http.Client{Timeout: time.Second * 5}
    ticker := time.NewTicker(POLLING_INTERVAL)
    go this.doWork(ticker)
}

func (this *Poller) doWork(tick *time.Ticker) {
    for time := range tick.C {
        // create request
        request := api.GetUpdatesRequest{Offset: this.lastId + 1}
        body, _ := json.Marshal(request)
        req, createErr := http.NewRequest("POST", API_ENDPOINT + this.Token + GET_UPDATES_PATH, bytes.NewBuffer(body))
        if createErr != nil {
            fmt.Println("Error creating http request, shutting down ..." + createErr.Error())
            break
        }

        // call to server
        req.Header.Set("Content-Type", "application/json")
        resp, postErr := this.client.Do(req)
        if postErr != nil {
            fmt.Println("Something is wrong with Telegram servers ..." + postErr.Error())
            continue
        }

        // parse response
        decoder := json.NewDecoder(resp.Body)
        var updates api.GetUpdatesResponse
        parseErr := decoder.Decode(&updates)
        if parseErr != nil {
            fmt.Println("Error parsing response ..." + parseErr.Error())
            continue
        }

        for _, upd := range updates.Result {
            this.handleUpdate(upd)
        }
        fmt.Printf("time: %v, updates: %#v\n", time, updates);
    }
}

func (this *Poller) handleUpdate(update api.Update) {
    if this.lastId < update.Update_id {
        this.lastId = update.Update_id
    }
}
