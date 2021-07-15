package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"tiki/cmd/api/config"
	BookingHandler "tiki/internal/api/booking/handler"
	bookingMySQL "tiki/internal/api/booking/storages/mysql"
	bookingUseCase "tiki/internal/api/booking/usecase"
	"tiki/internal/api/middleware"
	UserHandler "tiki/internal/api/user/handler"
	userMySQL "tiki/internal/api/user/storages/mysql"
	userUseCase "tiki/internal/api/user/usecase"
	"tiki/internal/api/utils"
	"tiki/internal/pkg/token/jwt"
)

// ToDoHandler implement HTTP server
type ToDoHandler struct {
	UserTp    *UserHandler.User
	BookingTp *BookingHandler.Booking
}

// CreateAPIEngine creates engine instance that serves API endpoints,
func CreateAPIEngine(cfg *config.Config) (*http.Server, error) {
	userTp, bookingTp, err := CreateHandler(cfg)
	if err != nil {
		return nil, err
	}

	generator := createJWTGenerator(cfg)
	userTp.TokenGenerator = generator

	handler := ToDoHandler{
		UserTp:    userTp,
		BookingTp: bookingTp,
	}

	apiDomainString := fmt.Sprintf("%v:%v", cfg.RestfulAPI.Host, cfg.RestfulAPI.Port)
	server := &http.Server{Addr: apiDomainString, Handler: middleware.AddCors(middleware.ValidToken(&handler, cfg, generator))}
	return server, nil
}

func initDBDB(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)
	return sql.Open("mysql", connectionString)
}

func CreateHandler(cfg *config.Config) (*UserHandler.User, *BookingHandler.Booking, error) {

	db, err := initDBDB(cfg)
	userStore := &userMySQL.MySQLDB{DB: db}
	bookingStore := &bookingMySQL.MySQLDB{DB: db}

	if err != nil {
		return nil, nil, err
	}

	userUC := userUseCase.User{
		Store: userStore,
	}
	bookingUC := bookingUseCase.Booking{
		Cfg:          cfg,
		Store:        bookingStore,
		UserStore:    userStore,
		GetTimeNowFn: utils.GetTimeNow,
	}

	return &UserHandler.User{
			UserUC: userUC,
		}, &BookingHandler.Booking{
			Cfg:       cfg,
			BookingUC: bookingUC,
		}, nil
}

func createJWTGenerator(cfg *config.Config) *jwt.Generator {
	return &jwt.Generator{
		Cfg: cfg,
	}
}

func (s *ToDoHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/login":
		switch req.Method {
		case http.MethodPost:
			s.UserTp.Login(resp, req)
			return
		}
	case "/screen":
		switch req.Method {
		case http.MethodPost:
			s.BookingTp.AddScreen(resp, req)
			return
		}
	case "/booking":
		switch req.Method {
		case http.MethodPost:
			s.BookingTp.Booking(resp, req)
			return
		}
	case "/check":
		switch req.Method {
		case http.MethodPost:
			s.BookingTp.Check(resp, req)
			return
		}
	default:
		http.NotFound(resp, req)
		return
	}
	http.Error(resp, "405 method not allowed", http.StatusMethodNotAllowed)
	return
}
