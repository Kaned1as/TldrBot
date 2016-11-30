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
    "strings"
)

var /* const */ LEET_REGEX = []*regexp.Regexp {
    regexp.MustCompile("0000(0|00)?"), // more is pointless, there is no 00 day or month or year
    regexp.MustCompile("0123(4|45)?"), // 6789 are pointless, we can match only HHmmss
    regexp.MustCompile("1111(1|11|1111|111111)?"), // so max 1111111111 matches MMddHHmmss
    regexp.MustCompile("1234(5|56)?"), // 789 are pointless, see #2
    regexp.MustCompile("1337"), // best regexp, truly
    regexp.MustCompile("2222(2|22|2222)?"), // so max 22222222 matches ddHHmmss
    regexp.MustCompile("(01)?2345"), // first day of month - 012345 matches ddHHmm
}

// Poller for bot updates
// repeatedly retrieves updates from the telegram server and updates its read position
// Note: runs in its own thread
type Poller struct {
    client http.Client   // http client to operate on
    Token string         // bot token
    lastId int64         // last Telegram update ID we know
    db *PersistenceLayer // database
    scoresToday []Score
    today time.Time
}

// starts poller
func (poll *Poller) Start(waiter *sync.WaitGroup) {
    // init
    poll.client = http.Client{Timeout: time.Second * 5}
    ticker := time.NewTicker(POLLING_INTERVAL)
    poll.db = newDb("data/scores.db")
    poll.today = time.Unix(0, 0)

    go poll.doWork(ticker, waiter)
}

func (poll *Poller) doWork(ticker *time.Ticker, waiter *sync.WaitGroup) {
    defer poll.db.Close()
    defer waiter.Done()

    for tick := range ticker.C {
        if tick.Day() > poll.today.Day() {
            poll.scoresToday = []Score{} // empty the list for today's scores
        }

        // create request
        request := api.GetUpdatesRequest{Offset: poll.lastId + 1}
        resp := poll.sendObject(request, API_ENDPOINT + poll.Token + GET_UPDATES_PATH)
        if resp == nil {
            continue // no response, should be logged
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
            poll.handleUpdate(upd)
        }
        fmt.Printf("time: %v, updates: %#v\n", tick, updates);
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
    currTime := time.Unix(msg.Date, 0)

    hmTimeStr := currTime.Format("1504")
    if !strings.HasPrefix(regex.String(), hmTimeStr) {
        return // score must match HHmm
    }
                          //    yyyyMMddHHmmss
    timeStr := currTime.Format("20060102150405")
    match := regex.FindStringSubmatch(timeStr)
    if match == nil {
        return // not enough to be l33t!
    }

    // make sure person hasn't scored this one yet...
    for _, score := range poll.scoresToday {
        if score.PersonId == msg.From.Id && currTime.Sub(score.Time).Seconds() < 60 {
            fmt.Printf("Exploit user detected: %s, score time %s, current time %s\n", 
                msg.From.First_name, score.Time.String(), currTime.String());
            return  // diff between scores is less than minute? You dirty jackass!
        }
    }

    // someone scored, save it and report
    scoredSize := utf8.RuneCountInString(scored) // e.g. we sent 111111 - wanted to score 11 seconds too, it'll be 6
    matchedSize := utf8.RuneCountInString(match[0]) // if time is 11:11:12 we'll get only 5
    totalMatched := uint8(min(scoredSize, matchedSize))

    score := Score{PersonId: msg.From.Id, Grade: totalMatched, Time: currTime}
    poll.scoresToday = append(poll.scoresToday, score) // save it to today's scores
    poll.db.SaveScore(&score) // save it to DB

    disablePreview := new(bool); *disablePreview = true
    report := fmt.Sprintf("%s is a l33t now, scored %d points, current time is %s",
        msg.From.First_name,
        totalMatched,
        currTime.String())

    request := api.SendMessage{
        Text: report,
        ChatId: msg.Chat.Id,
        ParseMode: "Markdown",
        ReplyToMessageId: &msg.Message_id,
        DisableWebPagePreview: disablePreview}

    resp := poll.sendObject(request, API_ENDPOINT + poll.Token + SEND_MESSAGE_PATH)
    if resp == nil {
        return // no response, should be logged
    }

    defer resp.Body.Close()
    fmt.Println(resp.Status)
    io.Copy(os.Stdout, resp.Body)
}

func (poll *Poller) sendObject(request interface{}, url string) (*http.Response) {
    body, _ := json.Marshal(request)
    req, createErr := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if createErr != nil {
        log.Fatal("Error creating http request, shutting down ..." + createErr.Error())
    }

    // call to server
    req.Header.Set("Content-Type", "application/json")
    resp, postErr := poll.client.Do(req)
    if postErr != nil {
        fmt.Println("Something is wrong with Telegram servers ..." + postErr.Error())
        return nil
    }

    return resp
}
