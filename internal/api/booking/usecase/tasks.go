package usecase

import (
	"context"
	"errors"
	"tiki/internal/api/booking/storages"
	"tiki/internal/api/dictionary"
	userStorages "tiki/internal/api/user/storages"
	"tiki/internal/api/utils"
	"tiki/internal/pkg/logger"
)

type Booking struct {
	Store           storages.Store
	UserStore       userStorages.Store
	GeneratorUUIDFn utils.GenerateNewUUIDFn
}

func (s *Booking) List(ctx context.Context, userID, createdDate string, page, limit int) ([]*storages.Booking, error) {
	bookings, err := s.Store.RetrieveBookings(ctx, userID, createdDate, page, limit)
	if err != nil {
		logger.TIKIErrorf(ctx, "booking storage failed to retrieve bookings of user_id %s created_date %s: %v", userID, createdDate, err)
		return nil, errors.New(dictionary.FailedGetRetrieveBookings)
	}

	return bookings, nil
}

func (s *Booking) Add(ctx context.Context, userID string, booking *storages.Booking) (*storages.Booking, error) {
	now := utils.GetTimeNowWithDefaultLayoutInString()

	booking.UserID = userID
	booking.ID = s.GeneratorUUIDFn()
	booking.CreatedDate = now

	if err := s.Store.AddBooking(ctx, booking); err != nil {
		logger.TIKIErrorln(ctx, err)
		return nil, errors.New(dictionary.StoreBookingFailed)
	}

	return booking, nil
}
