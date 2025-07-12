package schemas

import "time"

type APIResponse struct {
	Ok       bool        `json:"ok"`
	D        interface{} `json:"d"`
	ErrorMsg string      `json:"error_message"`

	ProcessedAt time.Time `json:"processed_at"`
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
	return NewAPIResponse(false, d, errorMsg)
}