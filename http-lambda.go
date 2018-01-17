package main

import (
	"context"
	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	logServerStarted  = "http-lambda server starting on address [%s] and port [%s]"
	logFunctionCalled = "Function %s called"
	logSignalReceived = "Signal [%s] received, shutting down server"

	urlFunctionParameter = "function-name"
	urlFunction          = "/function/{" + urlFunctionParameter + "}"

	cfgListenAddress = ""
	cfgListenPort    = "18080"
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
	app := kingpin.New("http-lambda", "Simple golang-based utility which enables AWS Lambda functions to be invoked from an HTTP endpoint")
	appHost := app.Flag("hostname", "hostname to bind to").Short('h').Default(cfgListenAddress).String()
	appPort := app.Flag("port", "port to bind to").Short('p').Default(cfgListenPort).String()
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger.Infof(logServerStarted, *appHost, *appPort)

	router := mux.NewRouter()
	router.HandleFunc(urlFunction, lambdaRequest).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    *appHost + ":" + *appPort,
		Handler: handlers.LoggingHandler(os.Stdout, router),
	}

	// See https://en.wikipedia.org/wiki/Signal_(IPC)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-signalChannel

		logger.Infof(logSignalReceived, s)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)

	}()

	logger.Info(server.ListenAndServe())
}
