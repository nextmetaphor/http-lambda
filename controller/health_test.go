package controller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sirupsen/logrus/hooks/test"
	"net/http/httptest"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func TestHealthRequest(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.Level = logrus.DebugLevel

	ctx := Context{
		Log:           logger,
		Router:        mux.NewRouter(),
	}

	t.Run("ValidHealth", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, urlHealth, nil)

		ctx.RegisterHealthHandlers()
		ctx.Router.ServeHTTP(writer, request)

		// check HTTP return code and body are correct
		assert.Equal(t, http.StatusOK, writer.Code, "incorrect status code returned")
		assert.True(t, (writer.Body == nil) || (writer.Body.Len() == 0), "expected empty body")

		// check correct message was logged
		assert.Equal(t, 1, len(hook.Entries), "expected 1 message to be logged")
		assert.Equal(t, logHealthCalled, hook.LastEntry().Message, "incorrect message logged")
		assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level, "expected debug message to be logged")
	})
}
