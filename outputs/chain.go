package outputs

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/nherson/brewski/measurement"
)

// ChainCallback can be used to include multiple sub-callback
// implementations into a single callback handler
type ChainCallback struct {
	callbacks []Callback
}

// NewChainCallback returns an empty ChainCallback
func NewChainCallback() *ChainCallback {
	return &ChainCallback{
		callbacks: make([]Callback, 0),
	}
}

// RegisterCallback adds a callback function to the list of callbacks in the chain
func (cc *ChainCallback) RegisterCallback(scb Callback) {
	cc.callbacks = append(cc.callbacks, scb)
}

// Handle iterates through all registered callbacks and executes them
func (cc *ChainCallback) Handle(s measurement.Sample) error {
	var errList *multierror.Error
	for _, cb := range cc.callbacks {
		err := cb.Handle(s)
		if err != nil {
			multierror.Append(errList, err)
		}
	}
	return errList.ErrorOrNil()
}
