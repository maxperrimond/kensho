package kensho

import (
	"context"
	"regexp"
)

const (
	email    = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	uuid     = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	colorHex = "^#(?:[0-9a-fA-F]{3}){1,2}$"
)

func validWithRegex(ctx context.Context, args ConstraintArgs, error *Error) *Error {
	if args.Value == nil {
		return nil
	}

	err := StringConstraint(ctx, args)
	if err != nil {
		return err
	}

	if args.Value == "" {
		return nil
	}

	pattern, ok := args.Arg.(string)
	if !ok {
		panic("the pattern is missing to validate with a regex")
	}

	if regexp.MustCompile(pattern).MatchString(args.Value.(string)) {
		return nil
	}

	return error
}

func RegexConstraint(ctx context.Context, args ConstraintArgs) *Error {
	return validWithRegex(ctx, args, &Error{
		Message: TranslateError("not_match_regex", nil),
		Error:   "not_match_regex",
	})
}

func EmailConstraint(ctx context.Context, args ConstraintArgs) *Error {
	return validWithRegex(ctx, ConstraintArgs{
		Root:    args.Root,
		Subject: args.Subject,
		Value:   args.Value,
		Arg:     email,
	}, &Error{
		Message: TranslateError("invalid_email", nil),
		Error:   "invalid_email",
	})
}

func UUIDConstraint(ctx context.Context, args ConstraintArgs) *Error {
	return validWithRegex(ctx, ConstraintArgs{
		Root:    args.Root,
		Subject: args.Subject,
		Value:   args.Value,
		Arg:     uuid,
	}, &Error{
		Message: TranslateError("invalid_uuid", nil),
		Error:   "invalid_uuid",
	})
}

func ColorHexConstraint(ctx context.Context, args ConstraintArgs) *Error {
	return validWithRegex(ctx, ConstraintArgs{
		Root:    args.Root,
		Subject: args.Subject,
		Value:   args.Value,
		Arg:     colorHex,
	}, &Error{
		Message: TranslateError("invalid_color", nil),
		Error:   "invalid_color",
	})
}
