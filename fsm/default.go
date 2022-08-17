package fsm

import (
	"context"
	"log"
)

const (
	Start State = "[*]"
	End   State = "[*]"
)

func Any(context.Context, Entity, State, Event, State) bool {
	return true
}

type defaultLogger struct {
}

func (d *defaultLogger) Info(c context.Context, msg string) {
	log.Println(msg)
}

func (d *defaultLogger) Warn(c context.Context, msg string) {
	log.Println(msg)
}

func (d *defaultLogger) Debug(c context.Context, msg string) {
	log.Println(msg)
}

func (d *defaultLogger) Error(c context.Context, msg string) {
	log.Println(msg)
}
