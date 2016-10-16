package main

import "fmt"

type Int int

func (i *Int) Incr(n int) Int {
	*i = (*i) + Int(n)
	return *i
}


func main() {
	i := Int(10)
	i.Incr(12)
	fmt.Println(i)

	b := &Bytes{b: []byte{1, 2, 3, 4, 5}, i: 0}
	c := b.Get(0, 3)
	c[1] = 11
	fmt.Println(b)
	fmt.Println(c)

	j := &Bytes{b: []byte{1, 2, 3, 4, 5}, i: 0}
	k := j.IGet(3)
	l := j.IGet(2)
	k[0] = 10
	l[0] = 11
	fmt.Println(j)
	fmt.Println(k)
	fmt.Println(l)
}
