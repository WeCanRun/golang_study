package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
)

type TaskState uint8

const (
	Scheduled TaskState = iota + 1
	Starting
	Running
	Finished
)

var stateList = [...]string{"", "Scheduled", "Starting", "Running", "Finished"}

func (t TaskState) String() string {
	return stateList[t]
}

func ParseTaskState(state string) TaskState {
	for i, s := range stateList {
		if s == state {
			return TaskState(i)
		}
	}
	return 0
}

type taskScheduler struct {
	ctx   context.Context
	Name  string    `json:"name"`
	State TaskState `json:"state"`
}

func NewTaskScheduler(ctx context.Context, taskName string) *taskScheduler {
	return &taskScheduler{
		ctx:   ctx,
		Name:  taskName,
		State: 0,
	}
}

func (t *taskScheduler) Key() string {
	return fmt.Sprintf("/scheduler/task/%s", t.Name)
}

func (t *taskScheduler) Notify(state TaskState) error {
	if state < Scheduled || state > Finished {
		return fmt.Errorf("error state: %d", state)
	}
	t.State = state
	log.Printf("console notify task %s to work\n", t.Name)
	return Put(t.ctx, t.Key(), t.State.String())
}

func (t *taskScheduler) WatchNotify() {
	watch := client.Watch(t.ctx, t.Key())
	for w := range watch {
		for _, e := range w.Events {
			switch e.Type {
			case mvccpb.PUT:
				t.Work(string(e.Kv.Value))
			case mvccpb.DELETE:
				t.close()
			}
		}
	}
}

func (t *taskScheduler) Work(state string) {
	t.State = ParseTaskState(state)
	switch state {
	case Scheduled.String():
		log.Printf("task %s is scheduled\n", t.Name)
	case Starting.String():
		log.Printf("task %s begin  working\n", t.Name)
	case Running.String():
		log.Printf("task %s is running\n", t.Name)
	case Finished.String():
		log.Printf("task %s finnish working\n", t.Name)
		Del(t.ctx, t.Key())
	}
}

func (t *taskScheduler) close() {
	t.Name = ""
	t.State = 0
}
