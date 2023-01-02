package repository

import (
	"dbtest/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	repository model.HeroDbInteractor
	mockSql    sqlmock.Sqlmock
}

func (suite *TestSuite) SetupTest() {
	var dbMockConnection *gorm.DB
	dbMockConnection, suite.mockSql = setUpMockDb()
	suite.repository = NewHeroRespository(dbMockConnection)
}

func TestSuiteMain(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetAll() {
	t := suite.T()

	testCases := []struct {
		name     string
		dest     []model.Hero
		sqlRegex string
		res      []model.Hero
	}{
		{
			name:     "TestGetAll OK",
			dest:     *new([]model.Hero),
			sqlRegex: "SELECT(.*)",
			res:      []model.Hero{{Id: 12, Name: "testname1", CreateDate: time.Now()}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(_ *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "name", "create_date"})
			addRowsResult(rows, tc.res)
			suite.mockSql.ExpectQuery(tc.sqlRegex).WillReturnRows(rows)

			heros, _ := suite.repository.GetAll()

			assert.NoError(t, suite.mockSql.ExpectationsWereMet())
			assert.NotNil(t, tc.res)
			assert.True(t, len(heros) > 0)
		})
	}
}

func (suite *TestSuite) TestGetById() {
	t := suite.T()

	testCases := []struct {
		name     string
		dest     model.Hero
		id       int
		sqlRegex string
		res      []model.Hero
	}{
		{
			name:     "TestGetById OK",
			dest:     *new(model.Hero),
			id:       1,
			sqlRegex: "SELECT(.*)",
			res:      []model.Hero{{Id: 12, Name: "testname1", CreateDate: time.Now()}},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(_ *testing.T) {
			rows := sqlmock.NewRows([]string{"id", "name", "create_date"})
			addRowsResult(rows, tc.res)

			suite.mockSql.ExpectQuery(tc.sqlRegex).WillReturnRows(rows)

			hero, rowsAffected := suite.repository.GetById(tc.id)

			assert.NoError(t, suite.mockSql.ExpectationsWereMet())
			assert.NotNil(t, hero)
			assert.Equal(t, int64(1), rowsAffected)
		})
	}
}

func (suite *TestSuite) TestSave() {
	t := suite.T()

	testCases := []struct {
		name     string
		dest     model.Hero
		sqlRegex string
	}{
		{
			name:     "TestSave OK",
			dest:     model.Hero{Name: "testName", CreateDate: time.Now()},
			sqlRegex: "INSERT INTO(.*)",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(_ *testing.T) {
			suite.mockSql.ExpectBegin()
			suite.mockSql.ExpectQuery(tc.sqlRegex).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))
			suite.mockSql.ExpectCommit()

			rowsAffected := suite.repository.Save(&tc.dest)

			assert.NoError(t, suite.mockSql.ExpectationsWereMet())
			assert.Equal(t, int64(1), rowsAffected)
		})
	}
}

func addRowsResult(rows *sqlmock.Rows, dest []model.Hero) {
	for _, v := range dest {
		rows.AddRow(v.Id, v.Name, v.CreateDate)
	}
}

func setUpMockDb() (*gorm.DB, sqlmock.Sqlmock) {
	dbMockConnection, mockSql, _ := sqlmock.New()

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlMock",
		DriverName:           "postgres",
		Conn:                 dbMockConnection,
		PreferSimpleProtocol: true,
	})

	gdb, _ := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{LogLevel: logger.Info},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	return gdb, mockSql
}
