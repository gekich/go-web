package dto

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

type CreateUserInput struct {
	Name  string `json:"name" validate:"required" example:"John Doe"`
	Email string `json:"email" validate:"required" example:"test@example.test"`
	Bio   string `json:"bio" validate:"required" example:"Creative strategist and digital marketing expert"`
}

type UpdateUserInput struct {
	Name string `json:"name" example:"John Doe"`
	Bio  string `json:"bio" example:"Creative strategist and digital marketing expert"`
}
