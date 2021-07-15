//+build integration

package integration

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"tiki/cmd/api/config"
	"tiki/internal/api"
	"tiki/internal/api/dictionary"
	"tiki/internal/api/utils"
)

const (
	HeaderContentType = "Content-Type"
	JSONContentType   = "application/json"
	Token             = "token"
)

var token string

type serverTestSuite struct {
	suite.Suite
	cfg *config.Config
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, &serverTestSuite{})
}

func (s *serverTestSuite) SetupSuite() {
	// Init configs
	state := "local"
	cfg, err := config.Load(&state)
	s.NoError(err)
	s.cfg = cfg

	var server *http.Server
	server, err = api.CreateAPIEngine(cfg)
	s.NoError(err)
	go func() {
		err = server.ListenAndServe()
		s.NoError(err)
	}()
}

func (s *serverTestSuite) Test1Login() {
	reqStr := `{"user_id": "tester", "password": "example"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/login", s.cfg.RestfulAPI.Host, s.cfg.RestfulAPI.Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(HeaderContentType, JSONContentType)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	responseData := utils.CommonResponse{}
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	err = json.Unmarshal(byteBody, &responseData)
	s.NoError(err)
	token = responseData.Data.(string)
	response.Body.Close()
}

func (s *serverTestSuite) TestCheckAvailable() {
	reqStr := `{"screen_id": 1,"location": {"row": 1,"column": 1}}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/check", s.cfg.RestfulAPI.Host, s.cfg.RestfulAPI.Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(HeaderContentType, JSONContentType)
	cookie := &http.Cookie{
		Name:   Token,
		Value:  token,
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusInternalServerError, response.StatusCode)
	responseData := utils.ErrorCommonResponse{}
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	err = json.Unmarshal(byteBody, &responseData)
	s.NoError(err)
	s.Equal(responseData.ErrorStr, dictionary.SeatBooked)
	response.Body.Close()
}

func (s *serverTestSuite) TestCheckBooking() {
	reqStr := `{"screen_id": 1,"number": 0,"locations": [{"row": 0,"column":0}]}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/booking", s.cfg.RestfulAPI.Host, s.cfg.RestfulAPI.Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(HeaderContentType, JSONContentType)
	cookie := &http.Cookie{
		Name:   Token,
		Value:  token,
		MaxAge: 300,
	}
	req.AddCookie(cookie)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusInternalServerError, response.StatusCode)
	responseData := utils.ErrorCommonResponse{}
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	err = json.Unmarshal(byteBody, &responseData)
	s.NoError(err)
	s.Equal(responseData.ErrorStr, dictionary.SeatBooked)
	response.Body.Close()
}
