package domain

import (
	"fmt"
)

type OwnerId string
type UserId string

type Owner struct {
	Id      OwnerId
	UserId  UserId
	Profile OwnerProfile
}

func NewOwner(userId UserId, profile OwnerProfile) *Owner {
	return &Owner{
		UserId:  userId,
		Profile: profile,
	}
}

func (o Owner) String() string {
	return fmt.Sprintf("Owner(id=%s,userId=%s,profile=%s)", o.Id, o.UserId, o.Profile)
}
