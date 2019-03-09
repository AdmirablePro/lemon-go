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
	Status       string `json:"status"`
	TaskID       string `json:"mid"`
	CookieID     string `json:"cid"`
	ResponseCode int    `json:"code"`
	Data         string `json:"data"`
	FetchedTime  int64  `json:"time"`
	User         string `json:"user"`
}

type Configuration struct {
	ClientID string `json:"uuid"`
}
