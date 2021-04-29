package handlers

import (
	"HadithAPI/cmd/data"
	"HadithAPI/cmd/scheduler"
	"fmt"
	"log"
	"net/http"
)

//Hadiths struct to act as a handler
type Hadiths struct {
	l *log.Logger
}

//HadList is the global pointer to the list from which hadith of the day is selected at random
var hadList *data.Hadiths

//Hadptr is the pointer to the hadith of the day at any given time
var hadptr *data.CleanHadith

//NewHadiths to create a new Hadiths struct handler
func NewHadiths(l *log.Logger) *Hadiths {
	return &Hadiths{l}
}

//function to fetch raw hadiths from json file
func (h *Hadiths) getHadiths() (*data.Hadiths, error) {
	hl, err := data.Fetch()
	if err != nil {
		h.l.Fatal(err)
		return nil, err
	}
	return &hl, nil
}

//Extracts one clean hadith from a channel and returns it
func (h *Hadiths) generate(c chan *data.CleanHadith) *data.CleanHadith {
	var n *data.CleanHadith
	n = <-c
	return n
}

//Initiate method is called in main to initiate the first hadith to be returned (With get) when
//the programis run.
func (h *Hadiths) Initiate() {
	var err error
	hadList, err = h.getHadiths()
	if err != nil {
		fmt.Printf("Could not fetch hadithlist from file %+v", http.StatusInternalServerError)
		return
	}
	hadptr = h.generate(hadList.Gen())
	fmt.Printf("Initial: %+v", *hadptr)
}

func (h *Hadiths) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Set("content-type", "application-json")

		var newOne data.CleanHadith

		hadChano := make(chan data.CleanHadith, 20)

		hadChano <- *hadptr
		fmt.Println(len(hadChano))
		go func() {
			for {
				t := <-scheduler.Jt.T.C
				fmt.Println("We got one in waiting")
				hadptr = h.generate(hadList.Gen())
				//	latestHadChano <- *Hadptr
				fmt.Printf("Latest: %+v", *hadptr)
				scheduler.ScheduleHadith(scheduler.Jt, t)
			}
		}()

		newOne = <-hadChano
		close(hadChano)

		fmt.Println("newOne.ToJSON")
		err := newOne.ToJSON(rw)
		if err != nil {
			http.Error(rw, "ToJSON in handler failed", http.StatusInternalServerError)
			return
		}

		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
	return
}
