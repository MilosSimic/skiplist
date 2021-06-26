# skiplist
Simple implementation of skiplist data structure

# Usage:
```go
package main

import (
	"fmt"
	sk "github.com/MilosSimic/skiplist
)

func main() {
	sl := sk.New(32, 8)
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
	fmt.Println("SIZE", sl.Size())

	fmt.Println(sl.Remove("a"))
	fmt.Println(sl.TombstoneIt("k"))
	fmt.Println(sl.Remove("s"))
	fmt.Println("SIZE", sl.Size())
	d, err := sl.Get("k")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(d)
	}

	node := sl.head
	prep := map[string]Entry{}
	sl.ToMap(node.next, prep)
	fmt.Println(prep)
}

```
