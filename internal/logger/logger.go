package logger

import (
	"context"
	"fmt"
	"log"
)

func Errorf(c context.Context, format string, args ...interface{}) {
	log.Println(fmt.Sprintf("ERROR: %s", fmt.Sprintf(format, args...)))
}
