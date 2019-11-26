package logger

import (
	"errors"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// func OtherWork() {
// 	logger = log.New(os.Stderr, "", 0)
// }

func DemoGlobal() {
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
	}
}

// DemoV1 function needs to use logger whenever there is an error so it is
// dependent on a logger.
func DemoV1() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
	}
}

// This way, the dependency will be provided.
func DemoV2(logger *log.Logger) {
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
	}
}

// DemoV3 works because Println is a variadic function that accepts any number
// of interfaces. It's implementation agnostic.
// Call this with:
//   logger := log.New(...)
//   DemoV3(log.Println)
func DemoV3(logFn func(...interface{})) {
	err := doTheThing()
	if err != nil {
		logFn("error in doTheThing():", err)
	}
}

// https://golang.org/pkg/log/#Logger.Printf
type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

// DemoV4 is still being implementation agnostic through interfaces.
// Call this with:
//   logger := log.New(...)
//   DemoV4(logger)
func DemoV4(logger Logger) {
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
		logger.Printf("error: %s\n", err)
	}
}

/*
// This can be implemented too like this, without splicitly defining the type:
func DemoV4(logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}) {
	// ...
}
*/

// var defaultThing Thing

// func DemoV5() {
// 	defaultThing.DemoV5()
// }

type Thing struct {
	Logger interface {
		Println(...interface{})
		Printf(string, ...interface{})
	}
	// Printer interface {
	// 	...
	// }
}

func (t Thing) DemoV5() {
	err := doTheThing()
	if err != nil {
		t.Logger.Println("error in doTheThing():", err)
		t.Logger.Printf("error: %s\n", err)
	}
}

func DemoV6(logger Logger) {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}
	err := doTheThing()
	if err != nil {
		logger.Println("error in doTheThing():", err)
		logger.Printf("error: %s\n", err)
	}
}

type ThingV2 struct {
	Logger interface {
		Println(...interface{})
		Printf(string, ...interface{})
	}
	once sync.Once
}

func (t *ThingV2) logger() Logger {
	t.once.Do(func() {
		if t.Logger == nil {
			t.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		}
	})
	return t.Logger
}

func (t *ThingV2) DemoV7() {
	err := doTheThing()
	if err != nil {
		t.logger().Println("error in doTheThing():", err)
		t.logger().Printf("error: %s\n", err)
	}
}

func doTheThing() error {
	// return nil
	return errors.New("error opening file: abc.txt")
}
