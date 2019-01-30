
package basicTest

import (
	"github.com/DATA-DOG/godog"
	mockModule "github.com/quoeamaster/echogogo_mock"
)

type testModel struct {
	filename string
	loader mockModule.MockLoader
}

func foundConfigFile(file string) error {
	return godog.ErrPending
}

func loadNparse(numMethods int) error {
	return godog.ErrPending
}

func getMockMethodConfig(methodName, verb string) error {
	return godog.ErrPending
}

func verifyMockMethodConfig(jsonOrXml string, id int, firstName string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a config file named "([^"]*)"$`, foundConfigFile)
	s.Step(`^successfully loaded and parsed, (\d+) methods are available$`, loadNparse)
	s.Step(`^get back the configured method "([^"]*)" having http verb "([^"]*)" should return a method definition$`, getMockMethodConfig)
	s.Step(`^the object returned should consists of a "([^"]*)" response stating id => (\d+) would get back an author named "([^"]*)"$`, verifyMockMethodConfig)
}