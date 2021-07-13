//+build integration

package integration

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"tiki/config"
	"tiki/internal/api"
	"tiki/internal/api/dictionary"
	"tiki/internal/api/utils"
)

const (
	HeaderContentType = "Content-Type"
	JSONContentType   = "application/json"
	Authorization     = "Authorization"
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
	reqStr := `{"user_id": "firstUser", "password": "example"}`
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

func (s *serverTestSuite) TestListBookings() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/bookings?created_date=2020-06-29", s.cfg.RestfulAPI.Host, s.cfg.RestfulAPI.Port), nil)
	s.NoError(err)

	req.Header.Set(HeaderContentType, JSONContentType)
	req.Header.Set(Authorization, token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}

func (s *serverTestSuite) TestAddBooking() {
	reqStr := `{"content": "a"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/bookings", s.cfg.RestfulAPI.Host, s.cfg.RestfulAPI.Port), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(HeaderContentType, JSONContentType)
	req.Header.Set(Authorization, token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	if http.StatusOK != response.StatusCode {
		responseData := utils.ErrorCommonResponse{}
		byteBody, err := ioutil.ReadAll(response.Body)
		s.NoError(err)

		err = json.Unmarshal(byteBody, &responseData)
		s.Equal(responseData.ErrorStr, dictionary.UserReachBookingLimit)
		s.NoError(err)
	} else {
		s.Equal(http.StatusOK, response.StatusCode)
	}
	response.Body.Close()
}