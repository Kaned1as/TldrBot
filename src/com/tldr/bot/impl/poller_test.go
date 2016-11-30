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
    poller.client = http.Client{Timeout: time.Second * 5}
    poller.today = now
    poller.db = newDb("/tmp/scores_test.db")

    msg := api.Message{}
    msg.Date = time.Date(now.Year(), now.Month(), now.Day(), 
                         13, 37, 00, 0, now.Location()).Unix()

    poller.handleL33t(&msg, "1337", LEET_REGEX[4])
}
