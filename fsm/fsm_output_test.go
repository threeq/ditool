package fsm_test

import (
	"github.com/threeq/ditool/fsm"
	"testing"
)

func simpleSM() *fsm.FSM {
	sm := fsm.NewFSM("test_1")

	sm.
		On(fsm.Start, "23").
		CondActionTo("df", nil, ok, "1")
	sm.
		On("1", "ok").
		CondActionTo("", fsm.Any, ok, "3")
	sm.
		On("3", "err").
		CondActionTo("err>1", nil, err, "3")
	sm.
		On("3", "panic").
		CondActionTo("err<1", nil, exception, "3")
	return sm
}

func TestFSM_Text(t *testing.T) {
	sm := simpleSM()

	println(sm.Text())
}

func TestFSM_PlantUML(t *testing.T) {
	sm := simpleSM()
	println(sm.PlantUML())
}
