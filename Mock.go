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
	"github.com/buger/jsonparser"
	"github.com/quoeamaster/echogogo_plugin"
	"io/ioutil"
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


/* ====================================== */
/* =	DoAction - request handling		= */
/* ====================================== */


var randomMessagesMock = [3]string{
	"Life is soooo Good.",
	"With great power comes great responsibility",
	"I shall shed my light over dark evil",
}

var mockInstructionsLoader mockLoader

/* (m *MockModule)  */
func DoAction(request http.Request, endPoint string, optionalMap ...map[string]interface{}) interface{}  {
	modelPtr := new(MockModuleModel)
	method, isJsonResponse := extractPathParamMock(echogogo.ExtractPathParameterFromUrl(request.URL.Path, endPoint))

	switch method {
	case "/configMockEndPoints":
		err := loadMockConfigMock(request)
		if err != nil {
			return err
		} else {
			// prepare the responseBody message
			prepareMockModuleModel(modelPtr, method, "mock instructions loaded", isJsonResponse, time.Now())
		}
	default:
		// add logics to check whether a mock instruction is available or not...
		if len(mockInstructionsLoader.mockInstructionsMap) > 0 {
			// looking for mock instructions
			methodNameWOSlash := method[1:]
			mockInModel := mockInstructionsLoader.GetMockInstructionByMethodNVerb(methodNameWOSlash, request.Method)
			// check validity... example an empty struct is NOT valid...
			if mockInModel.Method == "" && len(mockInModel.Conditions) == 0 {
				prepareMockModuleModel(modelPtr, method, "no such mock API to simulate~", isJsonResponse, time.Now())
			} else {
				bodyMsg, err := getMockResultMock(&mockInModel, isJsonResponse, request)
				if err != nil {
					return err
				}
				prepareMockModuleModel(modelPtr, method, bodyMsg, isJsonResponse, time.Now())
			}
		} else {
			prepareMockModuleModel(modelPtr, method, randomMessageGeneratorMock(), isJsonResponse, time.Now())
		}
	}
	return *modelPtr
}

// verify the conditions and check if any mocked result should be returned instead
func getMockResultMock(model *mockInstructionModel, isJsonResponse bool, request http.Request) (result string, err error) {
	contentInBytes, err := ioutil.ReadAll(request.Body)

	for _, cond := range model.Conditions {
		paramsMatched := false

		// either 0 or 1 set of params ONLY
		if len(cond.Params) == 0 {
			paramsMatched = true
			// return the response, assume additional params are ignored
			if isJsonResponse {
				result = cond.ReturnJson
			} else {
				result = cond.ReturnXml
			}
			break

		} else {
			paramsMap := cond.Params[0]
			// assume all matched and set to false when unmatch case occurs
			paramsMatched = true
			for param, paramVal := range paramsMap {
				valBytes, _, _, err1 := jsonparser.Get(contentInBytes, param)
				if err1 != nil {
					if strings.Index(err1.Error(), "Key path not found") != -1 {
						paramsMatched = false
						break
					} else {
						err = err1
						return
					}
				}	// end -- if (err1 valid, check if it was the case of unmatched key -> continue)
				if paramVal != string(valBytes) {
					paramsMatched = false
					break
				}
			}	// end -- for (per param match check)
			if paramsMatched {
				if isJsonResponse {
					result = cond.ReturnJson
				} else {
					result = cond.ReturnXml
				}
				break
			} else {
				result = "mock api found, but non-matchable params found, hence no results~"
			}	// end -- if (paramsMatch - all match scenario)
		}
	}	// end -- for (conditions)
	return
}

// DoAction routing method - handles /configMockEndPoints
func loadMockConfigMock(request http.Request) (err error) {
	contentInBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	// try to load and parse it...
	mockInstructionsLoader = *newMockLoader(contentInBytes)
	err = mockInstructionsLoader.Load(nil)

	return
}

func prepareMockModuleModel(modelPtr *MockModuleModel, method, responseBody string, isJsonResponse bool, timestamp time.Time) *MockModuleModel {
	modelPtr.MockMethodName = method
	if timestamp.Nanosecond() > 0 {
		modelPtr.Timestamp = timestamp.UTC()
		modelPtr.TimestampEpoch = modelPtr.Timestamp.UnixNano()
	}
	if responseBody != "" {
		modelPtr.ResponseBody = responseBody
	}
	modelPtr.isJsonResponse = isJsonResponse

	return modelPtr
}

// default "mock" result if no other mocking instructions are available
func randomMessageGeneratorMock() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := random.Intn(len(randomMessagesMock))
	return randomMessagesMock[idx]
}
