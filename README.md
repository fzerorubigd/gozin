# Gozin - Golang dynamic select

[![Build Status](https://travis-ci.org/fzerorubigd/gozin.svg)](https://travis-ci.org/fzerorubigd/gozin)
[![Coverage Status](https://coveralls.io/repos/github/fzerorubigd/gozin/badge.svg?branch=master)](https://coveralls.io/github/fzerorubigd/gozin?branch=master)
[![GoDoc](https://godoc.org/github.com/fzerorubigd/gozin?status.svg)](https://godoc.org/github.com/fzerorubigd/gozin)
[![Go Report Card](https://goreportcard.com/badge/github.com/fzerorubigd/gozin/die-github-cache-die)](https://goreportcard.com/report/github.com/fzerorubigd/gozin)

Gozin is a thin wrapper around the reflection api for select to create dynamic select at runtime. 

```go 
func main() {
	t1 := make(chan int)
	t2 := make(chan int)
	close(t1)
	err := Select(
		Default(func() {
			fmt.Println("Default")
		}), Send(t2, 111, func() {
			fmt.Println("Send")
		}), Receive(t1, func(data interface{}, ok bool) {
			fmt.Println("Rec", data, ok)
		}))

	if err != nil {
		panic(err)
	}

	// Output: Rec 0 false
}
```

In Persian Gozin ([گزین](https://www.vajehyab.com/dehkhoda/%DA%AF%D8%B2%DB%8C%D9%86-5)) means chosen.
 
