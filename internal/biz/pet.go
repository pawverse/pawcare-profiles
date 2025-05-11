package biz

import "time"

type PetId string

type Pet struct {
	Id          PetId
	OwnerId     OwnerId
	Name        string
	DateOfBirth time.Time
	Species     string
	Breed       string
	Weight      float64
	Gender      EGender
}

type EGender string

const (
	EGenderMale   EGender = "male"
	EGenderFemale EGender = "female"
	EGenderOther  EGender = "other"
)
