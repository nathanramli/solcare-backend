package params

type Login struct {
	Address       string `json:"address" validate:"required"`
	SignedMessage string `json:"signedMessage" validate:"required"`
}
