package models

type Flat struct {
	ID      int `json:"id"`
	HouseID int `json:"house_id"`
	//Number  int        `json:"number"`
	Price  int        `json:"price"`
	Room   int        `json:"room"`
	Status FlatStatus `json:"status"`
}

type FlatStatus = string

const (
	Created      FlatStatus = "created"
	Approved     FlatStatus = "approved"
	Declined     FlatStatus = "declined"
	OnModeration FlatStatus = "on moderation"
)
