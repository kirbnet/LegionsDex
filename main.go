package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//Define Global Variables
var checklist Checklist

//Templates
var tpl = template.Must(template.ParseFiles("static/index.html"))
var hometpl = template.Must(template.ParseFiles("static/home.html"))
var detailtpl = template.Must(template.ParseFiles("static/detail.html"))
var drilldowntpl = template.Must(template.ParseFiles("static/drilldown.html"))

type Checklist struct {
	Figures []Figure `json:"figures"`
}

func (checklist *Checklist) AddItem(figure Figure) {
	checklist.Figures = append(checklist.Figures, figure)
}

type Figure struct {
	Name    string   `json:"name"`
	Faction string   `json:"faction"`
	Race    string   `json:"race"`
	Role    string   `json:"role"`
	Release []string `json:"released"`
	Url     string   `json:"url"`
}
type HomePageData struct {
	RaceTotal     int
	RoleTotal     int
	FactionTotal  int
	ReleaseTotal  int
	FigureTotal   int
	LightTotal    int
	DarkTotal     int
	SplinterTotal int
	GoblinTotal   int
	OrcTotal      int
	ElfTotal      int
	UndeadTotal   int
	DwarfTotal    int
	VampireTotal  int
	AnthroTotal   int
}
type ListPageData struct {
	Type       string
	Total      string
	List       map[string]int
	SortedList []string
}
type DetailPageData struct {
	Title      string
	Type       string
	Query      string
	Total      string
	Checklist  Checklist
	List1Title string
	List1      map[string]int
	List2Title string
	List2      map[string]int
	List3Title string
	List3      map[string]int
}

func loadDatabase() {
	db, err := ioutil.ReadFile("figurechecklist.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(db, &checklist)
	if err != nil {
		log.Fatal(err)
	}
}
func sortChecklist(lst Checklist) Checklist {
	sortedChecklist := lst
	sort.Slice(sortedChecklist.Figures, func(i, j int) bool {
		return sortedChecklist.Figures[i].Name < sortedChecklist.Figures[j].Name
	})
	return sortedChecklist
}

func main() {
	loadDatabase()
	raceData(checklist)
	factionData(checklist)
	roleData(checklist)
	releaseData(checklist)

	//Page Server
	//If there is a preconfigured port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//Mux Http Handler
	router := mux.NewRouter()
	//Request handlers
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/race/", raceDirHandler)
	router.HandleFunc("/race/{race}", raceHandler)
	router.HandleFunc("/faction/", factionDirHandler)
	router.HandleFunc("/faction/{faction}", factionHandler)
	router.HandleFunc("/role/", roleDirHandler)
	router.HandleFunc("/role/{role}", roleHandler)
	router.HandleFunc("/release/", releaseDirHandler)
	router.HandleFunc("/release/{release}", releaseHandler)

	//Handling Combinations of Requests, stopping at only 2 deep
	router.HandleFunc("/race/{race}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/race/{race}/release/{release}", drilldownHandler)
	router.HandleFunc("/race/{race}/role/{role}", drilldownHandler)

	router.HandleFunc("/release/{release}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/release/{release}/race/{race}", drilldownHandler)
	router.HandleFunc("/release/{release}/role/{role}", drilldownHandler)

	router.HandleFunc("/role/{role}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/role/{role}/release/{release}", drilldownHandler)
	router.HandleFunc("/role/{role}/race/{race}", drilldownHandler)

	router.HandleFunc("/faction/{faction}/race/{race}", drilldownHandler)
	router.HandleFunc("/faction/{faction}/release/{release}", drilldownHandler)
	router.HandleFunc("/faction/{faction}/role/{role}", drilldownHandler)

	//Define Static Resources
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
	//Start Port Listener/Web Server
	http.ListenAndServe(":"+port, router)
}

//SECTION: FUNCTIONS BY RACE
//raceData is a map of the Races with a Count of total instances
func raceData(lst Checklist) map[string]int {
	races := make(map[string]int)
	for i := range lst.Figures {
		_, exists := races[lst.Figures[i].Race]
		if exists {
			races[lst.Figures[i].Race] += 1
		} else {
			races[lst.Figures[i].Race] = 1
		}
	}
	return races
}

//checklistByRace creates a new checklist limited to single Race
func checklistByRace(lst Checklist, race string) Checklist {
	var raceMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range lst.Figures {
		if figure.Race == race {
			raceMembers.AddItem(figure)
		}
	}
	return sortChecklist(raceMembers)
}

//SECTION: FUNCTIONS BY FACTION
func factionData(lst Checklist) map[string]int {
	factions := make(map[string]int)
	for i := range lst.Figures {
		_, exists := factions[lst.Figures[i].Faction]
		if exists {
			factions[lst.Figures[i].Faction] += 1
		} else {
			factions[lst.Figures[i].Faction] = 1
		}
	}
	return factions
}

//checklistByFaction creates a new checklist limited to single Race
func checklistByFaction(lst Checklist, faction string) Checklist {
	var factionMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range lst.Figures {
		if figure.Faction == faction {
			factionMembers.AddItem(figure)
		}
	}
	return sortChecklist(factionMembers)
}
func roleData(lst Checklist) map[string]int {
	roles := make(map[string]int)
	for i := range lst.Figures {
		_, exists := roles[lst.Figures[i].Role]
		if exists {
			roles[lst.Figures[i].Role] += 1
		} else {
			roles[lst.Figures[i].Role] = 1
		}
	}
	return roles
}

//checklistByRole creates a new checklist limited to single Role
func checklistByRole(lst Checklist, role string) Checklist {
	var roleMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range lst.Figures {
		if figure.Role == role {
			roleMembers.AddItem(figure)
		}
	}
	return sortChecklist(roleMembers)
}
func releaseData(lst Checklist) map[string]int {
	releases := make(map[string]int)
	for i := range lst.Figures {
		for j := range lst.Figures[i].Release {
			_, exists := releases[lst.Figures[i].Release[j]]
			if exists {
				releases[lst.Figures[i].Release[j]] += 1
			} else {
				releases[lst.Figures[i].Release[j]] = 1
			}
		}
	}
	return releases
}

//checklistByRelease creates a new checklist limited to single Release
func checklistByRelease(lst Checklist, release string) Checklist {
	var releaseMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range lst.Figures {
		for _, r := range figure.Release {
			if r == release {
				releaseMembers.AddItem(figure)
			}
		}

	}
	return sortChecklist(releaseMembers)
}

//PAGE HANDLER FUNCTIONS
//Main page and default handler.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	releasesOf := releaseData(checklist)
	factionsOf := factionData(checklist)
	racesOf := raceData(checklist)
	rolesOf := roleData(checklist)
	//FACTION FUN

	lightSide := forcesOfLight(checklist)
	darkSide := forcesOfDarkness(checklist)
	splinterSide := forcesOfSplinter(checklist)
	//RACE FUN

	allGoblins := groupSearch(checklist, "race", []string{"GOBLIN", "GREATER GOBLIN", "SWALE GOBLIN", "WOODLAND GOBLIN (FUZZMUNK)"})
	allOrcs := groupSearch(checklist, "race", []string{"ORC", "HUMAN - HALF-ORC", "LICHEN ORC", "ORAPHIM", "ORC AND HUMAN", "SHADOW ORC", "UUBYR"})
	allElves := groupSearch(checklist, "race", []string{"ELF", "SHADOW ELF", "FAERIE ELF", "ELF - WHISPERLING", "FROST ELF", "WHISPERLING", "WOOD ELF"})
	allDwarves := groupSearch(checklist, "race", []string{"DWARF"})
	allVampires := groupSearch(checklist, "race", []string{"VAMPIRE", "UUBYR", "VARGG", "VOGYRR"})
	allSkeletons := groupSearch(checklist, "race", []string{"SKELETON", "ARAKKIGHAST", "GHOST", "GHOUL", "LICH", "POISON SKELETON", "TURPICULUS", "UMANGEIST", "UNDEAD HORSE"})
	allAnthros := groupSearch(checklist, "race", []string{"AVIAN", "BOARRIOR", "CENTAUR", "DRAGOSYR", "EAGLE", "FAUN", "ELDER FROST DEER", "JAGUALLIAN", "MINOTAUR", "MOOSE", "SATYR", "SWALE GOBLIN", "WOODLAND GOBLIN (FUZZMUNK)"})

	var pagedata HomePageData
	pagedata.FactionTotal = len(factionsOf)
	pagedata.FigureTotal = len(checklist.Figures)
	pagedata.RaceTotal = len(racesOf)
	pagedata.RoleTotal = len(rolesOf)
	pagedata.ReleaseTotal = len(releasesOf)
	pagedata.LightTotal = len(lightSide.Figures)
	pagedata.DarkTotal = len(darkSide.Figures)
	pagedata.SplinterTotal = len(splinterSide.Figures)
	pagedata.GoblinTotal = len(allGoblins.Figures)
	pagedata.OrcTotal = len(allOrcs.Figures)
	pagedata.ElfTotal = len(allElves.Figures)
	pagedata.AnthroTotal = len(allAnthros.Figures)
	pagedata.DwarfTotal = len(allDwarves.Figures)
	pagedata.VampireTotal = len(allVampires.Figures)
	pagedata.UndeadTotal = len(allSkeletons.Figures)

	if err := hometpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

//SECTION: FUNCTIONS BY RACE
func raceDirHandler(w http.ResponseWriter, r *http.Request) {
	//Get races from data
	races := raceData(checklist)
	//Sort races for display
	//sortedRaces := SortMapByKeys(races)
	//valueSortedRaces := SortMapByValue(races)
	valuekeySortedRaces := SortMapByValueThenKey(races)

	pagedata := &ListPageData{"race", strconv.Itoa(len(races)), races, valuekeySortedRaces}
	if err := tpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}
func raceHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	race := reqvars["race"]
	//Get races from data
	chk := checklistByRace(checklist, race)

	rolesOfRace := roleData(chk)
	//valuekeySortedRolesOfRace := SortMapByValueThenKey(rolesOfRace)
	factionsOfRace := factionData(chk)
	//valuekeySortedFactionsOfRace := SortMapByValueThenKey(factionsOfRace)
	releasesOfRace := releaseData(chk)
	//valuekeySortedReleasesOfRace := SortMapByValueThenKey(releasesOfRace)

	var pagedata DetailPageData
	pagedata.Title = strings.ToTitle(race) + " Race"
	pagedata.Type = "race"
	pagedata.Query = race
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Checklist = chk
	pagedata.List1Title = "role"
	pagedata.List1 = rolesOfRace
	pagedata.List2Title = "faction"
	pagedata.List2 = factionsOfRace
	pagedata.List3Title = "release"
	pagedata.List3 = releasesOfRace

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

//SECTION: FUNCTIONS BY FACTION
func factionDirHandler(w http.ResponseWriter, r *http.Request) {
	//Get factions from data
	factions := factionData(checklist)
	//Sort
	valuekeySortedFactions := SortMapByValueThenKey(factions)
	//Present page
	pagedata := &ListPageData{"faction", strconv.Itoa(len(factions)), factions, valuekeySortedFactions}
	if err := tpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}
func factionHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	faction := reqvars["faction"]
	//Get factions from data
	chk := checklistByFaction(checklist, faction)

	rolesOfFaction := roleData(chk)
	//valuekeySortedRolesOfRace := SortMapByValueThenKey(rolesOfRace)
	racesOfFaction := raceData(chk)
	//valuekeySortedFactionsOfRace := SortMapByValueThenKey(factionsOfRace)
	releasesOfFaction := releaseData(chk)
	//valuekeySortedReleasesOfRace := SortMapByValueThenKey(releasesOfRace)

	var pagedata DetailPageData
	pagedata.Title = strings.ToTitle(faction) + " Faction"
	pagedata.Type = "faction"
	pagedata.Query = faction
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Checklist = chk
	pagedata.List1Title = "role"
	pagedata.List1 = rolesOfFaction
	pagedata.List2Title = "race"
	pagedata.List2 = racesOfFaction
	pagedata.List3Title = "release"
	pagedata.List3 = releasesOfFaction

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

func forcesOfLight(chk Checklist) Checklist {
	//define factions
	lightFactions := []string{"ARMY OF LEODYSSEUS", "ORDER OF EATHYRON", "CONVOCATION OF BASSYLIA", "XYLONA'S FLOCK"}
	return factionGroup(chk, lightFactions)
}
func forcesOfDarkness(chk Checklist) Checklist {
	//define factions
	darkFactions := []string{"LEGION OF ARETHYR", "CONGREGATION OF NECRONOMINUS", "ILLYTHIA'S BROOD", "CIRCLE OF POXXUS"}
	return factionGroup(chk, darkFactions)
}
func forcesOfSplinter(chk Checklist) Checklist {
	//define factions
	splinterFactions := []string{"SONS OF THE RED STAR", "HOUSE OF THE NOBLE BEAR"}
	return factionGroup(chk, splinterFactions)
}
func factionGroup(chk Checklist, facs []string) Checklist {
	var factionMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range chk.Figures {
		//iterate through members of group
		for _, faction := range facs {
			if figure.Faction == faction {
				factionMembers.AddItem(figure)
			}
		}
	}
	return sortChecklist(factionMembers)
}
func roleDirHandler(w http.ResponseWriter, r *http.Request) {
	//Get roles from data
	roles := roleData(checklist)
	//Sort
	valuekeySortedRoles := SortMapByValueThenKey(roles)
	//Present page
	pagedata := &ListPageData{"role", strconv.Itoa(len(roles)), roles, valuekeySortedRoles}
	if err := tpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}
func roleHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	role := reqvars["role"]
	//Get factions from data
	chk := checklistByRole(checklist, role)

	//TODO: Get other data sets
	factionsOfRole := factionData(chk)
	//valuekeySortedRolesOfRace := SortMapByValueThenKey(rolesOfRace)
	racesOfRole := raceData(chk)
	//valuekeySortedFactionsOfRace := SortMapByValueThenKey(factionsOfRace)
	releasesOfRole := releaseData(chk)
	//valuekeySortedReleasesOfRace := SortMapByValueThenKey(releasesOfRace)

	var pagedata DetailPageData
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Title = strings.ToTitle(role) + " Role"
	pagedata.Type = "role"
	pagedata.Query = role
	pagedata.Checklist = chk
	pagedata.List1Title = "race"
	pagedata.List1 = racesOfRole
	pagedata.List2Title = "faction"
	pagedata.List2 = factionsOfRole
	pagedata.List3Title = "release"
	pagedata.List3 = releasesOfRole
	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}
func releaseDirHandler(w http.ResponseWriter, r *http.Request) {
	//Get releases from data
	releases := releaseData(checklist)
	//Sort
	valuekeySortedReleases := SortMapByValueThenKey(releases)
	//Present page
	pagedata := &ListPageData{"release", strconv.Itoa(len(releases)), releases, valuekeySortedReleases}
	if err := tpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}
func releaseHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	release := reqvars["release"]
	//Get factions from data
	chk := checklistByRelease(checklist, release)

	//TODO: Get other data sets
	factionsOfRelease := factionData(chk)
	//valuekeySortedRolesOfRace := SortMapByValueThenKey(rolesOfRace)
	racesOfRelease := raceData(chk)
	//valuekeySortedFactionsOfRace := SortMapByValueThenKey(factionsOfRace)
	rolesOfRelease := roleData(chk)
	//valuekeySortedReleasesOfRace := SortMapByValueThenKey(releasesOfRace)

	var pagedata DetailPageData
	pagedata.Title = strings.ToTitle(release) + " Release"
	pagedata.Type = "release"
	pagedata.Query = release
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Checklist = chk
	pagedata.List1Title = "race"
	pagedata.List1 = racesOfRelease
	pagedata.List2Title = "role"
	pagedata.List2 = rolesOfRelease
	pagedata.List3Title = "faction"
	pagedata.List3 = factionsOfRelease
	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

//DRILLDOWN: Searching by 2 parameters
func drilldownHandler(w http.ResponseWriter, r *http.Request) {
	var remainingStats []string
	//parse request data
	reqvars := mux.Vars(r)
	faction := reqvars["faction"]
	race := reqvars["race"]
	release := reqvars["release"]
	role := reqvars["role"]

	//Create new checklists
	chk := checklist
	titlePart := "Drilldown: "
	if faction != "" {
		chk = checklistByFaction(chk, faction)
		titlePart += strings.ToTitle(faction) + " Faction; "
	} else {
		remainingStats = append(remainingStats, "faction")
	}
	if race != "" {
		chk = checklistByRace(chk, race)
		titlePart += strings.ToTitle(race) + " Race; "
	} else {
		remainingStats = append(remainingStats, "race")
	}

	if release != "" {
		chk = checklistByRelease(chk, release)
		titlePart += strings.ToTitle(release) + " Release; "
	} else {
		remainingStats = append(remainingStats, "release")
	}
	if role != "" {
		chk = checklistByRole(chk, role)
		titlePart += strings.ToTitle(role) + " Role; "
	} else {
		remainingStats = append(remainingStats, "role")
	}

	releasesOf := releaseData(chk)
	factionsOf := factionData(chk)
	racesOf := raceData(chk)
	rolesOf := roleData(chk)

	var pagedata DetailPageData
	pagedata.Title = titlePart
	pagedata.Type = "drilldown"
	pagedata.Query = r.URL.Path
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Checklist = chk

	//TODO: Straighten out which tertiary lists are displayed

	//pagedata.List1Title = "race"
	//pagedata.List1 = racesOf
	pagedata.List2Title = remainingStats[0]
	switch remainingStats[0] {
	case "role":
		pagedata.List2 = rolesOf
	case "race":
		pagedata.List2 = racesOf
	case "release":
		pagedata.List2 = releasesOf
	case "faction":
		pagedata.List2 = factionsOf
	}
	pagedata.List3Title = remainingStats[1]
	switch remainingStats[1] {
	case "role":
		pagedata.List3 = rolesOf
	case "race":
		pagedata.List3 = racesOf
	case "release":
		pagedata.List3 = releasesOf
	case "faction":
		pagedata.List3 = factionsOf
	}
	if err := drilldowntpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

//
func groupSearch(chk Checklist, searchType string, matches []string) Checklist {
	var newMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range chk.Figures {
		//iterate through members of group
		for _, q := range matches {
			switch searchType {
			case "faction":
				if figure.Faction == q {
					newMembers.AddItem(figure)
				}
			case "race":
				if figure.Race == q {
					newMembers.AddItem(figure)
				}
			case "role":
				if figure.Role == q {
					newMembers.AddItem(figure)
				}
				/* case "release":
				if figure.Release == q {
					newMembers.AddItem(figure)
				}*/
			}
		}

	}
	return sortChecklist(newMembers)
}

//GENERIC SUPPORT FUNCTIONS
func SortMapByKeys(m map[string]int) []string {
	//First, make a slice of just the keys, which can be sorted
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}
	//Sorting in Ascending Order
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func SortMapByValue(m map[string]int) []string {
	//First, make a slice of just the keys, which can be sorted
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}
	//Sorting in Descending Order
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}

func SortMapByValueThenKey(m map[string]int) []string {
	//First, make a slice of just the keys, which can be sorted
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}
	//Sorting in Descending Order
	sort.SliceStable(keys, func(i, j int) bool {
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}
