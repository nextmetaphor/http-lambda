package controller

import (
	"net/http"
)

const (
	logHealthCalled = "Health check called"
	urlHealth       = "/health"
)

func (ctx Context) healthRequest(writer http.ResponseWriter, request *http.Request) {
	ctx.Log.Debug(logHealthCalled)
	writer.WriteHeader(http.StatusOK)
}

func (ctx Context) RegisterHealthHandlers() {
	ctx.Router.HandleFunc(urlHealth, ctx.healthRequest).Methods(http.MethodGet)
}
