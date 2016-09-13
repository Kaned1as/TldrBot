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
    "log"
    "github.com/advancedlogic/GoOse"
    "strings"
    "unicode/utf8"
)

// Poller for bot updates
// repeatedly retrieves updates from the telegram server and updates its read position
// Note: runs in its own thread
type Poller struct {
    client http.Client
    articleFetcher goose.Goose
    Msg chan CaughtUrl
    lastId int64
}

// starts poller
func (poll *Poller) Start(waiter *sync.WaitGroup) {
    poll.client = http.Client{Timeout: time.Second * 5}
    poll.articleFetcher = goose.New()
    ticker := time.NewTicker(POLLING_INTERVAL)
    go poll.doWork(ticker, waiter)
}

func (this *Poller) doWork(tick *time.Ticker, waiter *sync.WaitGroup) {
    defer waiter.Done()

    for time := range tick.C {
        // create request
        request := api.GetUpdatesRequest{Offset: this.lastId + 1}
        body, _ := json.Marshal(request)
        req, createErr := http.NewRequest("POST", API_ENDPOINT + BOT_TOKEN + GET_UPDATES_PATH, bytes.NewBuffer(body))
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

        if len(updates.Result) > 0 {
            fmt.Printf("time: %v, updates: %#v\n", time, updates);
        }

        for _, upd := range updates.Result {
            this.handleUpdate(upd)
        }
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
    go poll.handlePage(url, &msg.Chat)
}

func newlineOrComma(r rune) bool {
    return r == '.' || r == '\n'
}

func (poll *Poller) handlePage(url string, chat *api.Chat)  {
    article, parseErr := poll.articleFetcher.ExtractFromURL(url)
    if parseErr != nil {
        fmt.Println("Something is wrong with this page, skipping ..." + parseErr.Error())
    }

    mainContent := article.CleanedText
    sentences := strings.FieldsFunc(mainContent, newlineOrComma)
    valuableSentences := []string{}
    for _, sentence := range sentences {
        if (utf8.RuneCountInString(sentence) > 5 && utf8.RuneCountInString(sentence) < 70) {
            continue  // too short sentence can be addressing like Mr. or Jr.
        }
        valuableSentences = append(valuableSentences, sentence)
    }

    poll.Msg <- CaughtUrl{strings.Join(valuableSentences, "."), chat}
}
