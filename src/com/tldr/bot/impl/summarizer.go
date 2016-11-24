package impl

import (
    "sync"
    "fmt"
    "time"
    "unicode/utf8"
    "telegram/api"
)

type Summarizer struct {
    // sentence count to split up to
    sentCount int
    MsgChannel chan *api.Message
}

func (summ *Summarizer) Start(waiter *sync.WaitGroup) {
    go summ.doWork(waiter)
}

func (summ *Summarizer) doWork(waiter *sync.WaitGroup) {
    defer waiter.Done()
    answerCount := 0
    for msg := range summ.MsgChannel {
        timeStr := time.Now().Format("20060201150405")

        matchLen := 0
        matchIdx := 0 // index in incoming string
        for _, matchChar := range timeStr {

        }
        fmt.Printf("Get ID %#v \n", answerCount)
        answerCount++
        go postWork(summary)
    }
}

func postWork(summary string) {
    fmt.Println(summary)
}
