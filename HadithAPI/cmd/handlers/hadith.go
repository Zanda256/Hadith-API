package handlers

import (
	"HadithAPI/cmd/data"
	"fmt"
	"log"
	"net/http"
)

//Hadiths struct to act as a handler
type Hadiths struct {
	l *log.Logger
}

//NewHadiths to create a new Hadiths struct handler
func NewHadiths(l *log.Logger) *Hadiths {
	return &Hadiths{l}
}

func (h *Hadiths) getHadiths() (*data.Hadiths, error) {
	hl, err := data.Fetch()
	if err != nil {
		h.l.Fatal(err)
		return nil, err
	}
	return &hl, nil
}
func (h *Hadiths) generate() <-chan data.CleanHadith {
	hl, err := h.getHadiths()
	if err != nil {
		fmt.Println(err)
	}
	hodch := hl.Gen()
	fmt.Println(<-hodch)
	return hodch
}

func (h *Hadiths) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rw.Header().Set("content-type", "application-json")
		newch := h.generate()
		new := <-newch
		fmt.Println("Entering JSON encoder")
		err := new.ToJSON(rw)
		if err != nil {
			http.Error(rw, "ToJSON in handler failed", http.StatusInternalServerError)
		}
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
