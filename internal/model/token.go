package model

type Token struct {
	UserID uint32 `db:"user_id"`
	Value  string `db:"token"`
}
