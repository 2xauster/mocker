package schemas

type MockCreateRequest struct {
	Topic string `json:"topic"`
	Instructions string `json:"instructions"`
	TimeMins int `json:"time_mins"`
	
	Questions []MockQuestionSchema `json:"questions"`

	AuthorID string `json:"author_id"`
}

type MockQuestionSchema struct {
	Problem string `json:"problem"`
	Points int `json:"points"`
	CorrectOptionID string `json:"correct_option_id"`
	Options []MockOptionSchema `json:"options"`
}

type MockOptionSchema struct {
	Number int `json:"number"`
	Option string `json:"option"`
}