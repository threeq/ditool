package fsm

import (
	"fmt"
	"strings"
)

func (fsm *FSM) Text() string {
	var transLine []string
	for state, eventTable := range fsm.transitionTable {
		for event, condTable := range eventTable {
			for _, t := range condTable {
				var l string
				if t.dec != "" {
					l = fmt.Sprintf("%v --> %v: %v [%v]", state, t.dst, event, t.dec)
				} else {
					l = fmt.Sprintf("%v --> %v: %v", state, t.dst, event)
				}
				transLine = append(transLine, l)
			}
		}
	}

	return strings.Join(transLine, "\n")
}

func (fsm *FSM) PlantUML() string {
	var transLine []string
	transLine = append(transLine, "\n@startuml")
	transLine = append(transLine, fmt.Sprintf("state %v {", fsm.name))

	for state, eventTable := range fsm.transitionTable {
		for event, condTable := range eventTable {
			for _, t := range condTable {
				var l string
				if t.dec != "" {
					l = fmt.Sprintf("%v --> %v: %v [%v]", state, t.dst, event, t.dec)
				} else {
					l = fmt.Sprintf("%v --> %v: %v", state, t.dst, event)
				}
				transLine = append(transLine, l)
			}
		}
	}

	transLine = append(transLine, "}")
	transLine = append(transLine, "@enduml")
	transLine = append(transLine, "\n请复制以上内容到网页上查看：https://www.plantuml.com/plantuml/uml")
	return strings.Join(transLine, "\n")
}
