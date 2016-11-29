package impl

import (
    "time"
    "testing"
    "net/http"
    "telegram/api"
)

func TestLeet(t *testing.T) {
	now := time.Now()

    poller := Poller{}
    poll.client = http.Client{Timeout: time.Second * 5}
    poll.today = now
    poll.db := newDb("/tmp/scores_test.db")

    msg := api.Message{}
    msg.Time = time.Date(now.Year(), now.Month(), now.Day(), 1, 2, 3, 4, time.UTC)

    poller.handleL33t(&msg, "1337", LEET_REGEX[4])
}
