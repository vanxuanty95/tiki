package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strings"
	"tiki/internal/api/booking/storages/model"
)

func ValidateRequest(req interface{}) error {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.Struct(req)
	var result []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			result = append(result, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
		return errors.New(strings.Join(result, ", "))
	}
	return nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateBookingRequest(request *model.BookingRequest) bool {
	if request == nil {
		return false
	}
	if request.Number == 0 && request.Locations == nil {
		return false
	}
	if request.Number != 0 && request.Locations != nil {
		return false
	}
	return true
}
