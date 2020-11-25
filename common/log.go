package common

import (
	"fmt"
	"log"
)

func LogWithModule(name, msg string, format ...interface{}) {
	log.Println("[" + name + "]" + fmt.Sprintf(msg, format...))
}
