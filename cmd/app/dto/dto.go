package dto

type TokenDTO struct {
	Token string `json:"token"`
}

type PaymentDTO struct {
	Id          int    `json:"id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type SuggestionDTO struct {
	Token  string `json:"token"`
	UserId int    `json:"userid"`
	Sugid  int    `json:"sugid"`
	Icon   string `json:"icon"`
	Title  string `json:"title"`
	Link   string `json:"link"`
}
