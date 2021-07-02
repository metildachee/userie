package model

type User struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	DOB         int32  `json:"dob"`
	Address     string `json:"address"`
	Description string `json:"description"`
	Ctime       int32  `json:"ctime"`
}
