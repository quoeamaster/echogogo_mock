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
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestBasicWorkflow(t *testing.T)  {
	// setup code
	filename := "sampleMockInstructions.json"

	rawBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1000)
	}
	loaderPtr := newMockLoader(rawBytes)
	err = loaderPtr.Load(nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1010)
	}

	// scenario(s)
	t.Run("1. basic workflow testing (load, parse, get mock-instructions)", func(t *testing.T) {
		model := loaderPtr.GetMockInstructionByMethodNVerb("getAuthorById", "GET")
		// verify the results... no assertions in golang... hm
		if !assertMockInstructionModelValid(t, model) {
			t.Errorf("mock instruction model is invalid, should consists of 'method' and 'conditions', got => %v\n", model)
		}
		// verify params level
		paramsArray := model.Conditions[0].Params
		if len(paramsArray) == 0 {
			t.Errorf("should consists of at least 1 set of parameters, got => %v\n", paramsArray)
		}
		for _, param := range paramsArray {
			if param["id"] != "13" && param["id"] != "999" {
				t.Errorf("parameter 'id' should be available, got => %v\n", param["id"])
			}
		}
		fmt.Printf("\t### method (%v) retrieved successfully!\n\n", model.Method)
	})

	// tear down code
}

// assertion method to check if model passed is valid or not
func assertMockInstructionModelValid(t *testing.T, model mockInstructionModel) (valid bool) {
	t.Helper()

	valid = true
	if model.Method == "" {
		valid = false
	} else if len(model.Conditions) == 0 {
		valid = false
	}
	return
}

