package logger

import (
	"context"
	"fmt"
	"log"
)

func Infof(c context.Context, format string, args ...interface{}) {
	log.Output(2, fmt.Sprintf(format, args...))
}

func Errorf(c context.Context, format string, args ...interface{}) {
	log.Println(fmt.Sprintf("ERROR: %s", fmt.Sprintf(format, args...)))
}

func Fatalf(c context.Context, format string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(format, args...))
}
