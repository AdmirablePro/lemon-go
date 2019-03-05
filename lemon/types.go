package main

type Task struct {
	TaskID string            `json:"mid"`
	Method string            `json:"method"`
	Host   string            `json:"host"`
	Path   string            `json:"path"`
	Header []string          `json:"header"`
	Param  map[string]string `json:"param"`
	Cookie string            `json:"cookie"`
}

type Result struct {
	Status       string `json:"status"`
	TaskID       string `json:"task_id"`
	ResponseCode int    `json:"code"`
	Data         string `json:"data"`
	FetchedTime  string `json:"time"`
	UserAgent    string `json:"ua"`
}
