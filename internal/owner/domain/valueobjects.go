package domain

import (
	"fmt"
	"time"
)

type OwnerProfile struct {
	Name        string    `validate:"required"`
	DateOfBirth time.Time `validate:"required"`
}

func NewOwnerProfile(name string, dateOfBirth time.Time) OwnerProfile {
	return OwnerProfile{
		Name:        name,
		DateOfBirth: dateOfBirth,
	}
}

func (o OwnerProfile) String() string {
	return fmt.Sprintf("OwnerProfile(name=%s,dob=%s)", o.Name, o.DateOfBirth.String())
}
