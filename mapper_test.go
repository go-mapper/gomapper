package mapper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

	err := Map[UserDTO, UserDAL](cfg)
	require.NoError(t, err)

	fmt.Printf("%+v\n", cfg)
}
