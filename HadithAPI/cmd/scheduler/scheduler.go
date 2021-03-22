package scheduler

import (
	"HadithAPI/cmd/data"
	"fmt"
	"time"
	"unsafe"
)

//IntervalPeriod between changing of hadiths
const IntervalPeriod time.Duration = 2 * time.Minute

//HourToTick - The set hour for the timer to fire
const HourToTick int = 17

//MinuteToTick - Minute set for the timer to fire
var MinuteToTick int = time.Now().Add(IntervalPeriod).Minute()

//SecondToTick - The second set for the timer to fire
const SecondToTick int = 59

//JobTicker struct to store the timer
type JobTicker struct {
	T *time.Timer
}

func getNextTickDuration() time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), HourToTick, MinuteToTick, SecondToTick, 0, time.Local)
	if nextTick.Before(now) {
		nextTick.Add(IntervalPeriod)
	}
	return nextTick.Sub(time.Now())
}

//NewJobTicker Creates a new JobTicker
func NewJobTicker() JobTicker {
	fmt.Println("new ticker here")
	return JobTicker{time.NewTimer(getNextTickDuration())}
}

//UpdateJobTicker resets the timer
func (jt JobTicker) UpdateJobTicker() {
	jt.T.Reset(getNextTickDuration())
	fmt.Println("ticker updated")
}

//ScheduleHadith is supposed to return a channel which spits only one hadith every after the timer fires
func ScheduleHadith(ch2 <-chan data.CleanHadith) <-chan data.CleanHadith {
	jt := NewJobTicker()
	var s data.CleanHadith
	hadCh := make(chan data.CleanHadith, 50*unsafe.Sizeof(s))
	var in data.CleanHadith
	in = <-ch2
	go func() chan data.CleanHadith {
		for v := range ch2 {
			select {
			case <-jt.T.C:
				hadCh <- v
				jt.UpdateJobTicker()
				return hadCh
			default:
				hadCh <- in
				return hadCh
			}
		}
		close(hadCh)
		return nil
	}()
	return hadCh
}
