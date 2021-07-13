package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"tiki/internal/api/booking/storages"
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
	BookingUC usecase.Booking
}

func (t *Booking) List(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userID, err := utils.GetValueFromCtx(ctx, token.UserIDField)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	createdDate := req.FormValue(CreatedDateField)
	if createdDate == "" {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, errors.New(dictionary.CreatedDateRequestEmpty))
		return
	}

	if ok := utils.ValidateDateFromString(createdDate, utils.DefaultLayout); !ok {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, errors.New(dictionary.DateRequestEmptyIsNotValid))
		return
	}

	page, limit := 0, 0
	pageStr := req.FormValue(PageField)
	if value, ok := utils.ValidateInputIsInteger(pageStr); pageStr != "" && ok {
		page = value
	}

	limitStr := req.FormValue(Limit)
	if value, ok := utils.ValidateInputIsInteger(limitStr); limitStr != "" && ok {
		limit = value
	}

	bookings, err := t.BookingUC.List(ctx, userID, createdDate, page, limit)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	utils.WriteJSON(ctx, resp, http.StatusOK, bookings, err)
}

func (t *Booking) Add(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	booking := &storages.Booking{}
	err := json.NewDecoder(req.Body).Decode(booking)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, err)
		return
	}

	valErr := utils.ValidateRequest(booking)
	if valErr != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, valErr)
		return
	}

	userID, err := utils.GetValueFromCtx(ctx, token.UserIDField)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	booking, err = t.BookingUC.Add(ctx, userID, booking)
	if err != nil {
		if err.Error() == errors.New(dictionary.UserReachBookingLimit).Error() {
			utils.WriteJSON(ctx, resp, http.StatusForbidden, nil, err)
			return
		}
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, err)
		return
	}

	utils.WriteJSON(ctx, resp, http.StatusOK, booking, nil)
}
