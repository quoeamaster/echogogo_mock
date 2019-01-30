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
)

// this struct is only accessible within this module / plugin
// (not to be exposed to the outside world.
type mockLoader struct {
	contentInBytes 		[]byte
	mockInstructionsMap	map[string][]mockInstructionModel
	//mockInstructions 	[]mockInstructionModel
}

type mockInstructionModel struct {
	Method 		string
	Conditions 	[]mockConditionModel
}

type mockConditionModel struct {
	// treat every param value as string first
	Params 		[]map[string]string
	ReturnJson 	string
	ReturnXml 	string
}

// constructor method for mockLoader
func newMockLoader(contentInBytes []byte) (loaderPtr *mockLoader) {
	loaderPtr = new(mockLoader)
	loaderPtr.contentInBytes = contentInBytes

	loaderPtr.mockInstructionsMap = make(map[string][]mockInstructionModel)

	return
}

// constructor method for mockInstructionModel
func newMockInstructionModel() (modelPtr *mockInstructionModel) {
	modelPtr = new(mockInstructionModel)
	modelPtr.Conditions = make([]mockConditionModel, 0)

	return
}

// return a String format to describe this mockLoader instance
func (m *mockLoader) String() string {
	bContent := new(bytes.Buffer)

	bContent.WriteString("mockInstructionsMap => \n")
	for key, value := range m.mockInstructionsMap {
		bContent.WriteString(fmt.Sprintf("%v:\n", key))
		// per value level...
		modelList := []mockInstructionModel(value)
		for _, model := range modelList {
			bContent.WriteString(fmt.Sprintf("    %v\n", model.Method))
			/*
			bContent.WriteString(fmt.Sprintf("    %v - [", model.Method))
			for _, condition := range model.Conditions {
				bContent.WriteString(fmt.Sprintf("\n        { returnJson: %v, \n          returnXml: %v }", condition.ReturnJson, condition.ReturnXml))
				bContent.WriteString(fmt.Sprintf("\n        %v", condition.Params))
			}
			bContent.WriteString(fmt.Sprintf("\n    ]\n\n"))
			*/
		}
	}

	return bContent.String()
}

// load and parse for Mock instructions. If the parameter given is non-nil,
// then the given content would be used for parsing.
func (m *mockLoader) Load(contentInBytes []byte) (err error) {
	err = nil
	rawBytes := m.contentInBytes

	// use this given []byte for parsing if not nil
	if contentInBytes != nil {
		rawBytes = contentInBytes
	}

	// 0. parse GET
	m.mockInstructionsMap["GET"] = make([]mockInstructionModel, 0)
	configsBytes, _, _, err := jsonparser.Get(rawBytes, "GET")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	err = m.parseMockInstructionFromBytes(configsBytes, "GET")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	// 1. parse PUT
	m.mockInstructionsMap["PUT"] = make([]mockInstructionModel, 0)
	configsBytes, _, _, err = jsonparser.Get(rawBytes, "PUT")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	err = m.parseMockInstructionFromBytes(configsBytes, "PUT")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	// 2. parse POST
	m.mockInstructionsMap["POST"] = make([]mockInstructionModel, 0)
	configsBytes, _, _, err = jsonparser.Get(rawBytes, "POST")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	err = m.parseMockInstructionFromBytes(configsBytes, "POST")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	// 3. parse DELETE
	m.mockInstructionsMap["DELETE"] = make([]mockInstructionModel, 0)
	configsBytes, _, _, err = jsonparser.Get(rawBytes, "DELETE")
	m.genericErrHandler(err)
	if err != nil {
		return
	}
	err = m.parseMockInstructionFromBytes(configsBytes, "DELETE")
	m.genericErrHandler(err)
	if err != nil {
		return
	}

	return
}


// method to parse mock instructions from the raw bytes
func (m *mockLoader) parseMockInstructionFromBytes(contentInBytes []byte, verb string) (err error) {
	_, err = jsonparser.ArrayEach(contentInBytes, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		if err1 != nil {
			err = err1
			return
		}
		// create a new mockInstructionModel instance
		modelPtr := newMockInstructionModel()
		// parse "method"
		sVal, err := jsonparser.GetString(value, "method")
		if err != nil {
			return
		}
		modelPtr.Method = sVal
		// parse "conditions"
		bVal, _, _, err := jsonparser.Get(value, "conditions")
		if err != nil {
			return
		}
		err = m.parseMockConditionsFromBytes(bVal, modelPtr)
		if err != nil {
			return
		}
		m.mockInstructionsMap[verb] = append(m.mockInstructionsMap[verb], *modelPtr)
		//m.mockInstructions = append(m.mockInstructions, *modelPtr)
	})
	return
}

// method to parse the conditions[x] contents. Involved parseMockConditionParamsFromBytes() method call(s)
func (m *mockLoader) parseMockConditionsFromBytes(contentInBytes []byte, mockInstructionModelPtr *mockInstructionModel) (err error) {
	_, err = jsonparser.ArrayEach(contentInBytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// create instance of mockConditionModel
		modelPtr := new(mockConditionModel)
		// parse "params" => [] of map
		paramsBytes, _, _, err := jsonparser.Get(value, "params")
		if err != nil {
			return
		}
		err = m.parseMockConditionParamsFromBytes(paramsBytes, modelPtr)
		if err != nil {
			return
		}
		// TODO: parse returnJson and returnXml
		// parse "returnJson"
		jVal, _, _, err := jsonparser.Get(value, "returnJson")
		if err != nil {
			return
		}
		modelPtr.ReturnJson = string(jVal)
		// parse "returnXml"
		xVal, _, _, err := jsonparser.Get(value, "returnXml")
		if err != nil {
			return
		}
		modelPtr.ReturnXml = string(xVal)

		mockInstructionModelPtr.Conditions = append(mockInstructionModelPtr.Conditions, *modelPtr)
	})
	return
}

// method to parse the conditions[x] > params[x] contents
func (m *mockLoader) parseMockConditionParamsFromBytes(contentInBytes []byte, mockConditionModelPtr *mockConditionModel) (err error) {
	mockConditionModelPtr.Params = make([]map[string]string, 0)

	_, err = jsonparser.ArrayEach(contentInBytes, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		if err1 != nil {
			err = err1
			return
		}
		// parse "key" by "key"
		paramsMap := make(map[string]string)
		err = jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			paramsMap[string(key)] = string(value)
			return nil
		})
		if err != nil {
			return
		}
		mockConditionModelPtr.Params = append(mockConditionModelPtr.Params, paramsMap)
	})
	return
}


// method to get back the configuration based on a "method" name and the http "verb"
func (m *mockLoader) GetMockInstructionByMethodNVerb(method string, verb string) (model mockInstructionModel) {
	modelList := m.mockInstructionsMap[verb]
	if modelList != nil && len(modelList) > 0 {
		for _, iModel := range modelList {
			if iModel.Method == method {
				model = iModel
				break
			}
		}
	}
	return
}


// generic error handler, if err != nil => print the error to console
func (m *mockLoader) genericErrHandler(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}