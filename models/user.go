package models

import (
	"errors"
	"time"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DOB         int32  `json:"dob"`
	Address     string `json:"address"`
	Description string `json:"description"`
	Ctime       int32  `json:"ctime"`
}

func (u *User) Validate() (err error) {
	if u == nil {
		return errors.New("empty user")
	}
	if u.ID != "" {
		return errors.New("invalid user id")
	}
	if u.Name == "" {
		return errors.New("invalid username")
	}
	if u.DOB >= int32(time.Now().UnixNano()) {
		return errors.New("invalid dob")
	}
	if u.Address == "" {
		return errors.New("invalid address")
	}
	if u.Description == "" {
		return errors.New("invalid description")
	}
	if u.Ctime > int32(time.Now().UnixNano()) {
		return errors.New("invalid ctime")
	}
	return nil
}
