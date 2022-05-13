package model

type Token struct {
	UserId uint32 `db:"user_id"`
	Value  string `db:"token"`
}
