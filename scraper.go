package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var figureChecklist Checklist

const CHECKLIST string = "https://sourcehorsemen.com/mythic-legions/checklist"
const BASEURL string = "https://sourcehorsemen.com"

type Checklist struct {
	Figures []Figure `json:"figures"`
}
type Figure struct {
	Name    string `json:"name"`
	Faction string `json:"faction"`
	Race    string `json:"race"`
	Role    string `json:"role"`
	Release string `json:"released"`
	Url     string `json:"url"`
}

func getChecklist(checkURL string) {
	//var links []string

	response, err := http.Get(checkURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		//Use goQuery for page parsing
		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		document.Find("ul.character-list li").Each(func(index int, selector *goquery.Selection) {
			link, _ := selector.Find("a").Attr("href")
			getFigureData(BASEURL + link)
		})
	}
}

func getFigureData(figURL string) {
	//Check that page is alive
	response, err := http.Get(figURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		//Use goQuery for page parsing
		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fig := Figure{}

		document.Find("div.billboardMessage").Each(func(index int, selector *goquery.Selection) {
			figureName := selector.Find("h2").Text()
			fig.Name = strings.TrimSpace(figureName)
		})
		document.Find("div.character-stats li").Each(func(index int, selector *goquery.Selection) {
			statText := selector.Text()
			stat := strings.Split(statText, ":")
			//fmt.Println(strconv.Itoa(index) + stat[0] + " " + stat[1])
			switch stat[0] {
			case "Race":
				fig.Race = strings.ToUpper(strings.TrimSpace(stat[1]))
			case "Faction":
				fig.Faction = strings.ToUpper(strings.TrimSpace(stat[1]))
			case "Role":
				fig.Role = strings.ToUpper(strings.TrimSpace(stat[1]))
			case "Released In":
				releaseRaw := strings.ToUpper(strings.TrimSpace(stat[1]))
				var releases []string
				if strings.Contains(releaseRaw, ",") {
					releases = strings.Split(releaseRaw, ", ")
				} else {
					releases = append(releases, releaseRaw)
				}
				fig.Release = releases
				//case "Accessories"
				//case "Additional Heads"
			}

		})
		fig.Url = figURL
		figureChecklist.Figures = append(figureChecklist.Figures, fig)
	}
}

//Resolve issues with data
func zhuzh(s string) string {
	var newStr string
	//TODO: Need to remove or encode slashes "/", specifically for "N/A"
	//TODO: Need to remove or change question marks "?" - Change to unknown?
	//TODO: Releases are a mess. Commas separate when a figure was released more than once.
	//TODO: Releases "Soul Spiller." Period before comment about other release
	if strings.Contains(s, "Soul Spiller") {
		newStr = "SOUL SPILLER"
	}
	//TODO: Same with "Wasteland"
	return newStr
}
func main() {
	getChecklist(CHECKLIST)
	fmt.Println(figureChecklist)
	//encoding
	checkByte, err := json.Marshal(figureChecklist)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("figurechecklist.json", checkByte, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(checkByte))
}
