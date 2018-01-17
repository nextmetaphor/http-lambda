# HTTP Lambda #
A simple golang-based utility which enables AWS Lambda functions to be invoked from an HTTP endpoint. Use if you want to expose Lambda functions over HTTP without using the AWS API Gateway.

## Getting Started

### Prerequisites
* Local [golang](https://golang.org/) installation; see [https://nextmetaphor.io/2016/12/09/getting-started-with-golang-on-macos/](https://nextmetaphor.io/2016/12/09/getting-started-with-golang-on-macos/) for details on how to install on macOS
* [`gvt` package manager](https://github.com/FiloSottile/gvt) for external vendor dependencies. Install with `go get -u github.com/FiloSottile/gvt` 

### Install

#### Building the Code
First restore the vendor dependencies:
```
gvt restore
```

Alternatively, manually install the vendor dependencies:
```bash
gvt fetch --revision v1.12.59 github.com/aws/aws-sdk-go
gvt fetch --revision v1.0.0 github.com/sirupsen/logrus
gvt fetch --revision v1.6.0 github.com/gorilla/mux
gvt fetch --revision v1.3.0 github.com/gorilla/handlers
gvt fetch --revision v2.2.6 github.com/alecthomas/kingpin
```

Then simply build the binary:
```bash
go build -i
```

## Deployment

### Running The http-lambda Server
Invoke the built server as follows; logs are output to `stderr`, access logs to `stdout`. The server listens binds to address `localhost` on port `18080`. 
```bash
./http-lambda 1>>http-lambda-access.log 2>>http-lambda.log
```

## Validation

### Testing the http
Simply use `cURL` to invoke a lambda function as follows. 

```
curl -X POST http://localhost:18080/function/myLambdaFunction -d '{"key1":"key1-value","key2":"key2-value"}' 
```

In the example above, lambda function `myLambdaFunction` is being called with an input of the following `POST` body.
```json
{
  "key1": "key1-value",
  "key2": "key2-value"
}
```

## Licence ##
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.