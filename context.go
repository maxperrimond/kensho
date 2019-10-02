package kensho

import "context"

type (
	ValidationContext struct {
		ctx           context.Context
		violationList ViolationList
		path          string
		arg           interface{}
		root          interface{}
		subject       interface{}
		value         interface{}
	}
)

func (ctx *ValidationContext) WithValue(value interface{}) *ValidationContext {
	clone := *ctx
	clone.value = value

	return &clone
}

func (ctx *ValidationContext) WithArg(arg interface{}) *ValidationContext {
	clone := *ctx
	clone.arg = arg

	return &clone
}

func (ctx *ValidationContext) Ctx() context.Context {
	return ctx.ctx
}

func (ctx *ValidationContext) ViolationList() ViolationList {
	return ctx.violationList
}

func (ctx *ValidationContext) BuildViolation(error string, parameters map[string]interface{}) *violationBuilder {
	return &violationBuilder{
		error:      error,
		message:    TranslateError(error, parameters),
		parameters: parameters,
		basePath:   ctx.path,
		onAdd: func(violation Violation) {
			ctx.violationList = append(ctx.violationList, &violation)
		},
	}
}

func (ctx *ValidationContext) Path() string {
	return ctx.path
}

func (ctx *ValidationContext) Arg() interface{} {
	return ctx.arg
}

func (ctx *ValidationContext) Root() interface{} {
	return ctx.root
}

func (ctx *ValidationContext) Subject() interface{} {
	return ctx.subject
}

func (ctx *ValidationContext) Value() interface{} {
	return ctx.value
}
