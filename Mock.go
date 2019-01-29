/*
 * Licensed to Echogogo under one or more contributor
 * license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. Echogogo licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/quoeamaster/echogogo_plugin"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// response model for the Mock module
type MockModuleModel struct {
	ResponseBody string
	MockMethodName string
	Timestamp time.Time
	TimestampEpoch int64

	isJsonResponse bool
}

// return the string representation of the Model
func (m* MockModuleModel) String() string {
	var bBufPtr = new(bytes.Buffer)

	bBufPtr.WriteString("{")

	bBufPtr.WriteString("\"ResponseBody\": \"")
	bBufPtr.WriteString(m.ResponseBody)
	bBufPtr.WriteString("\"")
	bBufPtr.WriteString(",\\n\"MockMethodName\": \"")
	bBufPtr.WriteString(m.MockMethodName)
	bBufPtr.WriteString("\"")
	bBufPtr.WriteString(",\\n\"Timestamp\": \"")
	bBufPtr.WriteString(m.Timestamp.String())
	bBufPtr.WriteString("\"")
	bBufPtr.WriteString(",\\n\"TimestampEpoch\": ")
	bBufPtr.WriteString(fmt.Sprintf("%v", m.TimestampEpoch))

	bBufPtr.WriteString("\\n}")

	return bBufPtr.String()
}


func extractPathParamMock(method string) (methodName string, isJsonResponse bool) {
	// check whether there are "paths"
	methodName = method
	isJsonResponse = true

	if strings.HasPrefix(method, "/json/") {
		idx := len("/json")
		methodName = method[idx:]
	} else if strings.HasPrefix(method, "/xml/") {
		idx := len("/xml")
		methodName = method[idx:]
		isJsonResponse = false
	}
	return
}

/* (m *MockModule)  */
func GetRestConfig() map[string]interface{} {
	/* TODO: either read from a file or simply overwrite it programmatically.... */
	mapModelPtr := make(map[string]interface{})

	mapModelPtr["consumeFormat"] = echogogo.FORMAT_JSON
	mapModelPtr["produceFormat"] = echogogo.FORMAT_XML_JSON
	mapModelPtr["path"] = "/mock"
	mapModelPtr["endPoints"] = []string {
		"GET::/{method}", "GET::/json/{method}", "GET::/xml/{method}",
		"POST::/{method}", "POST::/json/{method}", "POST::/xml/{method}",
		"PUT::/{method}", "POST::/json/{method}", "POST::/xml/{method}",
		"DELETE::/{method}", "POST::/json/{method}", "POST::/xml/{method}",
		"POST::/configMockEndPoints",
	}
	return mapModelPtr
}


/* (m *MockModule)  */
func DoAction(request http.Request, endPoint string, optionalMap ...map[string]interface{}) interface{}  {
	modelPtr := new(MockModuleModel)
	method, isJsonResponse := extractPathParamMock(echogogo.ExtractPathParameterFromUrl(request.URL.Path, endPoint))

	// TODO: add logics to check whether a mock instruction is available or not...
	modelPtr.ResponseBody = randomMessageGeneratorMock()

	// setup the return model
	modelPtr.MockMethodName = method
	modelPtr.Timestamp = time.Now().UTC()
	modelPtr.TimestampEpoch = modelPtr.Timestamp.UnixNano()
	modelPtr.isJsonResponse = isJsonResponse

	return *modelPtr
}


var randomMessagesMock = [3]string{
	"Life is soooo Good.",
	"With great power comes great responsibility",
	"I shall shed my light over dark evil",
}

// default "mock" result if no other mocking instructions are available
func randomMessageGeneratorMock() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := random.Intn(len(randomMessagesMock))
	return randomMessagesMock[idx]
}
