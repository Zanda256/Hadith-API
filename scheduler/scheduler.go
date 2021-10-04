package scheduler

import (
	"fmt"
	"time"
)

//IntervalPeriod between changing of hadiths
const IntervalPeriod time.Duration = 1 * time.Minute

//JobTicker struct to store the timer
type JobTicker struct {
	T *time.Ticker
}

var (
	//Jt is the ticker used to schedule new hadith selection
	Jt JobTicker
)

//NewJobTicker Creates a new JobTicker
func NewJobTicker() JobTicker {
	fmt.Println("new ticker here")
	return JobTicker{
		T: time.NewTicker(IntervalPeriod),
	}
}

//UpdateJobTicker resets the timer
func (jt JobTicker) UpdateJobTicker(t time.Time) {
	jt.T.Reset(IntervalPeriod)
	fmt.Println("Ticker Updated")
}

//ScheduleHadith is supposed to return a channel which spits only one hadith every after the timer fires
func ScheduleHadith(jt JobTicker, t time.Time) {
	jt.UpdateJobTicker(t)
	return
}
