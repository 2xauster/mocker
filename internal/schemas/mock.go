package schemas

type MockCreateRequest struct {
	Topic string `json:"topic" validate:"required,min=1,max=200"`
	Instructions string `json:"instructions" validate:"required,max=40000"`
	TimeMins int `json:"time_mins" validate:"required,numeric,min=1"`
	
	Questions []MockQuestionSchema `json:"questions" validate:"required,min=1"`

	AuthorID string `json:"author_id"`
}

type MockQuestionSchema struct {
	Problem string `json:"problem" validate:"required,min=1"`
	Points int `json:"points" validate:"required,numeric,min=1"`
	CorrectOptionID string `json:"correct_option_id" validate:"required,min=1"`
	Options []MockOptionSchema `json:"options" validate:"required,min=4"`
}

type MockOptionSchema struct {
	Number int `json:"number" validate:"required,numeric,min=1"`
	Option string `json:"option" validate:"required,min=1"`
}