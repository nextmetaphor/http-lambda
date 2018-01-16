package main

import (
	"github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gorilla/handlers"
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"io/ioutil"
)

const (
	logServerStarted  = "http-lambda server starting on address [%s] and port [%d]"
	logFunctionCalled = "Function %s called"
	logSIGTERMReceived = "syscall.SIGTERM received, shutting down server"
	logSignalReceived = "Signal [%s] received, shutting down server"

	urlFunctionParameter = "function-name"
	urlFunction          = "/function/{" + urlFunctionParameter + "}"

	// TODO these should be read from command-line
	cfgListenAddress = ""
	cfgListenPort = 18080
)

var (
	logger = logrus.New()
)

func lambdaRequest(writer http.ResponseWriter, request *http.Request) {
	// refer to
	// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/lambda/aws-go-sdk-lambda-example-run-function.go

	functionName := mux.Vars(request)[urlFunctionParameter]

	logger.Infof(logFunctionCalled, functionName)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String(*sess.Config.Region)})

	payload, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(functionName), Payload: payload})
	if err != nil {
		logger.Error(err)
	}
	if (result == nil) || (result.StatusCode == nil) {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(int(*result.StatusCode))
		writer.Write(result.Payload)
	}
}

func main() {
	logger.Infof(logServerStarted, cfgListenAddress, cfgListenPort)

	router := mux.NewRouter()
	router.HandleFunc(urlFunction, lambdaRequest).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    cfgListenAddress + ":" + strconv.Itoa(cfgListenPort),
		Handler: handlers.LoggingHandler(os.Stdout, router),
	}

	// See https://en.wikipedia.org/wiki/Signal_(IPC)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-signalChannel
		switch s {
		case syscall.SIGTERM:
			logger.Info(logSIGTERMReceived)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)

		default:
			logger.Infof(logSignalReceived, s)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		}

	}()

	logger.Info(server.ListenAndServe())
}
