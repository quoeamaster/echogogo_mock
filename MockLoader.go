package main

// this struct is only accessible within this module / plugin
// (not to be exposed to the outside world.
type MockLoader struct {

}

func (m *MockLoader) Load(contentInBytes []byte) (err error) {
	return nil
}
