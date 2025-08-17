package schemas

type SessionCreateRequest struct {
	MockID string `json:"mock_id" validate:"required"`
}

type AnswerAddRequest struct {
	MockID string `json:"mock_id" validate:"required"`
	QuestionID string `json:"question_id" validate:"required"`
	OptionID string `json:"option_id" validate:"required"`
}