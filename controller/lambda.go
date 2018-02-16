package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"io/ioutil"
	"net/url"
)

const (
	logFunctionCalled  = "Function %s called"
	logFunctionInvalid = "Function name is invalid: (%s)"

	urlFunction          = "/function/{" + urlFunctionParameter + "}"
	urlFunctionParameter = "function-name"
)

func (ctx Context) RegisterLambdaHandlers() {
	ctx.Router.HandleFunc(urlFunction, ctx.lambdaRequest).Methods(http.MethodPost)
}

func (ctx Context) getFunctionName(requestParameters map[string]string) (functionName string, err error) {
	if requestParameters != nil {
		rawFunctionName := requestParameters[urlFunctionParameter]
		functionName, err = url.PathUnescape(rawFunctionName)
	}

	return functionName, err
}

func (ctx Context) lambdaRequest(writer http.ResponseWriter, request *http.Request) {
	// refer to
	// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/lambda/aws-go-sdk-lambda-example-run-function.go

	fnName, fnErr := ctx.getFunctionName(mux.Vars(request))

	if fnErr != nil {
		ctx.Log.Debugf(logFunctionInvalid, fnErr.Error())

		writer.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.Log.Debugf(logFunctionCalled, fnName)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String(*sess.Config.Region)})

	payload, err := ioutil.ReadAll(request.Body)
	if err != nil {
		ctx.Log.Error(err)
		return
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(fnName), Payload: payload})
	if err != nil {
		ctx.Log.Error(err)
	}
	if (result == nil) || (result.StatusCode == nil) {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(int(*result.StatusCode))
		writer.Write(result.Payload)
	}
}
