package gomapper

import (
	"fmt"
	"testing"
)

type Email string

type UserDTO struct {
	ID    int64
	Age   int
	Name  string
	Email Email
}

type UserDAL struct {
	ID        int64
	Age       uint8
	Username  string
	EmailAddr string
}

func TestMap(t *testing.T) {
	cfg := NewConfig()

	Map[UserDTO, UserDAL](cfg)

	fmt.Printf("%+v\n", cfg)
}
