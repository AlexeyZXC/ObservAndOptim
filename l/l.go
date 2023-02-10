package l

import (
	"context"
	"elastic/m"
	"elastic/store_logs"
	"fmt"
	"time"
)

type ElLogger struct {
	store store_logs.LogsStore
}

func NewElasticLogger() (ElLogger, error) {
	store, err := store_logs.NewLogsStore()
	if err != nil {
		return ElLogger{}, err
	}

	return ElLogger{store: store}, nil
}

func (i ElLogger) Log(format string, a ...any) {
	str := fmt.Sprintf("Info: "+format, a)

	logTime := time.Now().Format("2006-01-02 15:04:03.000")
	str = logTime + ": " + str

	ctx := context.Background()
	msg := m.Logs{Message: str}

	i.store.Add(ctx, msg)
}

func (i ElLogger) Error(format string, a ...any) {
	str := fmt.Sprintf("Error: "+format, a)
	ctx := context.Background()

	logTime := time.Now().Format("2006-01-02 15:04:03.000")
	str = logTime + ": " + str

	msg := m.Logs{Message: str}

	i.store.Add(ctx, msg)
}

// var LogFunc = log.Println
// var FatalFunc = log.Fatal

// func F(i interface{}) {
// 	FatalFunc(Parse(i))
// }
// func Parse(i interface{}) interface{} {
// 	typeOf := reflect.TypeOf(i)
// 	if typeOf.Kind() == reflect.Map {
// 		b, err := json.MarshalIndent(i, "", "  ")
// 		if err != nil {
// 			return i
// 		}
// 		return string(b)
// 	}
// 	if typeOf.Kind() == reflect.Struct {
// 		b, err := json.MarshalIndent(i, "", "  ")
// 		if err != nil {
// 			return i
// 		}
// 		name := typeOf.Name()
// 		result := fmt.Sprintf("<%v>%v", name, string(b))
// 		return result
// 	}
// 	if typeOf.Kind() == reflect.Slice {
// 		v := reflect.ValueOf(i)
// 		result := "["
// 		for i := 0; i < v.Len(); i++ {
// 			val := v.Index(i)
// 			result += fmt.Sprintf("\n%v,", Parse(val.Interface()))
// 		}
// 		result = result[:len(result)-1] + "\n"
// 		result += "]"
// 		return result
// 	}
// 	if typeOf.Kind() == reflect.Array {
// 		v := reflect.ValueOf(i)
// 		result := "["
// 		for i := 0; i < v.Len(); i++ {
// 			val := v.Index(i)
// 			result += fmt.Sprintf("%v,", Parse(val.Interface()))
// 		}
// 		result += "]"
// 		return result
// 	}
// 	return i
// }
// func L(i interface{}) {
// 	LogFunc(Parse(i))
// }
