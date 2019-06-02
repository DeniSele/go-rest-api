package model

type UserInfo struct {
	ID          string `json:"id"`
	Firstname   string `json:"firstname"`
	Secondname  string `json:"secondname"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type PersonalAccount struct {
	ID      string `json:"id"`
	IDUser  string `json:"idUser"`
	Balance string `json:"balance"`
}

type Transaction struct {
	ID          string `json:"id"`
	IDPersAcc   string `json:"idPersAcc"`
	TransType   string `json:"transType"`
	Date        string `json:"date"`
	TargetAccID string `json:"targetAccID"`
	Value       string `json:"value"`
}

type RangeDate struct {
	StartDate  string `json:"startDate"`
	FinishDate string `json:"finishDate"`
}
