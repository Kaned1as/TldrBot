package impl

import (
    "telegram/api"
    "time"
)

const API_ENDPOINT string = "https://api.telegram.org/bot"
const GET_UPDATES_PATH = "/getUpdates"
const SEND_MESSAGE_PATH = "/sendMessage"
const POLLING_INTERVAL = time.Millisecond * 300

// filled in main()
var BOT_TOKEN string

type CaughtUrl struct {
    // text that was catched from URL
    Text string

    // chat where this occurred
    Chat *api.Chat
}

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}