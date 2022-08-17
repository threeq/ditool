package fsm_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/threeq/ditool/fsm"
	"testing"
)

func TestFSM(t *testing.T) {
	sm := fsm.NewFSM("test_1")

	sm.
		On("1", "ok").
		CondActionTo("", fsm.Any, ok, "3")
	sm.
		On("3", "err").
		CondActionTo("", nil, err, "3")
	sm.
		On("3", "panic").
		CondActionTo("", nil, exception, "3")

	e := entity("1")
	assert.Nil(t, sm.Emit(context.Background(), e, "ok"))
	assert.Equal(t, "3", e.(*testEntity).s)
	assert.Equal(t, sm.Emit(context.Background(), e, "err").Error(),
		"action error")
	assert.Panics(t, func() {
		sm.Emit(context.Background(), e, "panic")
	})

	assert.Equal(t, sm.Emit(context.Background(), entity("1"), "no").Error(),
		"1@no 事件没有定义",
	)
	assert.Equal(t, sm.Emit(context.Background(), entity("2"), "ok").Error(),
		"2 状态没有定义")

}

func TestFsm_Condition(t *testing.T) {

}

func ok(ctx context.Context, entity fsm.Entity, src fsm.State, e fsm.Event, dst fsm.State) (fsm.Entity, error) {
	entity.(*testEntity).s = dst
	return entity, nil
}

func err(ctx context.Context, entity fsm.Entity, src fsm.State, e fsm.Event, dst fsm.State) (fsm.Entity, error) {
	return nil, errors.New("action error")
}

func exception(ctx context.Context, entity fsm.Entity, src fsm.State, e fsm.Event, dst fsm.State) (fsm.Entity, error) {
	panic("panic")
}

type testEntity struct {
	s fsm.State
}

func (t *testEntity) ID() string {
	return ""
}

func (t *testEntity) State() fsm.State {
	return t.s
}

func entity(state fsm.State) fsm.Entity {
	return &testEntity{
		s: state,
	}
}
