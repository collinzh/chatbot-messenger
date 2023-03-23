package util

import "time"

type ScheduledTask struct {
	stop chan bool
}

func (s *ScheduledTask) Stop() {
	s.stop <- true
	close(s.stop)
}

func ScheduleTask(runnable func(), delay time.Duration) *ScheduledTask {
	stop := make(chan bool)
	go func() {
		runnable()
		for true {
			select {
			case <-stop:
				return
			case <-time.After(delay):
				runnable()
			}
		}
	}()

	return &ScheduledTask{stop: stop}
}
