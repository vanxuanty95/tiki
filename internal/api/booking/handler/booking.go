package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"tiki/cmd/api/config"
	"tiki/internal/api/booking/storages"
	"tiki/internal/api/booking/storages/model"
	"tiki/internal/api/booking/usecase"
	"tiki/internal/api/dictionary"
	"tiki/internal/api/utils"
	"tiki/internal/pkg/token"
)

const (
	CreatedDateField = "created_date"
	PageField        = "page"
	Limit            = "limit"
)

type Booking struct {
	Cfg       *config.Config
	BookingUC usecase.Booking
}

func (t *Booking) AddScreen(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	screen := &storages.Screen{}
	err := json.NewDecoder(req.Body).Decode(screen)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, err)
		return
	}

	valErr := utils.ValidateRequest(screen)
	if valErr != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, valErr)
		return
	}

	if screen.NumberSeatRow+1 < t.Cfg.Distance || screen.NumberSeatColumn+1 < t.Cfg.Distance {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, errors.New(dictionary.LesserThanDistance))
		return
	}

	userID, err := utils.GetValueFromCtx(ctx, token.UserIDField)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	id, err := t.BookingUC.CreateNewScreen(ctx, userID, screen)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	utils.WriteJSON(ctx, resp, http.StatusOK, map[string]int{"id": id}, nil)
}

func (t *Booking) Booking(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	bookingRq := &model.BookingRequest{}
	err := json.NewDecoder(req.Body).Decode(bookingRq)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, err)
		return
	}

	valErr := utils.ValidateRequest(bookingRq)
	if valErr != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, valErr)
		return
	}

	ok := utils.ValidateBookingRequest(bookingRq)
	if !ok {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, errors.New(dictionary.FailedValidateBookingRequest))
		return
	}

	userID, err := utils.GetValueFromCtx(ctx, token.UserIDField)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	seats, err := t.BookingUC.BookingNewSeat(ctx, userID, bookingRq)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	if len(seats) == 0 {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, errors.New(dictionary.SeatBooked))
		return
	}

	utils.WriteJSON(ctx, resp, http.StatusOK, seats, nil)
}

func (t *Booking) Check(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	checkRq := &model.CheckAvailableRequest{}
	err := json.NewDecoder(req.Body).Decode(checkRq)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, err)
		return
	}

	valErr := utils.ValidateRequest(checkRq)
	if valErr != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, valErr)
		return
	}

	userID, err := utils.GetValueFromCtx(ctx, token.UserIDField)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	err = t.BookingUC.CheckAvailable(ctx, userID, checkRq)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	utils.WriteJSON(ctx, resp, http.StatusOK, dictionary.OkSeat, nil)
}
