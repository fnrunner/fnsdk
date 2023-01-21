package fn

import "context"

var _ context.Context = &Context{}

type Context struct {
	context.Context
}
