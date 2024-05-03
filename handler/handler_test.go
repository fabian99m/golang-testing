package handler

import (
	"dbtest/app"
	"dbtest/domain/dto"
	"dbtest/model"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	. "github.com/ovechkin-dm/mockio/mock"
)

var heroUseCaseMock model.HeroUseCase
var fileUseCaseMock model.FileUseCase

func TestMain(m *testing.M) {
	runTests := m.Run()
	os.Exit(runTests)
}

func TestGetAll(test *testing.T) {
	//test.Parallel()

	testCases := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "TestGetAll OK",
			expectedStatus: 200,
		},
		{
			name:           "TestGetAll Fail",
			expectedStatus: 500,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		server, api := SetUpRouter()

		test.Run(tc.name, func(testCase *testing.T) {
			prepareTest(testCase)
			When(heroUseCaseMock.GetAllHeros()).ThenReturn(dto.ResponseDto{Status: tc.expectedStatus})

			NewHandler(api, &app.AppConfig{HeroUseCase: heroUseCaseMock, FileUseCase: fileUseCaseMock})

			request, _ := http.NewRequest(http.MethodGet, "/v1/heros", nil)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, request)

			assert.NotNil(test, responseRecorder)
			assert.Equal(test, tc.expectedStatus, responseRecorder.Code)
		})
	}
}

func TestGetHeroById(test *testing.T) {
	//test.Parallel()

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
			expectedStatus: 401,
			id:             "aaaaa",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		server, api := SetUpRouter()

		test.Run(tc.name, func(testCase *testing.T) {
			prepareTest(testCase)

			When(heroUseCaseMock.GetHeroById(AnyInt())).ThenReturn(dto.ResponseDto{Status: http.StatusOK})

			NewHandler(api, &app.AppConfig{HeroUseCase: heroUseCaseMock, FileUseCase: fileUseCaseMock})

			request, _ := http.NewRequest(http.MethodGet, "/v1/heros/"+tc.id, nil)
			responseRecorder := httptest.NewRecorder()
			server.ServeHTTP(responseRecorder, request)

			assert.NotNil(test, responseRecorder)
			assert.Equal(test, tc.expectedStatus, responseRecorder.Code)
		})
	}
}

func TestSaveHero(test *testing.T) {
	//test.Parallel()
	
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
			prepareTest(testCase)

			When(heroUseCaseMock.SaveHero(AnyOfType(dto.HeroDto{}))).ThenReturn(dto.ResponseDto{Status: http.StatusCreated})
			server, api := SetUpRouter()

			NewHandler(api, &app.AppConfig{HeroUseCase: heroUseCaseMock, FileUseCase: fileUseCaseMock})
			req, _ := http.NewRequest(http.MethodPost, "/v1/heros", bytes.NewBuffer(tc.dto))

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.NotNil(test, w)
			assert.Equal(test, tc.expectedStatus, w.Code)
		})
	}
}

func SetUpRouter() (*gin.Engine, *gin.RouterGroup) {
	server := gin.Default()
	return server, server.Group("v1")
}

func getJson(dest any) []byte {
	jsonString, _ := json.Marshal(dest)
	return jsonString
}

func prepareTest(t *testing.T) {
	SetUp(t)
	heroUseCaseMock = Mock[model.HeroUseCase]()
	fileUseCaseMock = Mock[model.FileUseCase]()
}