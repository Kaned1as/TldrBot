package impl

import (
    "telegram/api"
    "encoding/json"
    "net/http"
    "bytes"
    "log"
    "fmt"
)

type TelegramRestClient interface {
    GetUpdates(lastId int64) *api.GetUpdatesResponse
    GetChat(chatId int64) *api.Chat
    GetChatMember(chatId int64, userId int64) *api.ChatMember
    SendObject(request interface{}, url string) *http.Response
}

type BotRestClient struct {
     http.Client   // http client to operate on
}

func (client *BotRestClient) GetUpdates(lastId int64) *api.GetUpdatesResponse {
    // create request
    request := api.GetUpdatesRequest{Offset: lastId + 1}
    resp := client.SendObject(request, API_ENDPOINT + BOT_TOKEN + GET_UPDATES_PATH)
    if resp == nil {
        return nil // no response, should be logged
    }

    // parse response
    decoder := json.NewDecoder(resp.Body)
    var updates api.GetUpdatesResponse
    parseErr := decoder.Decode(&updates)
    if parseErr != nil {
        fmt.Println("Error parsing response ..." + parseErr.Error())
        return nil
    }

    return &updates
}

func (client *BotRestClient) GetChat(chatId int64) *api.Chat {
    request := api.GetChatRequest{ChatId: chatId}
    resp := client.SendObject(request, API_ENDPOINT + BOT_TOKEN + GET_CHAT_PATH)
    if resp == nil {
        return nil // no response, should be logged
    }

    // parse response
    decoder := json.NewDecoder(resp.Body)
    var chatResp api.GetChatResponse
    parseErr := decoder.Decode(&chatResp)
    if parseErr != nil {
        fmt.Println("Error parsing response ..." + parseErr.Error())
        return nil
    }

    if (!chatResp.Ok) {
        fmt.Println("Non-successful GetChat response ...")
        return nil
    }

    return chatResp.Result
}

func (client *BotRestClient) GetChatMember(chatId int64, userId int64) *api.ChatMember {
    request := api.GetChatMemberRequest{ChatId: chatId, UserId: userId}
    resp := client.SendObject(request, API_ENDPOINT + BOT_TOKEN + GET_CHAT_MEMBER_PATH)
    if resp == nil {
        return nil // no response, should be logged
    }

    // parse response
    decoder := json.NewDecoder(resp.Body)
    var chatResp api.GetChatMemberResponse
    parseErr := decoder.Decode(&chatResp)
    if parseErr != nil {
        fmt.Println("Error parsing response ..." + parseErr.Error())
        return nil
    }

    if (!chatResp.Ok) {
        fmt.Println("Non-successful GetChatMember response ...")
        return nil
    }

    return chatResp.Result
}

func (client *BotRestClient) SendObject(request interface{}, url string) *http.Response {
    body, _ := json.Marshal(request)
    req, createErr := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if createErr != nil {
        log.Fatal("Error creating http request, shutting down ..." + createErr.Error())
    }

    // call to server
    req.Header.Set("Content-Type", "application/json")
    resp, postErr := client.Do(req)
    if postErr != nil {
        fmt.Println("Something is wrong with Telegram servers ..." + postErr.Error())
        return nil
    }

    return resp
}