package main

import (
	"errors"
	"testing"
)

func TestHandlerErr(t *testing.T) {
	err := errors.New("我是一个错误我显示在第几行吗")
	handlerErr(err)
}
