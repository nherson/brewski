package handlers

import (
	"fmt"
	"strings"

	"github.com/nherson/brewski/measurement"
)

// StdoutCallback is a basic printing callback handler
type StdoutCallback struct{}

// NewStdoutCallback returns a new stdout callback...
// it just prints things to stdout.
func NewStdoutCallback() *StdoutCallback {
	return &StdoutCallback{}
}

// Handle is a simple callback that just prints to stdout
func (sc *StdoutCallback) Handle(s measurement.Sample) error {
	stringSnippets := []string{}
	stringSnippets = append(stringSnippets, fmt.Sprintf("device=%s", s.DeviceName()))
	for _, d := range s.Datapoints() {
		stringSnippets = append(stringSnippets, fmt.Sprintf("%s=%f", d.Name(), d.Value()))
	}
	dataString := strings.Join(stringSnippets, ",")
	fmt.Printf("data read! %s\n", dataString)
	return nil
}
