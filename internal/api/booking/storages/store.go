package storages

import (
	"context"
)

//go:generate mockgen -package mock -destination mock/store_mock.go . Store
type Store interface {
	InsertScreen(ctx context.Context, s *Screen) (int, error)
	GetScreenByID(ctx context.Context, id int) (*Screen, error)

	GetAllSeatByScreenID(ctx context.Context, id int) ([]*Seat, error)
	InsertSeats(ctx context.Context, seats []*Seat) error
}
