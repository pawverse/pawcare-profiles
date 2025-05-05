package domain

import (
	"fmt"
	"time"
)

type PetProfile struct {
	Name        string
	DateOfBirth time.Time
	Species     string
	Breed       string
	Weight      float64
	Gender      EGender
}

func NewPetProfile(name string, dateOfBirth time.Time, species string, breed string, weight float64, gender EGender) PetProfile {
	return PetProfile{
		Name:        name,
		DateOfBirth: dateOfBirth,
		Species:     species,
		Breed:       breed,
		Weight:      weight,
		Gender:      gender,
	}
}

func (p PetProfile) String() string {
	return fmt.Sprintf("PetProfile(name=%s,dob=%s,species=%s,breed=%s,weight=%f,gender=%s)", p.Name, p.DateOfBirth, p.Species, p.Breed, p.Weight, p.Gender)
}
