package mapper

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutput(t *testing.T) {
	buf := new(bytes.Buffer)

	cfg := NewConfig()

	err := Map[UserDTO, UserDAL](cfg)
	require.NoError(t, err)

	Output(buf, cfg)

	fmt.Println(buf.String())
}
