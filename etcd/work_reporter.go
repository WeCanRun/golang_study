package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"strconv"
	"strings"
)

type workReporter struct {
	ctx      context.Context `json:"ctx"`
	Name     string          `json:"name"`
	progress int             `json:"progress"`
}

func NewWorkReporter(ctx context.Context, name string) *workReporter {
	return &workReporter{
		ctx:      ctx,
		Name:     name,
		progress: 0,
	}
}

func (w *workReporter) Key() string {
	return fmt.Sprintf("/work/%s", w.Name)
}

func (w *workReporter) FmtProgress() string {
	return fmt.Sprintf("%d", w.progress) + "%"
}

func (w *workReporter) ParseProgress(progress string) (int, error) {
	if !strings.HasSuffix(progress, "%") {
		panic("progress isn't end of %")
	}
	parse, err := strconv.Atoi(progress[0 : len(progress)-1])
	return parse, err
}

func (w *workReporter) Report(progress int) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("the progress of %s is error", w.Name)
	}
	w.progress = progress
	return Put(w.ctx, w.Key(), w.FmtProgress())
}

func (w *workReporter) WatchNotify() {
	watchChan := client.Watch(w.ctx, w.Key())
	for watch := range watchChan {
		for _, e := range watch.Events {
			switch e.Type {
			case mvccpb.PUT:
				progress := string(e.Kv.Value)
				parseInt, err := w.ParseProgress(progress)
				if err != nil {
					log.Printf("parseInt paogress of %s fail, progress: %s\n", w.Name, progress)
					break
				}
				log.Printf("the paogress of %s is %s\n", w.Name, progress)
				w.progress = parseInt
			}
		}
	}
}
