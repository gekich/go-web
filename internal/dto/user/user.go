package dto

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Bio   string `json:"bio" validate:"required"`
}

type UpdateUserInput struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}
