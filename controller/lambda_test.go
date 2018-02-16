package controller

import (
	"testing"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/url"
)

func TestGetFunctionName(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.Level = logrus.DebugLevel

	ctx := Context{
		Log:           logger,
		Router:        mux.NewRouter(),
	}

	t.Run("SimpleFunctionName", func(t *testing.T) {
		functionName := "myTestFunction"
		fn, fnErr := ctx.getFunctionName(map[string]string {urlFunctionParameter: functionName})

		assert.Equal(t, functionName, fn)
		assert.Nil(t, fnErr)
	})

	t.Run("MissingFunctionName", func(t *testing.T) {
		fn, fnErr := ctx.getFunctionName(nil)

		assert.Equal(t, "", fn)
		assert.Nil(t, fnErr)
	})

	t.Run("ValidEncodedFunctionName", func(t *testing.T) {
		encodedFunctionName := "arn%3Aaws%3Aiam%3A%3Aaccount-id%3Arole%2Flambda_function"
		unencodedFunctionName := "arn:aws:iam::account-id:role/lambda_function"
		fn, fnErr := ctx.getFunctionName(map[string]string {urlFunctionParameter: encodedFunctionName})

		assert.Equal(t, unencodedFunctionName, fn)
		assert.Nil(t, fnErr)
	})

	t.Run("InalidEncodedFunctionName", func(t *testing.T) {
		encodedFunctionName := "arn%EHaws%3Aiam%3A%3Aaccount-id%3Arole%2Flambda_function"
		unencodedFunctionName := ""
		errorString := url.EscapeError("%EH")
		fn, fnErr := ctx.getFunctionName(map[string]string {urlFunctionParameter: encodedFunctionName})

		assert.Equal(t, unencodedFunctionName, fn)
		assert.Equal(t, fnErr, errorString)
	})
}