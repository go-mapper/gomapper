package gomapper

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOutput(t *testing.T) {
	buf := new(bytes.Buffer)

	cfg := NewConfig()

	Map[UserDTO, UserDAL](cfg)

	Output(buf, cfg)

	fmt.Println(buf.String())
}
