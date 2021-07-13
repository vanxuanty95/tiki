package storages

import (
	"context"
)

//go:generate mockgen -package mock -destination mock/store_mock.go . Store
type Store interface {
	RetrieveBookings(ctx context.Context, userID, createdDate string, page, limit int) ([]*Booking, error)
	AddBooking(ctx context.Context, t *Booking) error
}
