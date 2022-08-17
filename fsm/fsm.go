package fsm

import (
	"context"
	"errors"
	"fmt"
	"github.com/threeq/ditool/lockx"
	"time"
)

const LockPrefix = "FSM:"

type Entity interface {
	ID() string
	State() State
}

type Logger interface {
	Info(c context.Context, msg string)
	Warn(c context.Context, msg string)
	Debug(c context.Context, msg string)
	Error(c context.Context, msg string)
}

type Event = string
type State = string

type Cond = func(ctx context.Context, entity Entity, src State, e Event, dst State) bool
type Action = func(ctx context.Context, entity Entity, src State, e Event, dst State) (Entity, error)

type FSM struct {
	name            string
	lockerFactory   lockx.LockerFactory
	log             Logger
	transitionTable map[State]map[Event][]*transition
}

func (fsm *FSM) Emit(c context.Context, entity Entity, event Event) error {
	src := entity.State()
	transitions, err := fsm.transitions(src, event)
	if err != nil {
		return err
	}

	// 支持并发控制
	if fsm.lockerFactory != nil {
		locker, err := fsm.lockerFactory.Mutex(c,
			lockx.Key(fsm.lockerEntityID(entity.ID())),
			lockx.Retry(lockx.LimitRetry(lockx.LinearBackoff(50*time.Millisecond), 60)))
		if err != nil {
			return err
		}
		err = locker.Lock()
		if err != nil {
			return err
		}
		defer func() {
			e := locker.Unlock()
			fsm.log.Warn(c, "状态机释放 "+fsm.lockerEntityID(entity.ID())+" 锁失败："+e.Error())
		}()
	}

	// 出 状态 条件判断
	var transition *transition
	for _, trans := range transitions {
		if trans.cdt(c, entity, trans.src, event, trans.dst) {
			transition = trans
			fsm.log.Debug(c, fmt.Sprintf("%v@%v:%v 条件检查成功", src, event, trans.dec))
			break
		}
		fsm.log.Debug(c, fmt.Sprintf("%v@%v:%v 条件检查失败", src, event, trans.dec))
	}
	if transition == nil {
		return errors.New(fmt.Sprintf("%v@%v 所有条件检查均失败", src, event))
	}

	// 进 状态 操作逻辑
	_, err = transition.act(c, entity, src, event, transition.dst)

	if err == nil {
		fsm.log.Info(c, fmt.Sprintf("实体【%v】状态转移成功: %v --> %v :%v [%v] ",
			entity.ID(), src, transition.dst, event, transition.dec))
	} else {
		fsm.log.Error(c, fmt.Sprintf("实体【%v】状态转移失败: %v --> %v :%v [%v] ==> %v",
			entity.ID(), src, transition.dst, event, transition.dec, err.Error()))
	}

	return err
}

func (fsm *FSM) On(state State, event Event) *StateTransitionBuilder {
	return &StateTransitionBuilder{
		fsm:   fsm,
		src:   state,
		event: event,
	}
}

type StateTransitionBuilder struct {
	fsm   *FSM
	src   State
	event Event
}

func (stb *StateTransitionBuilder) CondActionTo(dec string, cond Cond, do Action, dst State) *StateTransitionBuilder {
	if cond == nil {
		cond = Any
	}

	src := stb.src
	event := stb.event
	if stb.fsm.transitionTable == nil {
		stb.fsm.transitionTable = make(map[State]map[Event][]*transition)
	}

	eventTable, ok := stb.fsm.transitionTable[src]
	if !ok {
		eventTable = make(map[Event][]*transition)
	}

	condTable, ok := eventTable[event]
	if !ok {
		condTable = []*transition{}
	}

	condTable = append(condTable, &transition{
		src,
		dst,
		event,
		cond,
		do,
		dec,
	})

	eventTable[event] = condTable
	stb.fsm.transitionTable[src] = eventTable

	return stb
}

func (fsm *FSM) transitions(src State, event Event) ([]*transition, error) {
	stateEvents, stateExist := fsm.transitionTable[src]
	if !stateExist {
		return nil, errors.New(fmt.Sprintf("%v 状态没有定义", src))
	}
	transitions, eventExist := stateEvents[event]
	if !eventExist {
		return nil, errors.New(fmt.Sprintf("%v@%v 事件没有定义", src, event))
	}
	return transitions, nil

}

type transition struct {
	src State
	dst State
	evt Event
	cdt Cond
	act Action
	dec string
}

func (fsm *FSM) lockerEntityID(id string) string {
	return LockPrefix + fsm.name + ":" + id
}

// --------------------------------------------------

func NewFSM(name string, Options ...Option) *FSM {
	sm := &FSM{
		name: name,
		log:  &defaultLogger{},
	}

	for _, option := range Options {
		option(sm)
	}

	return sm
}

type Option func(fsm *FSM)

func ConcurrentLocker(lf lockx.LockerFactory) Option {
	return func(fsm *FSM) {
		fsm.lockerFactory = lf
	}
}

func Log(log Logger) Option {
	return func(fsm *FSM) {
		fsm.log = log
	}
}
