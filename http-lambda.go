package main

import (
	"context"
	"github.com/alecthomas/kingpin"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nextmetaphor/http-lambda/controller"
	"github.com/nextmetaphor/http-lambda/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	logInsecureServerStarting = "http-lambda server starting on address [%s] and port [%s] with a insecure configuration"
	logSecureServerStarting   = "http-lambda server starting on address [%s] and port [%s] with a secure configuration: cert[%s] key[%s]"
	logSignalReceived         = "Signal [%s] received, shutting down server"

	cfgDefaultListenAddress   = ""
	cfgDefaultListenPort      = "18080"
	cfgDefaultLogLevel        = log.LevelWarn
	cfgDefaultCertificateFile = "http-lambda.crt"
	cfgDefaultKeyFile         = "http-lambda.key"
)

func main() {

	app := kingpin.New("http-lambda", "Simple golang-based utility which enables AWS Lambda functions to be invoked from an HTTP endpoint")
	appHost := app.Flag("hostname", "hostname to bind to").Short('h').Default(cfgDefaultListenAddress).String()
	appPort := app.Flag("port", "port to bind to").Short('p').Default(cfgDefaultListenPort).String()
	appCertFile := app.Flag("certFile", "TLS certificate file").Short('c').Default(cfgDefaultCertificateFile).String()
	appKeyFile := app.Flag("keyFile", "TLS key file").Short('k').Default(cfgDefaultKeyFile).String()
	appSecure := app.Flag("secure", "whether to use secure TLS connection").Short('s').Default("false").Bool()
	appLogLevel := app.Flag("logLevel", "log level: debug, info, warn or error").Short('l').Default(cfgDefaultLogLevel).String()
	kingpin.MustParse(app.Parse(os.Args[1:]))

	ctx := controller.Context{Log: log.Get(*appLogLevel), Router: mux.NewRouter()}

	ctx.RegisterHealthHandlers()
	ctx.RegisterLambdaHandlers()

	server := &http.Server{
		Addr:    *appHost + ":" + *appPort,
		Handler: handlers.LoggingHandler(os.Stdout, ctx.Router),
	}

	// See https://en.wikipedia.org/wiki/Signal_(IPC)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-signalChannel

		ctx.Log.Infof(logSignalReceived, s)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)

	}()

	if *appSecure {
		ctx.Log.Infof(logSecureServerStarting, *appHost, *appPort, *appCertFile, *appKeyFile)
		ctx.Log.Info(server.ListenAndServeTLS(*appCertFile, *appKeyFile))
	} else {
		ctx.Log.Infof(logInsecureServerStarting, *appHost, *appPort)
		ctx.Log.Info(server.ListenAndServe())
	}
}
