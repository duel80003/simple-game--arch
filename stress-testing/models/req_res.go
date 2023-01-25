package models

type Request struct {
	Topic string                 `json:"topic"`
	Data  map[string]interface{} `json:"data"`
}

type Response struct {
	Topic string                 `json:"topic"`
	Data  map[string]interface{} `json:"data"`
}

type ErrorRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
