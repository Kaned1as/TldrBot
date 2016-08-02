package impl

import (
    "github.com/JesusIslam/tldr"
    "sync"
    "fmt"
)

type Summarizer struct {
    // sentence count to split up to
    sentCount int
    bag *tldr.Bag
    Msg chan string
}

func (summ *Summarizer) Start(waiter *sync.WaitGroup) {
    summ.bag = tldr.New()
    summ.sentCount = 5

    go summ.doWork(waiter)
}

func (summ *Summarizer) doWork(waiter *sync.WaitGroup) {
    defer waiter.Done()
    answerCount := 0
    for text := range summ.Msg {
        summary, err := summ.bag.Summarize(text, summ.sentCount)
        if (err != nil) {
            fmt.Println("Error inlining summary ..." + err.Error())
            continue
        }
        fmt.Printf("Summarize ID %#v \n", answerCount)
        go postWork(summary)
    }
}

func postWork(summary string) {
    fmt.Println(summary)
}