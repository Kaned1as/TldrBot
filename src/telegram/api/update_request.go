package api

type GetUpdatesRequest struct {
    offset int // last update_id known to server
    limit int  // optional, from 1 to 100
    timeout int // timeout for polling (usually 0)
}
