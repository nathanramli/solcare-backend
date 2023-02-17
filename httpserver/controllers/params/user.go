package params

import "mime/multipart"

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

type UpdateUserAvatar struct {
	Picture *multipart.FileHeader `form:"picture" validate:"required"`
}

type RequestKyc struct {
	Nik            string                `form:"nik" validate:"required,len=16"`
	IdCard         *multipart.FileHeader `form:"idCard" validate:"required"`
	Face           *multipart.FileHeader `form:"face" validate:"required"`
	FaceWithIdCard *multipart.FileHeader `form:"faceWithIdCard" validate:"required"`
}
