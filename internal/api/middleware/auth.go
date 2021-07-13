package middleware

import (
	"context"
	"net/http"
	"tiki/cmd/api/config"
	"tiki/internal/pkg/logger"
	"tiki/internal/pkg/token"
)

var (
	loginPath = "/login"
)

const (
	CookieKey = "token"
)

func ValidToken(next http.Handler, cfg *config.Config, generator token.Generator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case loginPath:
			tikiLogger.Infoln("login request, skip validate")
		default:
			tokenStr, err := r.Cookie(CookieKey)
			if err != nil {
				if err == http.ErrNoCookie {
					tikiLogger.Errorf("validate token err: %v", err)
					http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
					return
				}
				// For any other type of error, return a bad request status
				tikiLogger.Errorf("validate token err: %v", err)
				http.Error(w, "401 Unauthorized", http.StatusBadRequest)
				return
			}

			id, err := generator.ValidateToken(tokenStr.Value)
			if err != nil {
				tikiLogger.Errorf("validate token err: %v", err)
				http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
				return
			}

			tikiLogger = logger.WithFields(map[string]interface{}{
				"user_id": id,
				"path":    r.URL.Path,
				"method":  r.Method,
			})

			r = r.WithContext(context.WithValue(r.Context(), logger.TIKILoggerConText, tikiLogger))
			r = r.WithContext(context.WithValue(r.Context(), token.UserIDField, id))
		}

		next.ServeHTTP(w, r)
	})
}
