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
