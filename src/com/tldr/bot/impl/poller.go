package impl

import (
    "fmt"
    "encoding/json"
    "telegram/api"
    "sync"
    "net/http"
    "time"
    "bytes"
    "unicode/utf16"
    "io/ioutil"
    "log"
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
func (poll *Poller) Start(waiter *sync.WaitGroup) {
    defer waiter.Done()

    poll.client = http.Client{Timeout: time.Second * 5}
    ticker := time.NewTicker(POLLING_INTERVAL)
    go poll.doWork(ticker)
}

func (this *Poller) doWork(tick *time.Ticker) {
    for time := range tick.C {
        // create request
        request := api.GetUpdatesRequest{Offset: this.lastId + 1}
        body, _ := json.Marshal(request)
        req, createErr := http.NewRequest("POST", API_ENDPOINT + this.Token + GET_UPDATES_PATH, bytes.NewBuffer(body))
        if createErr != nil {
            log.Fatal("Error creating http request, shutting down ..." + createErr.Error())
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

func (poll *Poller) handleUpdate(update api.Update) {
    if poll.lastId < update.Update_id {
        poll.lastId = update.Update_id
    }

    // we're interested only in messages
    msg := update.Message
    if (msg == nil) {
        return // not message update, skip
    }

    var url string
    for _, entity := range msg.Entities {
        if entity.Type == "url" {
            textAsUtf16 := utf16.Encode([]rune(msg.Text))
            urlAsUtf16 := textAsUtf16[entity.Offset:entity.Offset + entity.Length]
            url = string(utf16.Decode(urlAsUtf16))
            break // we need only first URL we encounter
        }
    }

    if url == "" {
        return // no URLs found
    }
    fmt.Printf("Got message with URL: %#v\n", url);

    // retrieve page
    go poll.getPage(url)
}
func (poll *Poller) getPage(url string) {
    resp, getErr := poll.client.Get(url)
    if getErr != nil {
        fmt.Println("Error retrieving url mentioned, skipping ..." + getErr.Error())
        return
    }

    bytes, readErr := ioutil.ReadAll(resp.Body)
    if (readErr != nil) {
        fmt.Println("Error reading response body!" + readErr.Error())
        return
    }

    html := string(bytes)
}
