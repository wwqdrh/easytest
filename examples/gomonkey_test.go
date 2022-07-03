package examples

import (
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func DoSomething(name string, args ...string) (string, error) {
	return "", errors.New("TODO")
}

func TestExec(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	outputExpect := "xxx-vethName100-yyy"
	guard := patches.ApplyFunc(DoSomething, func(_ string, _ ...string) (string, error) {
		return outputExpect, nil
	})
	defer guard.Reset()
	output, err := DoSomething("asd", "1", "2", "3")
	assert.Nil(t, err)
	assert.Equal(t, outputExpect, output)
}
