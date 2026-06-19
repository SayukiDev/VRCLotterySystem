package http

type CommonResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any
}
