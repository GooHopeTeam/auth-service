package model

import (
	"fmt"
	"time"
)

type User struct {
	ID             uint32    `db:"id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
}

func (u *User) ToString() string {
	return fmt.Sprintf("User(%d, %s, %s, %s)", u.ID, u.Email, u.HashedPassword, u.CreatedAt)
}
