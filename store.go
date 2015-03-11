package session

import "time"

type Store interface {
	Add(id Id) bool
	Exist(id Id) bool
	Clear(id Id) bool

	Get(id Id, key string) interface{}
	Set(id Id, key string, value interface{}) error
	Del(id Id, key string) bool

	SetMaxAge(maxAge time.Duration)
	SetIdMaxAge(id Id, maxAge time.Duration)

	Run() error
}
