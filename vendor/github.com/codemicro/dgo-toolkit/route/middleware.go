package route

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
)

type MiddlewareTrigger uint

type Middleware struct {
	Name string
	Run     func(interface{}) error
	Trigger MiddlewareTrigger
}

const (
	MiddlewareTriggerValid MiddlewareTrigger = 1 << iota
	MiddlewareTriggerInvalid
	MiddlewareTriggerReactionAdd
	MiddlewareTriggerReactionRemove

	MiddlewareTriggerAllMessages = MiddlewareTriggerInvalid | MiddlewareTriggerValid
	MiddlewareTriggerAllReactions = MiddlewareTriggerReactionAdd | MiddlewareTriggerReactionRemove

	MiddlewareTriggerEverything = MiddlewareTriggerAllMessages | MiddlewareTriggerAllReactions
)

func (b *Kit) runMiddlewares(trigger MiddlewareTrigger, context interface{}) error {
	var combinedErrors error
	for _, x := range b.middlewareSet {
		if x.Trigger & trigger != 0 {
			if err := x.Run(context); err != nil {
				combinedErrors = multierror.Append(combinedErrors, fmt.Errorf("%s: %s", x.Name, err.Error()))
			}
		}
	}
	return combinedErrors
}
