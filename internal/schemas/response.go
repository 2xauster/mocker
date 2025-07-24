package schemas

import (
	"log/slog"
	"time"
)

type APIResponse struct {
	Ok       bool        `json:"ok"`
	D        interface{} `json:"d"`
	ErrorMsg string      `json:"error_message"`

	ProcessedAt time.Time `json:"processed_at"`
}

type ErrorStringSchema struct {
	Err interface{} `json:"err"`
	Code int `json:"code"`
	Type string `json:"type"`
}

func NewAPIResponse(ok bool, d interface{}, errorMsg string) APIResponse {
	return APIResponse{
		Ok: ok,
		D: d,
		ErrorMsg: errorMsg,

		ProcessedAt: time.Now(),
	}
}

func NewErrorAPIResponse(d interface{}, errorMsg string) APIResponse {
	slog.Error("Error response", "msg", errorMsg, "err", d)

	return NewAPIResponse(false, d, errorMsg)
}