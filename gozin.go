// Package gozin is a runtime select for golang based on reflection. The select command in go
// is an essential part of concurrency, but having a dynamic select is only possible via reflection.
// If you have several channels, and you know them at the coding time, then the go select
// is your ultimate answer. but if you want to watch channels that are created programmatically,
// then you have to use reflection. هف
// Package gozin is a simple thin wrapper around the reflection call. It is slower than
// go select (what you expect!?) but if you need  it then you have to use it.
// In Persian, gozin (گزین) means chosen.
package gozin

import (
	"errors"
	"reflect"
)

// Case is an internal interface used to handle the cases. you need to call gozin.{Send,Receive,Default}
// functions to create a case.
type Case interface {
	getSelect() reflect.SelectCase
	call(interface{}, bool)
}

type selectCase struct {
	data  reflect.SelectCase
	onRec func(interface{}, bool)
	on    func()
}

func (s *selectCase) getSelect() reflect.SelectCase {
	return s.data
}

func (s *selectCase) call(data interface{}, ok bool) {
	switch s.data.Dir {
	case reflect.SelectDefault:
		if s.on != nil {
			s.on()
		}
	case reflect.SelectSend:
		if s.on != nil {
			s.on()
		}
	case reflect.SelectRecv:
		if s.onRec != nil {
			s.onRec(data, ok)
		}
	}
}

// Receive returns a receive case. it panics if the channel is not a channel or the channel is write only ( chan <-)
// the on function called if this case is selected. the first parameter of the on function is the
// data read from the channel, and the second parameter is false if the channel is closed
func Receive(channel interface{}, on func(interface{}, bool)) Case {
	v := reflect.ValueOf(channel)
	if v.Kind() != reflect.Chan {
		panic("channel should be a go channel")
	}

	if d := v.Type().ChanDir(); d == reflect.SendDir {
		panic("channel is write only")
	}

	return &selectCase{
		data: reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: v,
		},
		onRec: on,
	}
}

// Send creates a send case. it panics if the channel is not a channel, or its read only, or the type of data is
// not exactly the type of channel element.
// the on function is called if the data sent to the channel
func Send(channel interface{}, data interface{}, on func()) Case {
	v := reflect.ValueOf(channel)
	if v.Kind() != reflect.Chan {
		panic("channel should be a go channel")
	}

	if d := v.Type().ChanDir(); d == reflect.RecvDir {
		panic("channel is read only")
	}

	d := reflect.ValueOf(data)
	if !d.Type().AssignableTo(v.Type().Elem()) {
		panic("the data type is not acceptable for channel")
	}

	return &selectCase{
		data: reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: v,
			Send: d,
		},
		on: on,
	}
}

// Default is the default case.
func Default(on func()) Case {
	return &selectCase{
		data: reflect.SelectCase{
			Dir: reflect.SelectDefault,
		},
		on: on,
	}
}

// Select is the main function to call. it returns an error if there is more than one default case
// is in the cases. other than that, the behavior is exactly the same as go select.
func Select(cases ...Case) error {
	c := make([]reflect.SelectCase, len(cases))
	var def bool
	for i := range cases {
		c[i] = cases[i].getSelect()
		if c[i].Dir == reflect.SelectDefault {
			if def {
				return errors.New("multiple default case")
			}
			def = true
		}
	}

	idx, val, ok := reflect.Select(c)
	var res interface{}
	if val.Kind() != reflect.Invalid && val.CanInterface() {
		res = val.Interface()
	}

	cases[idx].call(res, ok)
	return nil
}
