package presenter

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
