package l

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

//var Log Logger = Logger{ }

func CreateZapLogger() (Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return Logger{}, err
	}
	log := Logger{logger: logger}

	return log, nil
}

func (i Logger) Sync() error {
	return i.logger.Sync()
}

func (i Logger) Log(msg string) {
	i.logger.Log(zap.DebugLevel, msg)
}

func (i Logger) Error(str string) {
	i.logger.Error(str)
}

func (i Logger) Infof(msg string, args ...interface{}) {
	// fs := []zapcore.Field{}
	// for _, v := range args {
	// 	f := zapcore.Field{Interface: v}
	// 	fs = append(fs, f)
	// }

	i.logger.Sugar().Infof(msg, args...)
}

var LogFunc = log.Println
var FatalFunc = log.Fatal

func F(i interface{}) {
	FatalFunc(Parse(i))
}

func L(i interface{}) {
	LogFunc(Parse(i))
}

func Parse(i interface{}) interface{} {
	typeOf := reflect.TypeOf(i)
	if typeOf.Kind() == reflect.Map {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return i
		}
		return string(b)
	}
	if typeOf.Kind() == reflect.Struct {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return i
		}
		name := typeOf.Name()
		result := fmt.Sprintf("<%v>%v", name, string(b))
		return result
	}
	if typeOf.Kind() == reflect.Slice {
		v := reflect.ValueOf(i)
		result := "["
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			result += fmt.Sprintf("\n%v,", Parse(val.Interface()))
		}
		result = result[:len(result)-1] + "\n"
		result += "]"
		return result
	}
	if typeOf.Kind() == reflect.Array {
		v := reflect.ValueOf(i)
		result := "["
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			result += fmt.Sprintf("%v,", Parse(val.Interface()))
		}
		result += "]"
		return result
	}
	return i
}
