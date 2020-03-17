package gopsql

import (
	"fmt"
	"testing"
)

func TestQuery_Prepare(t *testing.T) {
	q := Builder{}
	res := q.Prepare("SELECT * from bhjb jh wjere 1 = ? and b = ? and f = ?")
	fmt.Println(res)
}
