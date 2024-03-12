package main

import (
	"fmt"

	"github.com/thockin/go-spew/spew"
)

type List struct {
	name string
	Data *Data
	Next *List
	Prev *List
}

type Data struct {
	val  int
	tree []*Data
}

var head *List

func init() {
	head = &List{name: "head"}
	head.Prev = head

	p := head
	for i := 0; i < 3; i++ {
		p.Data = &Data{val: i * 100}
		for j := 0; j < 3; j++ {
			n := &Data{val: i*100 + j*10}
			p.Data.tree = append(p.Data.tree, n)
		}
		p.Next = &List{name: fmt.Sprintf("+%d", i+1)}
		p.Next.Prev = p
		p = p.Next
	}
}

func main() {
	spew.CleanConfig.Dump(head)
}
