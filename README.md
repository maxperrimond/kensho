# Kenshō - A Go validator

[![Build Status](https://travis-ci.org/maxperrimond/kensho.svg?branch=master)](https://travis-ci.org/maxperrimond/kensho)
[![Coverage Status](https://coveralls.io/repos/github/maxperrimond/kensho/badge.svg?branch=master)](https://coveralls.io/github/maxperrimond/kensho?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxperrimond/kensho)](https://goreportcard.com/report/github.com/maxperrimond/kensho)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**A *Work in progress* project so the `master` branch might change in the future and BC breaks some part or change some behaviors.**

A simple Go library for validation, but gives the possibility to validate deeply, collections, any field type of struct by following tag or file.

## TODO

 - More docs/examples
 - More tests
 - Add validators

## Features

 - Struct validation
 - Able to add custom constraints
 - Validator argument
 - Configuration from a file (ex: json)
 - Deep validation
 - List of struct validation
 - Translation

## Installation

    go get github.com/maxperrimond/kensho

## Usage

Here some code as an example:

```go
package main

import (
	"fmt"

	"github.com/maxperrimond/kensho"
)

type Group struct {
	Name  string  `valid:"required"`
	Users []*User `valid:"valid"` // Ask to valid users if there is some
}

type User struct {
	Email     string `valid:"required,email"`
	FirstName string `valid:"required,min=2,max=64"`
	LastName  string `valid:"required,min=2,max=64"`
}

func main() {
	// Simple struct
	user := &User{
		Email:     "foo.bar@example.com",
		FirstName: "foo",
		LastName:  "bar",
	}

	// Validate user
	ok, violations, _ := kensho.Validate(user)

	fmt.Printf("Result: %t\n", ok)
	fmt.Println(violations)

	user.Email = "this is not an email"
	user.FirstName = ""

	// Validate user after inserting bad data
	ok, violations, _ = kensho.Validate(user)

	fmt.Printf("Result: %t\n", ok)
	fmt.Println(violations)

	users := []*User{
		{
			Email:     "john@example.com",
			FirstName: "john",
			LastName:  "bar",
		},
		{
			Email:     "pierre@example.com",
			FirstName: "pierre",
			LastName:  "bar",
		},
	}

	// Validate collection of users
	ok, violations, _ = kensho.Validate(users)

	fmt.Printf("Result: %t\n", ok)
	fmt.Println(violations)

	// Nested struct
	group := &Group{
		Name:  "foo",
		Users: append(users, user), // With the bad user
	}

	// Validate the group
	ok, violations, _ = kensho.Validate(group)

	fmt.Printf("Result: %t\n", ok)
	fmt.Println(violations)
}
```

### Available constraints

Tag | Arg | Description
--- | --- | ---
valid | | Validate nested struct
string | | Check is string
struct | | Check is struct
required | | Is required
length | int | Match length
min | int | Match min length
max | int | Match max length
regex | string (regex pattern) | Match pattern
email | | Match email
uuid | | Match UUID
iso3166 | optional: `num`, `alpha3`, `alpha2` (default) | Match country code based on [IS03166](https://en.wikipedia.org/wiki/ISO_3166)

And more to come

### Custom constraint

```go
package main

import (
	"context"
	"errors"

	"github.com/maxperrimond/kensho"
)

// Define your constraint
func poneyConstraint(ctx *kensho.ValidationContext) error {
	if ctx.Value() != "poney" {
		ctx.BuildViolation("invalid_poney", nil).AddViolation()
	}

	return nil
}

func init() {
	// add it with the tag of your choice
    kensho.AddConstraint("poney", poneyConstraint)
}
```

Note: If you use an existent tag, it will override it.

### Context

You can pass a context during validation so it can be accessible in constraints:

```go
ok, violations, err := kensho.ValidateWithContext(ctx, subject)
```
