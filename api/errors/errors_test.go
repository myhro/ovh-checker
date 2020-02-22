package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ErrorsTestSuite struct {
	suite.Suite
}

func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

func (s *ErrorsTestSuite) TestValidationMessage() {
	table := []struct {
		in  error
		out string
	}{
		{
			in:  errors.New("strconv.ParseInt: parsing \"xyz\": invalid syntax"),
			out: "parsing \"xyz\": invalid syntax",
		},
		{
			in:  errors.New("Key: 'OffersRequest.Last' Error:Field validation for 'Last' failed on the 'required' tag"),
			out: "field validation for 'last' failed on the 'required' tag",
		},
		{
			in:  errors.New("Key: 'OffersRequest.Country' Error:Field validation for 'Country' failed on the 'required' tag\nKey: 'OffersRequest.Last' Error:Field validation for 'Last' failed on the 'required' tag"),
			out: "field validation for 'country' failed on the 'required' tag",
		},
		{
			in:  errors.New("invalid character '}' looking for beginning of object key string"),
			out: "invalid json",
		},
		{
			in:  errors.New("EOF"),
			out: "invalid json",
		},
		{
			in:  errors.New("not mapped failure"),
			out: "unknown error",
		},
	}

	for _, tt := range table {
		assert.Equal(s.T(), tt.out, ValidationMessage(tt.in))
	}
}
