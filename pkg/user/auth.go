package user

type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
