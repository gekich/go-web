package dto

type UserResponse struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Bio   *string `json:"bio"`
}

type UpdateUserInput struct {
	Name string  `json:"name" example:"John Doe"`
	Bio  *string `json:"bio" example:"Creative strategist and digital marketing expert"`
}
