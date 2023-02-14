package entity

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Salt     string `json:"salt"`
	Password string `json:"password"`
}
