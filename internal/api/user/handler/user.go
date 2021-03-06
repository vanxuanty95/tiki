package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"tiki/internal/api/dictionary"
	"tiki/internal/api/middleware"
	"tiki/internal/api/user/storages"
	userUseCase "tiki/internal/api/user/usecase"
	"tiki/internal/api/utils"
	"tiki/internal/pkg/logger"
	"tiki/internal/pkg/token"
	"time"
)

type User struct {
	UserUC         userUseCase.User
	TokenGenerator token.Generator
}

func (s *User) Login(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	user := &storages.User{}
	err := json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, err)
		return
	}

	valErr := utils.ValidateRequest(user)
	if valErr != nil {
		utils.WriteJSON(ctx, resp, http.StatusBadRequest, nil, valErr)
		return
	}

	isValid, err := s.UserUC.IsValidate(ctx, user.ID, user.Password)
	if err != nil {
		utils.WriteJSON(ctx, resp, http.StatusUnauthorized, nil, err)
		return
	}

	if !isValid {
		utils.WriteJSON(ctx, resp, http.StatusUnauthorized, nil, errors.New(dictionary.IncorrectLogin))
		return
	}

	tokenStr, err := s.TokenGenerator.CreateToken(user.ID)
	if err != nil {
		logger.TIKIPrintf(ctx, "failed to create token: %s", err)
		utils.WriteJSON(ctx, resp, http.StatusInternalServerError, nil, errors.New(dictionary.FailedToCreateToken))
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:    middleware.CookieKey,
		Value:   tokenStr,
		Expires: time.Now().Add(s.TokenGenerator.GetTimeExpire()),
	})

	utils.WriteJSON(ctx, resp, http.StatusOK, tokenStr, nil)
}
