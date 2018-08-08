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

	formError := err.ToFormErrors()

	fmt.Printf("Result: %t\n", ok)
	fmt.Printf("Email errors: %v\n", formError.Fields["Email"].Errors)
	fmt.Printf("First name errors: %v\n", formError.Fields["FirstName"].Errors)

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

	formError = err.ToFormErrors()

	fmt.Printf("Result: %t\n", ok)
	fmt.Printf("Email errors: %v\n", formError.Fields["Users"].Fields["2"].Fields["Email"].Errors)
	fmt.Printf("First name errors: %v\n", formError.Fields["Users"].Fields["2"].Fields["FirstName"].Errors)
}
