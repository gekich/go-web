package dto

type ProfileResponse struct {
	ID    string  `json:"id"`
	Email string  `json:"email"`
	Name  string  `json:"name"`
	Bio   *string `json:"bio"`
}
