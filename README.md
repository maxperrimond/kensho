# Kensho - A Go validator

## Installation

    TODO

## Usage

Before any explanation, here a simple example:

```go
package main

import (
	"fmt"

	"github.com/maxperrimond/kensho"
)

type User struct {
	Email     string `valid:"required,email"`
	FirstName string `valid:"required,min=2,max=64"`
	LastName  string `valid:"required,min=2,max=64"`
}

func main() {
	user1 := &User{
		Email:     "foo.bar@example.com",
		FirstName: "foo",
		LastName:  "bar",
	}

	ok, _ := kensho.Validate(user1)

	fmt.Printf("Result: %t\n", ok)

	user1.Email = "this is not an email"
	user1.FirstName = ""

	ok, err := kensho.Validate(user1)

	fmt.Printf("Result: %t\n", ok)
	fmt.Printf("Email errors: %v\n", err.Fields["Email"].Errors)
	fmt.Printf("First name errors: %v\n", err.Fields["FirstName"].Errors)
}
```

## Contribute
