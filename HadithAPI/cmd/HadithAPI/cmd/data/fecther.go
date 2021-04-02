package data

import (
	"HadithAPI/cmd/scheduler"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//RawHadith struct to store a hadith extracted from json file
type RawHadith struct {
	Text     string `json:"En_Text"`
	Narrator string `json:"En_Sanad"`
	ID       int    `json:"Hadith_ID"`
}

//Hadiths type to store newly extracted hadiths from the json file
type Hadiths []*RawHadith

//CleanHadith struct to store a hadith after it has been parsed
type CleanHadith struct {
	CText     string `json:"Hadith"`
	CNarrator string `json:"Narrator"`
}

//Fetch function to extract Hadith text and Narrator from json file
func Fetch() (Hadiths, error) {
	hadithList := make([]*RawHadith, 500)
	pw, _ := os.Getwd()
	filename := pw + "/data/hadiths.json"
	hadBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Cannot read json file. %+v", err)
		return nil, err
	}
	err = json.Unmarshal(hadBytes, &hadithList)
	if err != nil {
		log.Printf("Cannot read json file. %+v", err)
		return nil, err
	}
	return hadithList, nil
}

func parseHadith(h *RawHadith) *CleanHadith {

	if strings.Contains(h.Narrator, ":") {
		h.Narrator = strings.ReplaceAll(h.Narrator, ":", ".")
	}
	if strings.ContainsAny(h.Text, "[]") {
		h.Text = strings.ReplaceAll(h.Text, "[", "(")
		h.Text = strings.ReplaceAll(h.Text, "]", ")")
	}
	cl := &CleanHadith{
		CText:     h.Text,
		CNarrator: h.Narrator,
	}
	return cl
}

func randomNoGenerator() int {
	fmt.Println("Inside random")
	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(808)
	return rn
}
func contains(s *Hadiths, id int) (*RawHadith, bool) {
	for _, v := range *s {
		if v.ID == id {
			return v, true
		}
	}
	return nil, false
}

//ToJSON marshals a clean hadith to a json object
func (raw *CleanHadith) ToJSON() (bytes.Buffer, error) {
	var b bytes.Buffer
	hadjson, err := json.Marshal(raw)
	if err != nil {
		fmt.Println("Unable to marshal hadith to json")
		return b, err
	}
	_, err = b.Write(hadjson)
	if err != nil {
		fmt.Println("Unable to write hadith to buffer")
		return b, err
	}
	return b, nil
}

var defaultTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

//Gen method returns a channel containing clean hadiths to be scheduled by the scheduler
func (hl *Hadiths) Gen() <-chan CleanHadith {
	fmt.Println(len(*hl))
	out := make(chan CleanHadith, 5)
	id := randomNoGenerator()
	for i := 0; i < 30; i++ {
		fmt.Println("In the loop")
		if dh, ok := contains(hl, id); ok {
			clhad := parseHadith(dh)
			fmt.Println(*clhad)
			c := *clhad
			out <- c
			fmt.Println("Out from loop")
			return out
		}
		id = randomNoGenerator()
		continue
	}
	fmt.Println("Entering Scheduling...")
	t, ok := scheduler.ScheduleHadith()
	switch {
	case ok:
		if defaultTime.Before(*t) {
			for _, had := range *hl {
				go func(h *RawHadith) chan CleanHadith {
					clhad := parseHadith(h)
					out <- *clhad
					defaultTime = *t
					return out
				}(had)
			}
		}
	default:
		fmt.Println("Hit default case")
	}
	fmt.Println("Scheduling...")

	return out
}
