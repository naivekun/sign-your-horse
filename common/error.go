package common

import (
	"log"
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func Raise(errmsg string) error {
	return &Error{msg: errmsg}
}

func Must(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
