package tests

import (
	"fmt"
	"testing"
	"time"
)

import (
	. "github.com/ubrabbit/go-common/common"
	. "github.com/ubrabbit/go-common/event"
)

type EventTest_1 struct {
	eventName string
	eventID   int
}

func (self *EventTest_1) ID() int {
	return self.eventID
}

func (self *EventTest_1) Name() string {
	return self.eventName
}

func (self *EventTest_1) Execute(args ...interface{}) {
	v1 := args[0].(int)
	v2 := args[1].(string)
	fmt.Println("EventTest_1  Execute:  ", v1, v2)
}

func TestEvent_1(t *testing.T) {
	fmt.Printf("\n\n=====================  TestEvent_1  =====================\n")

	InitEvent()

	AddEvent(&EventTest_1{eventName: "event_1", eventID: 1})
	AddEvent(&EventTest_1{eventName: "event_2", eventID: 2})
	TriggerEvent("event_1", 111, "event_1 trigger_0")

	err := AddEvent(&EventTest_1{eventName: "event_1", eventID: 1})
	if err != nil {
		fmt.Println("AddEvent Error: ", err)
	}

	time.Sleep(1 * time.Millisecond)
	TriggerEvent("event_1", 111, "event_1 trigger_0-1")
	RemoveEvent(1)
	TriggerEvent("event_1", 111, "event_1 trigger_1")
	TriggerEvent("event_1", 111, "event_1 trigger_2")
	TriggerEvent("event_1", 111, "event_1 trigger_3")
	RemoveEvent(3)
	RemoveEvent(4)
	TriggerEvent("event_2", 222, "event_2 trigger")

	time.Sleep(3 * time.Second)
	LogInfo("Finished")
}
