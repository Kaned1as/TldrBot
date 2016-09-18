package impl

import (
    "github.com/JesusIslam/tldr"
    "sync"
    "fmt"
    "telegram/api"
    "net/http"
    "time"
    "encoding/json"
    "bytes"
    "log"
    "io"
    "os"
)

type Summarizer struct {
    // sentence count to split up to
    sentCount int
    // channel where works comes from
    Msg chan CaughtUrl
    // http client to post work back
    client http.Client
}

func (summ *Summarizer) Start(waiter *sync.WaitGroup) {
    summ.sentCount = 5
    summ.client = http.Client{Timeout: time.Second * 5}

    go summ.doWork(waiter)
}

func (summ *Summarizer) doWork(waiter *sync.WaitGroup) {
    defer waiter.Done()
    for caught := range summ.Msg {
        bag := tldr.New()
        summary, err := bag.Summarize(caught.Text, summ.sentCount)
        if err != nil {
            fmt.Println("Error inlining summary ..." + err.Error())
            continue
        }
        go summ.postWork(summary, caught.Chat)
    }
}

func (summ *Summarizer)  postWork(summary string, chat *api.Chat) {
    disablePreview := new(bool); *disablePreview = true
    request := api.SendMessage{Text: summary, ChatId: chat.Id, ParseMode: "Markdown",  DisableWebPagePreview: disablePreview}
    body, _ := json.Marshal(request)
    req, createErr := http.NewRequest("POST", API_ENDPOINT + BOT_TOKEN + SEND_MESSAGE_PATH, bytes.NewBuffer(body))
    if createErr != nil {
        log.Fatal("Error creating http request, shutting down ..." + createErr.Error())
    }

    // call to server
    req.Header.Set("Content-Type", "application/json")
    resp, postErr := summ.client.Do(req)
    if postErr != nil {
        fmt.Println("Something is wrong with Telegram servers ..." + postErr.Error())
        return
    }

    defer resp.Body.Close()
    fmt.Println(resp.Status)
    io.Copy(os.Stdout, resp.Body)
}