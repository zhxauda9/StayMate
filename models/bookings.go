package models

type Booking struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	RoomID   int    `json:"room_id"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
}
