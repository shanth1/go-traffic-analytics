package domain

type SendFeedbackCmd struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	UserID  UserID `json:"user_id"`
}
