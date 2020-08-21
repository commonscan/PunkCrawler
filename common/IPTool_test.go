package common

import (
	"fmt"
	"testing"
)

func TestGenerateIP(t *testing.T) {
	fmt.Println(GenerateIP("202.194.14.1/24"))
}
