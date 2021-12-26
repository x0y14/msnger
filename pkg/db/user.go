package db

import "time"

type UserData struct {
	Id          string    `db:"id"`
	DisplayName string    `db:"displayName"`
	StatusText  string    `db:"statusText"`
	PictureUrl  string    `db:"pictureUrl"`
	CreatedAt   time.Time `db:"createdAt"`
	UpdatedAt   time.Time `db:"updatedAt"`
}
