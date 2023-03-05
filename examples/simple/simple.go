package simple

import (
	"github.com/go-mapper/gomapper"
)

type UserDTO struct {
	ID    int64
	Age   int
	Name  string
	Email string
}

type UserDAL struct {
	ID    int64
	Age   uint8
	Name  string
	Email string
}

func ConfigureMapper(cfg *gomapper.Config) {
	gomapper.Map[UserDTO, UserDAL](cfg)
}
