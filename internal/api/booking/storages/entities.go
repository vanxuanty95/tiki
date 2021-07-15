package storages

import "time"

// Screen reflects bookings in DB
type Screen struct {
	ID               int       `json:"id"`
	NumberSeatRow    int       `json:"number_seat_row" validate:"required,min=1"`
	NumberSeatColumn int       `json:"number_seat_column" validate:"required,min=1"`
	UserID           string    `json:"user_id"`
	CreatedDate      time.Time `json:"created_date"`
}

// Seat reflects bookings in DB
type Seat struct {
	ID         int       `json:"-"`
	Row        int       `json:"row_id"`
	Column     int       `json:"column_id"`
	UserID     string    `json:"user_id"`
	ScreenID   int       `json:"screen_id"`
	BookedDate time.Time `json:"created_date"`
}
