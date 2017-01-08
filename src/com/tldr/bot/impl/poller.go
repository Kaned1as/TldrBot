package impl

import (
    "fmt"
    "regexp"
    "telegram/api"
    "sync"
    "net/http"
    "time"
    "unicode/utf8"
    "io"
    "os"
    "strings"
)

var /* const */ LEET_REGEX = []*regexp.Regexp {
    regexp.MustCompile("0000(00|0)?"), // more is pointless, there is no 00 day or month or year
    regexp.MustCompile("0123(45|4)?"), // 6789 are pointless, we can match only HHmmss
    regexp.MustCompile("1111(111111|1111|11|1)?"), // so max 1111111111 matches MMddHHmmss
    regexp.MustCompile("1234(56|5)?"), // 789 are pointless, see #2
    regexp.MustCompile("1337"), // best regexp, truly
    regexp.MustCompile("2222(2222|22|2)?"), // so max 22222222 matches ddHHmmss
    regexp.MustCompile("(01)?2345"), // first day of month - 012345 matches ddHHmm
}

// Poller for bot updates
// repeatedly retrieves updates from the telegram server and updates its read position
// Note: runs in its own thread
type Poller struct {
    client TelegramRestClient
    lastId int64         // last Telegram update ID we know
    db *PersistenceLayer // database
    scoresToday []Score
    today time.Time
}

// starts poller
func (poll *Poller) Start(waiter *sync.WaitGroup) {
    // init
    poll.client = &BotRestClient{http.Client {Timeout: time.Second * 5}}
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
            poll.today = tick
        }

        updates := poll.client.GetUpdates(poll.lastId)
        if updates == nil {
            continue
        }

        for _, upd := range updates.Result {
            poll.handleUpdate(upd)
        }
        fmt.Printf("time: %v, updates: %#v\n", tick, updates)
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

    if msg.Text == "/stat" || msg.Text == "/stat@l33t_count_bot" {
        go poll.handleStat(msg)
        return
    }

    if msg.Text == "/totals" || msg.Text == "/totals@l33t_count_bot" {
        go poll.handleTotals(msg)
        return
    }

    fmt.Printf("Got message with text: %#v\n", msg.Text)
    
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

func (poll *Poller) handleStat(msg *api.Message) {
    total, highest, latest := poll.db.GetScores(msg.From.Id)
    report := fmt.Sprintf("Stats for %s:\n" +
                          "  Total scored: %d points\n",
        msg.From.First_name,
        total)

    if highest != nil {
        report += fmt.Sprintf("  Highest score: %d points at %s\n",
            highest.Grade,
            highest.Time.Format("2006-Jan-02 15:04:05"))
    }

    if (latest != nil) {
        report += fmt.Sprintf("  Latest score: %d points at %s",
            latest.Grade,
            latest.Time.Format("2006-Jan-02 15:04:05"))
    }

    disablePreview := new(bool); *disablePreview = true
    request := api.SendMessage{
        Text: report,
        ChatId: msg.Chat.Id,
        ParseMode: "Markdown",
        ReplyToMessageId: &msg.Message_id,
        DisableWebPagePreview: disablePreview}

    resp := poll.client.SendObject(request, API_ENDPOINT + BOT_TOKEN + SEND_MESSAGE_PATH)
    if resp == nil {
        return // no response, should be logged
    }

    defer resp.Body.Close()
    fmt.Println(resp.Status)
    io.Copy(os.Stdout, resp.Body)
}

func (poll *Poller) handleTotals(msg *api.Message) {
    totals := poll.db.GetTotals()

    disablePreview := new(bool); *disablePreview = true
    answer := api.SendMessage{
        ChatId: msg.Chat.Id,
        ParseMode: "Markdown",
        ReplyToMessageId: &msg.Message_id,
        DisableWebPagePreview: disablePreview}

    // no stats available
    if (len(totals) == 0) {
        answer.Text = "No stats available for this chat! Sorry..."
        poll.client.SendObject(answer, API_ENDPOINT + BOT_TOKEN + SEND_MESSAGE_PATH)
        return
    }

    answer.Text = "Total stats for this chat members:\n"
    counter := 0
    for idx := range totals {
        member := poll.client.GetChatMember(msg.Chat.Id, totals[idx].PersonId)
        if member == nil { // no such person in this chat!
            continue
        }
        counter++
        answer.Text += fmt.Sprintf("  %d. %s - %d points\n", counter, member.User.First_name, totals[idx].TotalScored);
    }

    poll.client.SendObject(answer, API_ENDPOINT + BOT_TOKEN + SEND_MESSAGE_PATH)
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
            fmt.Printf("Exploit user detected: %s, score time %s, message time %s\n",
                msg.From.First_name, score.Time.String(), currTime.String())
            return  // diff between scores is less than minute? You dirty jackass!
        }
    }

    // someone scored, save it and report
    fmt.Printf("Timestring: %s, Message: %s Matched: %s\n", timeStr, scored, match)
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

    resp := poll.client.SendObject(request, API_ENDPOINT + BOT_TOKEN + SEND_MESSAGE_PATH)
    if resp == nil {
        return // no response, should be logged
    }

    defer resp.Body.Close()
    fmt.Println(resp.Status)
    io.Copy(os.Stdout, resp.Body)
}
