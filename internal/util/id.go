package util

import "github.com/aidarkhanov/nanoid"

var alphabet string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var size int = 14

func NewId() string {
	id, _ := nanoid.Generate(alphabet, size)
	return id
}
