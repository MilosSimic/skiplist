# skiplist
Simple implementation of skiplist algorithm

# Usage:
```go
package main

import (
	"fmt"
)

func main() {
	sl := skiplist.New(32, 8)
	sl.Add("a", []byte("A"))
	sl.Add("b", []byte("B"))
	sl.Add("c", []byte("C"))
	sl.Add("d", []byte("D"))
	sl.Add("e", []byte("E"))
	sl.Add("f", []byte("F"))
	sl.Add("g", []byte("G"))
	sl.Add("h", []byte("H"))
	sl.Add("i", []byte("I"))
	sl.Add("j", []byte("J"))
	sl.Add("k", []byte("K"))
	sl.Add("l", []byte("L"))
	sl.Add("m", []byte("M"))
	sl.Add("n", []byte("N"))
	sl.Add("o", []byte("O"))
	sl.Add("p", []byte("P"))
	sl.Add("q", []byte("Q"))
	sl.Add("r", []byte("R"))
	sl.Add("s", []byte("S"))

	fmt.Println(sl.Contains("l"))
	fmt.Println(sl.Contains("z"))

	fmt.Println(string(sl.Get("l")[:]))
	fmt.Println(string(sl.Get("z"))[:])

	fmt.Println(sl.Remove("l"))
	fmt.Println(sl.Remove("z"))

	node := sl.head
	prep := map[string]Value{}
	sl.Prep(node.next, prep)
	fmt.Println(prep)

}

```
