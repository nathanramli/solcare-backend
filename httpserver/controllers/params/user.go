package params

type Login struct {
	Address       string `json:"address" validate:"required"`
	SignedMessage string `json:"signedMessage" validate:"required"`
}

type UpdateUser struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Gender    *bool  `json:"gender" validate:"required"`
}
