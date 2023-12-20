package manager

import (
	"context"
)

type CallbackInitializer struct {
	callback func(context.Context) error
}

func NewCallbackInitializer(callback func(context.Context) error) *CallbackInitializer {
	return &CallbackInitializer{
		callback: callback,
	}
}

func (ci *CallbackInitializer) Start(ctx context.Context) error {
	if ci.callback == nil {
		return nil
	}
	return ci.callback(ctx)
}
