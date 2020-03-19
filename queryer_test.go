package gopsql

import (
	"fmt"
	"testing"
)

func TestQuery_Prepare(t *testing.T) {
	res := prepare("SELECT * from bhjb jh wjere 1 = ? and b = ? and f = ?")
	fmt.Println(res)
}
