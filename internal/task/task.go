package task

import (
	"time"

	"github.com/SayukiDev/VRCLotterySystem/internal/data"
	"github.com/SayukiDev/VRCLotterySystem/internal/provider"
	"github.com/SayukiDev/VRCLotterySystem/log"
	"go.uber.org/zap"
)

type Task struct {
	Tick   time.Duration
	p      *provider.Provider
	logger *zap.Logger
}

func NewTask(p *provider.Provider, tick time.Duration) *Task {
	return &Task{
		Tick:   tick,
		p:      p,
		logger: log.SubLogger("Task"),
	}
}

func (t *Task) DrawingTask() error {
	var showed bool
	var date time.Time
	t.p.Data.RLock(func(d *data.Content) {
		showed = d.Showed
		date = d.Date
	})
	if showed {
		return nil
	}
	if date.After(time.Now()) {
		return nil
	}
	ids, err := t.p.Drawing()
	if err != nil {
		return err
	}
	err = t.p.AddResults(ids)
	if err != nil {
		return err
	}

	// Todo: send message to discord
	return nil
}

func (t *Task) Start() {
	go func() {
		for {
			err := t.DrawingTask()
			if err != nil {
				t.logger.Error("Failed in DrawingTask", zap.Error(err))
			}
			time.Sleep(t.Tick)
		}
	}()
}
