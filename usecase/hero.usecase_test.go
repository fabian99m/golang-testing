package usecase

import (
	"net/http"
	"os"
	"testing"

	"dbtest/domain/dto"
	"dbtest/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockRepo *mockRepository
var heroUsecase model.HeroUseCase

func TestMain(m *testing.M) {
	mockRepo = new(mockRepository)
	heroUsecase = NewHeroUseCase(mockRepo)

	runTests := m.Run()
	os.Exit(runTests)
}

func TestGetHeroById(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name     string
		id       int
		rows     int64
		expected dto.ResponseDto
	}{
		{
			name:     "TestGetHeroById Ok",
			id:       1,
			rows:     int64(1),
			expected: dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: dto.HeroDto{}},
		},
		{
			name:     "TestGetHeroById Not Found",
			id:       111,
			rows:     int64(0),
			expected: dto.ResponseDto{Status: http.StatusNotFound, Codigo: "1003", Mensaje: "Consulta exitosa", Data: dto.HeroDto{}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			mockRepo.On("GetById", tc.id, mock.Anything).Return(tc.rows).Once()
			res := heroUsecase.GetHeroById(tc.id)

			assert.NotNil(testCase, res)
			assert.Equal(testCase, tc.expected.Status, res.Status, "Codigos deben ser iguales...")
		})
	}
}

func TestGetAll(test *testing.T) {
	test.Parallel()

	testCases := []struct {
		name     string
		rows     int64
		expected dto.ResponseDto
	}{
		{
			name:     "TestGetAll Ok",
			rows:     int64(1),
			expected: dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: []dto.HeroDto{}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			mockRepo.On("GetAll", mock.Anything).Return(tc.rows).Once()
			res := heroUsecase.GetAllHeros()

			assert.NotNil(testCase, res)
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
		{
			name:     "TestSaveHero ErrorFecha",
			dto:      dto.HeroDto{Name: "test", CreateDate: "2022/12-12"},
			rows:     int64(1),
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

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) GetById(id int, dest *model.Hero) int64 {
	args := m.Called(id, dest)
	return args.Get(0).(int64)
}

func (m *mockRepository) GetAll(dest *[]model.Hero) {
	m.Called(dest)
}

func (m *mockRepository) Save(dest *model.Hero) int64 {
	args := m.Called(dest)
	return args.Get(0).(int64)
}
