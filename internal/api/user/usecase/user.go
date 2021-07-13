package usecase

import (
	"context"
	"errors"
	"tiki/internal/api/dictionary"
	"tiki/internal/api/user/storages"
	"tiki/internal/api/utils"
	"tiki/internal/pkg/logger"
)

type User struct {
	Store storages.Store
}

func (s *User) IsValidate(ctx context.Context, userID, password string) (bool, error) {
	user, err := s.Store.GetUserByID(ctx, userID)
	if err != nil {
		logger.TIKIError(ctx, err)
		return false, errors.New(dictionary.FailedToGetUser)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return false, nil
	}

	return true, nil
}
