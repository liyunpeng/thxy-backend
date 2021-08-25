package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestA(t *testing.T) {
	re:=regexp.MustCompile("[0-9]+")
	fmt.Println(re.FindAllString("abc123def", -1))
}
