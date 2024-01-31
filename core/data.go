package core

type httpResponse struct {
	Error   bool
	Message string
}

type Buyer struct {
	Name    string
	PhoneNr int
	Address string
}

type RegisterBuyer struct {
	Name    string `json:"name"`
	PhoneNr int    `json:"phone_nr"`
	Address string `json:"address"`
}
