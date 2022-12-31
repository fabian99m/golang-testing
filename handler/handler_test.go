package handler

import (
	"bytes"
	"dbtest/domain/dto"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func SetUpRouter() (*gin.Engine, *gin.RouterGroup) {
	server := gin.Default()
	return server, server.Group("v1")
}

func TestGetAll(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "TestGetAll OK",
			expectedStatus: 200,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			useCaseMock := new(useCaseMockRepository)
			useCaseMock.On("GetAllHeros").Return(dto.ResponseDto{Status: http.StatusOK}).Once()

			server, api := SetUpRouter()

			NewHeroHandler(api, useCaseMock)
			req, _ := http.NewRequest(http.MethodGet, "/v1/heros", nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.NotNil(test, w)
			assert.Equal(test, tc.expectedStatus, w.Code)
		})
	}
}

func TestGetHeroById(test *testing.T) {
	test.Parallel()
	testCases := []struct {
		name           string
		expectedStatus int
		id             string
	}{
		{
			name:           "TestGetHeroById OK",
			expectedStatus: 200,
			id:             "1",
		}, {
			name:           "TestGetHeroById Error Id",
			expectedStatus: 400,
			id:             "as",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		test.Run(tc.name, func(testCase *testing.T) {

			useCaseMock := new(useCaseMockRepository)
			useCaseMock.On("GetHeroById", 1).Return(dto.ResponseDto{Status: http.StatusOK}).Once()

			server, api := SetUpRouter()

			NewHeroHandler(api, useCaseMock)
			req, _ := http.NewRequest(http.MethodGet, "/v1/heros/"+tc.id, nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.NotNil(test, w)
			assert.Equal(test, tc.expectedStatus, w.Code)
		})
	}
}

func TestSaveHero(test *testing.T) {
	test.Parallel()
	testCases := []struct {
		name           string
		expectedStatus int
		dto            []byte
	}{
		{
			name:           "TestSaveHero OK",
			expectedStatus: http.StatusCreated,
			dto:            getJson(dto.HeroDto{Name: "test", CreateDate: "2022-12-12"}),
		},
		{
			name:           "TestSaveHero BadRequest",
			expectedStatus: http.StatusBadRequest,
			dto:            getJson(""),
		},
	}
	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			useCaseMock := new(useCaseMockRepository)
			useCaseMock.On("SaveHero", mock.Anything).Return(dto.ResponseDto{Status: http.StatusCreated}).Once()

			server, api := SetUpRouter()

			NewHeroHandler(api, useCaseMock)
			req, _ := http.NewRequest(http.MethodPost, "/v1/heros", bytes.NewBuffer(tc.dto))

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.NotNil(test, w)
			assert.Equal(test, tc.expectedStatus, w.Code)
		})
	}
}

func getJson(dest any) []byte {
	s, _ := json.Marshal(dest)
	return s
}

type useCaseMockRepository struct {
	mock.Mock
}

func (m *useCaseMockRepository) GetAllHeros() dto.ResponseDto {
	args := m.Called()
	return args.Get(0).(dto.ResponseDto)
}

func (m *useCaseMockRepository) GetHeroById(id int) dto.ResponseDto {
	args := m.Called(id)
	return args.Get(0).(dto.ResponseDto)
}

func (m *useCaseMockRepository) SaveHero(hero dto.HeroDto) dto.ResponseDto {
	args := m.Called(hero)
	return args.Get(0).(dto.ResponseDto)
}
