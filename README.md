# Kensho - A Go validator

**A *Work in progress* project so the `master` branch might change in the future and BC breaks some part or change some behaviors.**

A simple Go library for validation, but gives the possibility to validate deeply, collections, any field type of struct by following tag or file.

## Features

 - Struct validation
 - Allow custom validator
 - Validator argument
 - Configuration from a file (ex: json)
 - Deep validation
 - List of struct validation

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
	ok, _ := kensho.Validate(user)

	fmt.Printf("Result: %t\n", ok)

	user.Email = "this is not an email"
	user.FirstName = ""

	// Validate user after inserting bad data
	ok, err := kensho.Validate(user)

	fmt.Printf("Result: %t\n", ok)
	fmt.Printf("Email errors: %v\n", err.Fields["Email"].Errors)
	fmt.Printf("First name errors: %v\n", err.Fields["FirstName"].Errors)

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
	ok, _ = kensho.Validate(users)

	fmt.Printf("Result: %t\n", ok)

	// Nested struct
	group := &Group{
		Name:  "foo",
		Users: append(users, user), // With the bad user
	}

	// Validate the group
	ok, err = kensho.Validate(group)

	fmt.Printf("Result: %t\n", ok)
	fmt.Printf("Email errors: %v\n", err.Fields["Users"].Fields["2"].Fields["Email"].Errors)
	fmt.Printf("First name errors: %v\n", err.Fields["Users"].Fields["2"].Fields["FirstName"].Errors)
}
```

### Available validators

Tag | Arg | Description
--- | --- | ---
valid | | Useful tag to force nested validation if you no need to validate the struct if self.
string | | If it's string
struct | | If it's struct
required | | Is required
length | int | Check the given length
min | int | Check if it has the required min length
max | int | Check if it has the required max length
regex | string (regex pattern) | Match the given pattern
email | | Match if it's an email
uuid | | Match if it's an UUID

For now

### Custom validator

To add a new validator:

```go
package main

import (
	"context"
	"errors"

	"github.com/maxperrimond/kensho"
)

// Define your validator
func poneyValidator(ctx context.Context, subject interface{}, value interface{}, arg interface{}) (bool, error) {
	if value == "poney" {
		return true, nil
	}

	return false, errors.New("this is not a poney!!!")
}

func init() {
	// add it with the tag of your choice
    kensho.AddValidator("poney", poneyValidator)
}
```

param | Description
---|---
ctx  | validation context
subject | parent struct
value | field value (or array value)
arg | validator argument from struct configuration

Note: If you use an existent tag, it will override it.

### Context

You can pass a context during validation and can be used by any validator:

```go
ok, err := kensho.ValidateWithContext(myCtx, myStruct)
```
