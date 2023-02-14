package v1

type RequestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ResultUser struct {
	Email    string `json:"email"`
	Salt     string `json:"salt"`
	Password string `json:"password"`
}
