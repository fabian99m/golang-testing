package usecase

import (
	"dbtest/domain/dto"
	"dbtest/model"
	"dbtest/model/mocks"

	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	. "github.com/ovechkin-dm/mockio/mock"
)

var mockRepo *mocks.HeroDbInteractor
var heroUsecase model.HeroUseCase
var HeroDbInteractorMock model.HeroDbInteractor

func TestMain(m *testing.M) {
	heroUsecase = NewHeroUseCase(HeroDbInteractorMock)
	runTests := m.Run()
	os.Exit(runTests)
}

func TestGetHeroById(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name     string
		id       int
		hero     model.Hero
		rows     int64
		expected dto.ResponseDto
	}{
		{
			name: "TestGetHeroById Ok",
			id:   1,
			hero: model.Hero{
				Id:         2,
				Name:       "test",
				CreateDate: time.Time{},
			},
			rows:     int64(1),
			expected: dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: dto.HeroDto{}},
		},
		{
			name:     "TestGetHeroById Not Found",
			id:       111,
			hero:     model.Hero{Id: 2, Name: "test", CreateDate: time.Time{}},
			rows:     int64(0),
			expected: dto.ResponseDto{Status: http.StatusNotFound, Codigo: "1003", Mensaje: "Consulta exitosa", Data: dto.HeroDto{}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			SetUp(testCase)
			mockint := Any[int]()
			HeroDbInteractorMock = Mock[model.HeroDbInteractor]()
			heroUsecase = NewHeroUseCase(HeroDbInteractorMock)
			When(HeroDbInteractorMock.GetById(mockint)).ThenReturn(tc.hero, tc.rows)

			response := heroUsecase.GetHeroById(tc.id)

			assert.NotNil(testCase, response)
			assert.Equal(testCase, tc.expected.Status, response.Status, "Codigos deben ser iguales...")
		})
	}
}

func TestGetAll(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name     string
		heros    []model.Hero
		rows     int64
		expected dto.ResponseDto
	}{
		{
			name:     "TestGetAll Ok",
			heros:    []model.Hero{{Id: 1, Name: "fa", CreateDate: time.Now()}, {Id: 2, Name: "test", CreateDate: time.Time{}}},
			rows:     int64(1),
			expected: dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa"},
		},
		{
			name:     "TestGetAll Empty",
			heros:    []model.Hero{},
			rows:     int64(1),
			expected: dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa"},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			mockRepo.On("GetAll", mock.Anything).Return(tc.heros, tc.rows).Once()
			res := heroUsecase.GetAllHeros()

			assert.NotNil(testCase, res.Data.([]dto.HeroDto))
			assert.Equal(testCase, len(res.Data.([]dto.HeroDto)), len(tc.heros))
			assert.Equal(testCase, tc.expected.Status, res.Status, "Codigos deben ser iguales...")
		})
	}
}

func TestSaveHero(test *testing.T) {
	test.Parallel()
	testCases := []struct {
		name     string
		dto      dto.HeroDto
		rows     int64
		expected dto.ResponseDto
	}{
		{
			name:     "TestSaveHero Ok",
			dto:      dto.HeroDto{Name: "test", CreateDate: "2022-12-12"},
			rows:     int64(1),
			expected: dto.ResponseDto{Status: http.StatusCreated, Codigo: "1000", Mensaje: "Consulta exitosa", Data: []dto.HeroDto{}},
		},
		{
			name:     "TestSaveHero ErrorGeneric",
			dto:      dto.HeroDto{Name: "test", CreateDate: "2022-12-12"},
			rows:     int64(0),
			expected: dto.ResponseDto{Status: http.StatusInternalServerError, Codigo: "1002", Mensaje: "Error", Data: []dto.HeroDto{}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			mockRepo.On("Save", mock.Anything).Return(tc.rows).Once()
			res := heroUsecase.SaveHero(tc.dto)

			assert.NotNil(testCase, res)
			assert.Equal(testCase, tc.expected.Status, res.Status, "Codigos deben ser iguales...")
		})
	}
}
