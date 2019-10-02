package kensho

import (
	"regexp"
)

const (
	email    = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	uuid     = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	colorHex = "^#(?:[0-9a-fA-F]{3}){1,2}$"
)

func validWithRegex(ctx *ValidationContext, pattern string, onViolation func()) error {
	if ctx.Value() == nil {
		return nil
	}

	err := StringConstraint(ctx)
	if err != nil {
		return err
	}

	if ctx.Value() == "" {
		return nil
	}

	if regexp.MustCompile(pattern).MatchString(ctx.Value().(string)) {
		return nil
	}

	onViolation()

	return nil
}

func RegexConstraint(ctx *ValidationContext) error {
	pattern, ok := ctx.Arg().(string)
	if !ok {
		panic("the pattern is missing to validate with a regex")
	}

	return validWithRegex(ctx, pattern, func() {
		ctx.BuildViolation("not_match_regex", map[string]interface{}{
			"regex": ctx.Arg().(string),
		}).AddViolation()
	})
}

func EmailConstraint(ctx *ValidationContext) error {
	return validWithRegex(ctx, email, func() {
		ctx.BuildViolation("invalid_email", nil).AddViolation()
	})
}

func UUIDConstraint(ctx *ValidationContext) error {
	return validWithRegex(ctx, uuid, func() {
		ctx.BuildViolation("invalid_uuid", nil).AddViolation()
	})
}

func ColorHexConstraint(ctx *ValidationContext) error {
	return validWithRegex(ctx, colorHex, func() {
		ctx.BuildViolation("invalid_color", nil).AddViolation()
	})
}
