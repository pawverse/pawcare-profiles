package biz

import (
	"fmt"
	"time"
)

type (
	OwnerId string
	UserId  string
)

type Owner struct {
	Id          OwnerId
	UserId      UserId
	Name        string
	DateOfBirth time.Time
}

func (o Owner) String() string {
	return fmt.Sprintf("Owner(id=%s, userId=%s, name=%s, dateOfBirth=%s)", o.Id, o.UserId, o.Name, o.DateOfBirth.String())
}
