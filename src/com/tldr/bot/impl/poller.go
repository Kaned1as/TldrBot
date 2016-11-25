package impl

import (
    "fmt"
    "regexp"
    "encoding/json"
    "telegram/api"
    "sync"
    "net/http"
    "time"
    "bytes"
    "log"
    "unicode/utf8"
    "io"
    "os"
)

var /* const */ LEET_REGEX = []*regexp.Regexp {
    regexp.MustCompile("0000(0*)"),
    regexp.MustCompile("0123(4?5?6?7?8?9?)"),
    regexp.MustCompile("1111(1*)"),
    regexp.MustCompile("1234(5?6?7?8?9?)"),
    regexp.MustCompile("1337"),
    regexp.MustCompile("2222(2*)"),
    regexp.MustCompile("2345(6?7?8?9?)"),
}

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
    poll.client = http.Client{Timeout: time.Second * 5}
    ticker := time.NewTicker(POLLING_INTERVAL)
    go poll.doWork(ticker, waiter)
}

func (this *Poller) doWork(tick *time.Ticker, waiter *sync.WaitGroup) {
    defer waiter.Done()

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
    if (msg == nil || msg.Text == "") {
        return // not message update, skip
    }

    fmt.Printf("Got message with text: %#v\n", msg.Text);
    
    // make sure it's l33t msg
    for _, regex := range LEET_REGEX {
        matches := regex.FindStringSubmatch(msg.Text)
        if matches == nil {
            continue
        }

        // someone has scored, kudos to him
        go poll.handleL33t(msg, matches[0], regex)
        return
    }
}

func (poll *Poller) handleL33t(msg *api.Message, scored string, regex *regexp.Regexp) {
    timeStr := time.Now().Format("20060201150405")
    match := regex.FindStringSubmatch(timeStr)
    if match == nil {
        return // not enough to be l33t!
    }
    scoredSize := utf8.RuneCountInString(scored) // e.g. we sent 111111 - wanted to score 11 seconds too, it'll be 6
    matchedSize := utf8.RuneCountInString(match[0]) // if time is 11:11:12 we'll get only 5
    totalMatched := min(scoredSize, matchedSize)

    disablePreview := new(bool); *disablePreview = true
    report := fmt.Sprintf("%s is a l33t now, scored %d points, current time is %s",
        msg.From.First_name,
        totalMatched,
        time.Now().String())
    request := api.SendMessage{
        Text: report,
        ChatId: msg.Chat.Id,
        ParseMode: "Markdown",
        ReplyToMessageId: &msg.Message_id,
        DisableWebPagePreview: disablePreview}
    body, _ := json.Marshal(request)
    req, createErr := http.NewRequest("POST", API_ENDPOINT + poll.Token + SEND_MESSAGE_PATH, bytes.NewBuffer(body))
    if createErr != nil {
        log.Fatal("Error creating http request, shutting down ..." + createErr.Error())
    }

    // call to server
    req.Header.Set("Content-Type", "application/json")
    resp, postErr := poll.client.Do(req)
    if postErr != nil {
        fmt.Println("Something is wrong with Telegram servers ..." + postErr.Error())
        return
    }

    defer resp.Body.Close()
    fmt.Println(resp.Status)
    io.Copy(os.Stdout, resp.Body)
}
