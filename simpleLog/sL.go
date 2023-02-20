package simpleLog

import "fmt"

type SimpleLogger struct {
}

func NewSimpleLogger() SimpleLogger {
	return SimpleLogger{}
}

func (sL SimpleLogger) Info(format string, a ...any) {
	fmt.Printf("Info: "+format+"\n", a)
}

func (sL SimpleLogger) Error(format string, a ...any) {
	fmt.Printf("Error: "+format+"\n", a)
}
