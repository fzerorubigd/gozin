package gozin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSelect(b *testing.B) {
	t1 := make(chan int, 2)
	t2 := make(chan int, 2)

	for i := 0; i < b.N; i++ {
		select {
		case <-t1:
		case t2 <- 10:
		default:

		}
	}
}

func BenchmarkGozin(b *testing.B) {
	t1 := make(chan int, 2)
	t2 := make(chan int, 2)

	for i := 0; i < b.N; i++ {
		_ = Select(
			Receive(t1, func(i interface{}, b bool) {

			}),
			Send(t2, 10, func() {

			}),
			Default(func() {

			}),
		)
	}
}

func ExampleSend() {
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

func TestSend(t *testing.T) {
	assert.Panics(t, func() { Send(100, 100, nil) })
	c := make(chan string)
	assert.Panics(t, func() { Send(c, 100, nil) })
	var c2 <-chan string = c
	assert.Panics(t, func() { Send(c2, "str", nil) })
}

func TestReceive(t *testing.T) {
	assert.Panics(t, func() { Receive(100, nil) })
	c := make(chan string)
	var c2 chan<- string = c
	assert.Panics(t, func() { Receive(c2, nil) })
}

func TestMultipleDefault(t *testing.T) {
	err := Select(Default(nil), Default(nil))
	assert.Error(t, err)
}

func TestSelect(t *testing.T) {
	called := false
	err := Select(
		Default(func() {
			called = true
		}),
	)
	assert.NoError(t, err)
	assert.True(t, called)

	c := make(chan string, 1)
	called = false
	err = Select(
		Default(func() {
			assert.Fail(t, "should not called")
		}),
		Send(c, "string", func() {
			called = true
		}),
	)
	assert.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, "string", <-c)
}
