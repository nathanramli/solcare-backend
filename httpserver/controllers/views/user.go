package views

type Login struct {
	Token string `json:"token"`
}

type FindUser struct {
	Address        string `json:"address"`
	CreatedAt      int64  `json:"createdAt"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Gender         bool   `json:"gender"`
	IsVerified     bool   `json:"isVerified"`
	IsWarned       bool   `json:"isWarned"`
	ProfilePicture string `json:"profilePicture"`
}

type FindKycRequest struct {
	Id                      uint   `json:"id"`
	Nik                     string `json:"nik"`
	Name                    string `json:"name"`
	UsersWalletAddress      string `json:"usersWalletAddress"`
	RequestedAt             int64  `json:"requestedAt"`
	IdCardPicture           string `json:"idCardPicture"`
	FacePicture             string `json:"facePicture"`
	SelfieWithIdCardPicture string `json:"selfieWithIdCardPicture"`
	Status                  uint8  `json:"status"`
}
