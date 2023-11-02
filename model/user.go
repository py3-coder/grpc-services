package model

import "time"

type UserDetails struct {
	ID        int       `bson:"id" json:"id"`
	FName     string    `bson:"fname" json:"fname"`
	City      string    `bson:"city" json:"city"`
	Phone     int64     `bson:"phone" json:"phone"`
	Height    float64   `bson:"height" json:"height"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
type User struct {
	ID     int     `bson:"id" json:"id"`
	FName  string  `bson:"fname" json:"fname"`
	City   string  `bson:"city" json:"city"`
	Phone  int64   `bson:"phone" json:"phone"`
	Height float64 `bson:"height" json:"height"`
}
type UserRespone struct {
	StatusCode int    `json:"statuscode"`
	Status     string `json:"status"`
}
