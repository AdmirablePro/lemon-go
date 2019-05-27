package main

type Task struct {
	TaskID     string            `json:"mid"`
	HTTPMethod string            `json:"method"`
	Scheme     string            `json:"scheme"`
	Host       string            `json:"host"`
	Path       string            `json:"path"`
	Header     map[string]string `json:"header"`
	Param      map[string]string `json:"param"`
	Cookie     string            `json:"cookie"`
	CookieID   string            `json:"cid"`
	Payload    string            `json:"data"`
}

type Result struct {
	Status       string `json:"status"`    // task status
	TaskID       string `json:"mid"`       // task id
	FetchedTime  int64  `json:"time"`      // fetch time
	CookieID     string `json:"cid"`       // cookie id
	ResponseCode int    `json:"code"`      // HTTP response code
	Data         string `json:"data"`      // HTTP response body
	User         string `json:"user"`      // User identifier
	ErrorCode    int    `json:"errorCode"` // error code
}

type Configuration struct {
	ClientID string `json:"uuid"`
}
