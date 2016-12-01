package impl

import (
    "time"
    "testing"
    "net/http"
    "telegram/api"
    "fmt"
)

type FakeBotClient struct {
    t *testing.T
    invocationCount uint8
    desiredCount uint8
}


func (client *FakeBotClient) GetUpdates(lastId int64) *api.GetUpdatesResponse {
    return nil
}

func (client *FakeBotClient) SendObject(request interface{}, url string) *http.Response {
    fmt.Printf("Object that was about to be sent: %#v\n", request)
    if client.invocationCount > client.desiredCount {
        client.t.Error("Should only fire once on leet match, no repetition")
    }

    client.invocationCount += 1
    return nil
}

func (client *FakeBotClient) verify() {
    if client.invocationCount < client.desiredCount {
        client.t.Errorf("Should be invoked at least %d times!", client.desiredCount)
    }

    client.t.Logf("Verified, fired exactly %d times", client.invocationCount)
}

func TestExploitNotWorking(suite *testing.T) {
	now := time.Now()
    fClient := &FakeBotClient{t: suite, desiredCount: 1}

    poller := Poller{}
    poller.client = fClient
    poller.today = now
    poller.db = newDb("/tmp/scores_test.db")

    msg := api.Message{}
    msg.From = api.User{First_name: "Bayaz"}

    msg.Date = time.Date(now.Year(), now.Month(), now.Day(), 13, 37, 00, 0, now.Location()).Unix()
    poller.handleL33t(&msg, "1337", LEET_REGEX[4])
    msg.Date = time.Date(now.Year(), now.Month(), now.Day(), 13, 37, 30, 0, now.Location()).Unix()
    poller.handleL33t(&msg, "1337", LEET_REGEX[4])
    msg.Date = time.Date(now.Year(), now.Month(), now.Day(), 13, 37, 59, 0, now.Location()).Unix()
    poller.handleL33t(&msg, "1337", LEET_REGEX[4])

    fClient.verify()
}


func TestScoreMustBeSix(suite *testing.T) {
    now := time.Now()

    poller := Poller{}
    poller.client = &FakeBotClient{t: suite}
    /*poller.client.SendObject = func(request interface{}, url string) *http.Response {
        fmt.Printf("Object that was about to be sent: %#v\n", request)
        sm := request.(api.SendMessage)
        if !strings.Contains(sm.Text, "scored 6 points") {
            suite.Error("Must score six points for this message time!")
        }
        return nil
    }*/
    poller.today = now
    poller.db = newDb("/tmp/scores_test.db")

    msg := api.Message{}
    msg.From = api.User{First_name: "Bayaz"}

    msg.Date = time.Date(now.Year(), now.Month(), now.Day(), 12, 34, 56, 0, now.Location()).Unix()
    poller.handleL33t(&msg, "123456", LEET_REGEX[3])
}
