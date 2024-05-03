package fileusecasetest

import (
	"bytes"
	"dbtest/domain/dto"
	"dbtest/model"
	"dbtest/usecase"

	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	. "github.com/ovechkin-dm/mockio/mock"
	assert "github.com/stretchr/testify/assert"
)

var underTestFile model.FileUseCase
var mockS3 usecase.S3API

func TestMain(m *testing.M) {
	runTests := m.Run()
	os.Exit(runTests)
}

func TestSaveFile(test *testing.T) {
	//test.Parallel()

	var content = []byte("hellooo")
	testCases := []struct {
		name          string
		content       *[]byte
		fileName      string
		output        *s3.PutObjectOutput
		errorExpected error
		emptyResponse bool
	}{
		{
			name:          "TestSaveFile Ok",
			content:       &content,
			fileName:      "file.txt",
			errorExpected: nil,
			output:        &s3.PutObjectOutput{},
			emptyResponse: false,
		},
		{
			name:          "TestSaveFile Error",
			content:       &content,
			fileName:      "file.txt",
			errorExpected: errors.New("Err"),
			output:        &s3.PutObjectOutput{},
			emptyResponse: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			prepareFileTest(testCase)
			When(mockS3.PutObject(Any[context.Context](), Any[*s3.PutObjectInput]())).ThenReturn(tc.output, tc.errorExpected)

			responseFileS3Name, error := underTestFile.SaveFile(tc.content, &tc.fileName)

			var flag bool
			if !tc.emptyResponse {
				flag = strings.Contains(responseFileS3Name, tc.fileName)
			} else {
				flag = responseFileS3Name == ""
			}
			assert.Equal(testCase, tc.errorExpected, error)
			assert.True(testCase, flag)
		})
	}
}

func TestGetFile(test *testing.T) {
	//test.Parallel()

	var data = io.NopCloser(bytes.NewReader([]byte("hellooo")))
	var textPlain = "text/plain"

	testCases := []struct {
		name          string
		key           string
		output        *s3.GetObjectOutput
		errorOutPut   error
		errorExpected error
		expected      *dto.FileResponseDto
	}{
		{
			name: "TestGetFile OK",
			key:  "abc",
			output: &s3.GetObjectOutput{
				ContentType:   &textPlain,
				ContentLength: aws.Int64(int64(len(textPlain))),
				Body:          data,
			},
			errorOutPut:   nil,
			errorExpected: nil,
			expected:      &dto.FileResponseDto{ContentType: textPlain, ContentLength: int64(len(textPlain)), Data: data},
		},
		{
			name:          "TestGetFile NoSuchKey",
			key:           "abc",
			output:        nil,
			errorOutPut:   &types.NoSuchKey{Message: aws.String("the key dont exists")},
			errorExpected: nil,
			expected:      &dto.FileResponseDto{},
		},
		{
			name:          "TestGetFile Error",
			key:           "abc",
			output:        nil,
			errorOutPut:   errors.New("Error en s3..."),
			errorExpected: errors.New("Error en s3..."),
			expected:      nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			prepareFileTest(testCase)
			When(mockS3.GetObject(Any[context.Context](), Any[*s3.GetObjectInput]())).ThenReturn(tc.output, tc.errorOutPut)

			response, error := underTestFile.GetFile(tc.key)

			assert.Equal(testCase, tc.errorExpected, error)
			assert.Equal(testCase, tc.expected, response)
		})
	}
}

func TestGetKeys(test *testing.T) {
	//test.Parallel()

	testCases := []struct {
		name          string
		token         string
		output        *s3.ListObjectsV2Output
		errorExpected error
		expected      *dto.FileResponseDto
	}{
		{
			name:          "TestGetKeys OK firtPage",
			token:         "",
			output:        &s3.ListObjectsV2Output{NextContinuationToken: aws.String("abc"), Contents: []types.Object{{Key: aws.String("abc")}}},
			errorExpected: nil,
			expected:      &dto.FileResponseDto{Next: "abc", Keys: []string{"abc"}},
		},
		{
			name:          "TestGetKeys OK lastPage",
			token:         "bnbn",
			output:        &s3.ListObjectsV2Output{Contents: []types.Object{{Key: aws.String("abc")}}},
			errorExpected: nil,
			expected:      &dto.FileResponseDto{Next: "", Keys: []string{"abc"}},
		},
		{
			name:          "TestGetKeys Error",
			token:         "",
			output:        nil,
			errorExpected: errors.New("Error en s3..."),
			expected:      nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		test.Run(tc.name, func(testCase *testing.T) {
			prepareFileTest(testCase)
			When(mockS3.ListObjectsV2(Any[context.Context](), Any[*s3.ListObjectsV2Input]())).ThenReturn(tc.output, tc.errorExpected)

			response, error := underTestFile.GetFiles(tc.token)

			assert.Equal(testCase, tc.errorExpected, error)
			assert.Equal(testCase, tc.expected, response)
		})
	}
}

func prepareFileTest(t *testing.T) {
	SetUp(t)
	mockS3 = Mock[usecase.S3API]()
	underTestFile = usecase.NewFileUseCase("testbucketname", mockS3)
}
