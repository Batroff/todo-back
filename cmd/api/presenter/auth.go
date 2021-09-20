package presenter

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Request AuthRequest `json:"request,omitempty"`
	Token   string      `json:"token"`
	Msg     string      `json:"msg"`
}

func NewAuthResponse(request AuthRequest, token, msg string) *AuthResponse {
	return &AuthResponse{
		Request: request,
		Token:   token,
		Msg:     msg,
	}
}
