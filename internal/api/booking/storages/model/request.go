package model

type BookingRequest struct {
	ScreenID  int         `json:"screen_id" validate:"required"`
	Number    int         `json:"number"`
	Locations *[]Location `json:"locations"`
}
type CheckAvailableRequest struct {
	ScreenID int       `json:"screen_id" validate:"required"`
	Location *Location `json:"location" validate:"required"`
}
type Location struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}
