package data

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unsafe"
)

//RawHadith struct to store a hadith extracted from json file
type RawHadith struct {
	Text     string `json:"En_Text"`
	Narrator string `json:"En_Sanad"`
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
		log.Fatalf("Cannot read json file. %s", err)
		return nil, err
	}
	json.Unmarshal(hadBytes, &hadithList)
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

//ToJSON marshals a clean hadith to a json object
func (raw *CleanHadith) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(raw)
}

//Gen method returns a channel containing clean hadiths to be scheduled by the scheduler
func (hl *Hadiths) Gen() <-chan CleanHadith {
	var s CleanHadith
	out := make(chan CleanHadith, 50*unsafe.Sizeof(s))
	go func() chan CleanHadith {
		for _, had := range *hl {
			clhad := parseHadith(had)
			out <- *clhad
		}
		close(out)
		return out
	}()
	return out
}
