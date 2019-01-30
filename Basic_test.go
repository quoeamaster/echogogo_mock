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
	t.Run("1", func(t *testing.T) {
		model := loaderPtr.GetMockInstructionByMethodNVerb("getAuthorById", "GET")
		// TODO: verify the results... no assertions in golang... hm
		fmt.Printf("%v\n", model)
	})


	// tear down code
}

