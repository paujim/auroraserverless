package entities

type Profile struct {
	ID          int64  `json:"-"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
