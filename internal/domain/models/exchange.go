package models

import "github.com/google/uuid"

type DummyLoginRequest struct {
	UserType UserType `json:"user_type"`
}

type DummyLoginResponse struct {
	Token string `json:"token"`
}

type HouseCreateRequest struct {
	Address   string  `json:"address"`
	Year      *int    `json:"year"`
	Developer *string `json:"developer"`
}

type HouseCreateResponse struct {
	House
}

type HouseGetFlatsRequest struct {
	ID int `json:"id"`
}

type HouseGetFlatsResponse struct {
	House
	Flats []Flat `json:"flats"`
}

type FlatCreateRequest struct {
	//ID      int  `json:"id"`
	HouseID int  `json:"house_id"`
	Price   *int `json:"price"`
	Room    int  `json:"room"`
}

type FlatCreateResponse struct {
	Flat
}

type FlatUpdateRequest struct {
	ID     int         `json:"id"`
	Price  *int        `json:"price"`
	Room   int         `json:"room"`
	Status *FlatStatus `json:"status"`
}

type FlatUpdateResponse struct {
	Flat
}

type LoginRequest struct {
	UserID   *uuid.UUID `json:"id"`
	Password *string    `json:"password"`
}

type LoginResponse struct {
	DummyLoginResponse
}

type RegisterRequest struct {
	Email    *string   `json:"email"`
	Password *string   `json:"password"`
	UserType *UserType `json:"user_type"`
}

type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

type SubscribeRequest struct {
	Email *string `json:"email"`
}
