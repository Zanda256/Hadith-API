package scheduler

import (
	"fmt"
	"time"
)

//IntervalPeriod between changing of hadiths
const IntervalPeriod time.Duration = 1 * time.Minute

var defaultTime time.Time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

//JobTicker struct to store the timer
type JobTicker struct {
	T *time.Ticker
}

var Jt JobTicker

//NewJobTicker Creates a new JobTicker
func NewJobTicker() JobTicker {
	fmt.Println("new ticker here")
	return JobTicker{
		T: time.NewTicker(IntervalPeriod),
	}
}

//UpdateJobTicker resets the timer
func (jt JobTicker) UpdateJobTicker() {
	jt.T.Reset(IntervalPeriod)
	fmt.Println("Ticker Updated")
}

//ScheduleHadith is supposed to return a channel which spits only one hadith every after the timer fires
func ScheduleHadith() (*time.Time, bool) {
	var Jt = NewJobTicker()
	value, more := <-Jt.T.C
	fmt.Println(value)
	if more {
		fmt.Println("still more inside")
		Jt.UpdateJobTicker()
		return &value, true
	}
	fmt.Println("Nothing here")
	return nil, false
}
