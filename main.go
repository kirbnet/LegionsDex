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

// Define Global Variables
var checklist Checklist

// Misc Globals
var lightFactions []string = []string{"ARMY OF LEODYSSEUS", "ORDER OF EATHYRON", "CONVOCATION OF BASSYLIA", "XYLONA'S FLOCK"}
var darkFactions []string = []string{"LEGION OF ARETHYR", "CONGREGATION OF NECRONOMINUS", "ILLYTHIA'S BROOD", "CIRCLE OF POXXUS"}
var splinterFactions []string = []string{"SONS OF THE RED STAR", "HOUSE OF THE NOBLE BEAR"}
var goblinRaces []string = []string{"GOBLIN", "GREATER GOBLIN", "SWALE GOBLIN", "WOODLAND GOBLIN (FUZZMUNK)", "GOBLINS"}
var orcRaces []string = []string{"ORC", "HUMAN - HALF-ORC", "LICHEN ORC", "ORAPHIM", "ORC AND HUMAN", "SHADOW ORC", "UUBYR"}
var elfRaces []string = []string{"ELF", "SHADOW ELF", "FAERIE ELF", "ELF - WHISPERLING", "FROST ELF", "WHISPERLING", "WOOD ELF"}
var dwarfRaces []string = []string{"DWARF", "DWARVES", "DWARVEN SKELETON"}
var vampireRaces []string = []string{"VAMPIRE", "UUBYR", "VARGG", "VOGYRR"}
var undeadRaces []string = []string{"SKELETON", "ARAKKIGHAST", "GHOST", "GHOUL", "LICH", "POISON SKELETON", "TURPICULUS", "UMANGEIST", "UNDEAD HORSE", "UNDEAD ANGEL"}
var anthroRaces []string = []string{"AVIAN", "BOARRIOR", "CENTAUR", "DRAGOSYR", "EAGLE", "FAUN", "ELDER FROST DEER", "JAGUALLIAN", "MINOTAUR", "MOOSE", "NORTHLANDS MINOTAUR", "SATYR", "SKORRIAN", "SWALE GOBLIN", "WOODLAND GOBLIN (FUZZMUNK)"}

// Templates
var tpl = template.Must(template.ParseFiles("static/index.html"))
var hometpl = template.Must(template.ParseFiles("static/home.html"))
var detailtpl = template.Must(template.ParseFiles("static/detail.html"))
var drilldowntpl = template.Must(template.ParseFiles("static/drilldown.html"))

// Struct just to hold figures
type Checklist struct {
	Figures []Figure `json:"figures"`
}

// Add a Figure to the Checklist
func (checklist *Checklist) AddItem(figure Figure) {
	checklist.Figures = append(checklist.Figures, figure)
}

// Figure Data
type Figure struct {
	Name    string   `json:"name"`
	Faction string   `json:"faction"`
	Race    string   `json:"race"`
	Role    string   `json:"role"`
	Release []string `json:"released"`
	Url     string   `json:"url"`
	Scale   string   `json:"scale"`
}

// Data for the Home Page
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

// Generic data for a main page list
type ListPageData struct {
	Type       string
	Total      string
	List       map[string]int
	SortedList []string
}

// Data for a single search term, Lists1-3 should correspond to other data types
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
	List4Title string
	List4      map[string]int
}

// Parse JSON data in Figures and Checklist
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

// Sort a Checklist by Figure names
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
	scaleData(checklist)

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
	router.HandleFunc("/races/{races}", racesHandler)
	router.HandleFunc("/faction/", factionDirHandler)
	router.HandleFunc("/faction/{faction}", factionHandler)
	router.HandleFunc("/factions/{factions}", factionsHandler)
	router.HandleFunc("/role/", roleDirHandler)
	router.HandleFunc("/role/{role}", roleHandler)
	router.HandleFunc("/release/", releaseDirHandler)
	router.HandleFunc("/release/{release}", releaseHandler)
	router.HandleFunc("/scale/", scaleDirHandler)
	router.HandleFunc("/scale/{scale}", scaleHandler)

	//Handling Combinations of Requests, stopping at only 2 deep
	router.HandleFunc("/race/{race}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/race/{race}/release/{release}", drilldownHandler)
	router.HandleFunc("/race/{race}/role/{role}", drilldownHandler)
	router.HandleFunc("/race/{race}/scale/{scale}", drilldownHandler)

	router.HandleFunc("/release/{release}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/release/{release}/race/{race}", drilldownHandler)
	router.HandleFunc("/release/{release}/role/{role}", drilldownHandler)
	router.HandleFunc("/release/{release}/scale/{scale}", drilldownHandler)

	router.HandleFunc("/role/{role}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/role/{role}/release/{release}", drilldownHandler)
	router.HandleFunc("/role/{role}/race/{race}", drilldownHandler)
	router.HandleFunc("/role/{role}/scale/{scale}", drilldownHandler)

	router.HandleFunc("/faction/{faction}/race/{race}", drilldownHandler)
	router.HandleFunc("/faction/{faction}/release/{release}", drilldownHandler)
	router.HandleFunc("/faction/{faction}/role/{role}", drilldownHandler)
	router.HandleFunc("/faction/{faction}/scale/{scale}", drilldownHandler)

	router.HandleFunc("/scale/{scale}/race/{race}", drilldownHandler)
	router.HandleFunc("/scale/{scale}/faction/{faction}", drilldownHandler)
	router.HandleFunc("/scale/{scale}/role/{role}", drilldownHandler)
	//Define Static Resources
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
	//Start Port Listener/Web Server
	http.ListenAndServe(":"+port, router)
}

// raceData is a map of the Races with a Count of total instances
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

// checklistByRace creates a new checklist limited to single Race
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

// Gets the Factions with Count from a Checklist
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

// checklistByFaction creates a new checklist limited to single Race
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

// Gets the Roles with count from a Checklist
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

// checklistByRole creates a new checklist limited to single Role
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

// Gets the Releases and Count from a Checklist
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

// checklistByRelease creates a new checklist limited to single Release
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

// Gets the Scales with count from a Checklist
func scaleData(lst Checklist) map[string]int {
	roles := make(map[string]int)
	for i := range lst.Figures {
		_, exists := roles[lst.Figures[i].Scale]
		if exists {
			roles[lst.Figures[i].Scale] += 1
		} else {
			roles[lst.Figures[i].Scale] = 1
		}
	}
	return roles
}

// checklistByScale creates a new checklist limited to single Scale
func checklistByScale(lst Checklist, scale string) Checklist {
	var scaleMembers Checklist
	//iterate through list of figures, and copy those that match
	for _, figure := range lst.Figures {
		if figure.Scale == scale {
			scaleMembers.AddItem(figure)
		}
	}
	return sortChecklist(scaleMembers)
}

// PAGE HANDLER FUNCTIONS
// Main page and default handler.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	releasesOf := releaseData(checklist)
	factionsOf := factionData(checklist)
	racesOf := raceData(checklist)
	rolesOf := roleData(checklist)
	//scalesOf := scaleData(checklist)
	//FACTION FUN
	lightSide := groupSearch(checklist, "faction", lightFactions)
	darkSide := groupSearch(checklist, "faction", darkFactions)
	splinterSide := groupSearch(checklist, "faction", splinterFactions)
	//RACE FUN
	allGoblins := groupSearch(checklist, "race", goblinRaces)
	allOrcs := groupSearch(checklist, "race", orcRaces)
	allElves := groupSearch(checklist, "race", elfRaces)
	allDwarves := groupSearch(checklist, "race", dwarfRaces)
	allVampires := groupSearch(checklist, "race", vampireRaces)
	allSkeletons := groupSearch(checklist, "race", undeadRaces)
	allAnthros := groupSearch(checklist, "race", anthroRaces)

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

// SECTION: FUNCTIONS BY RACE
// Page listing directory of Races
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

// Page listing figures and other data of a specified Race
func raceHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	race := reqvars["race"]
	//Get races from data
	chk := checklistByRace(checklist, race)

	rolesOfRace := roleData(chk)
	factionsOfRace := factionData(chk)
	releasesOfRace := releaseData(chk)
	scalesOfRace := scaleData(chk)

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
	pagedata.List4Title = "scale"
	pagedata.List4 = scalesOfRace

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Page displaying information about figures from several pre-specified Races
func racesHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	races := reqvars["races"]
	//Get races from data
	var group []string
	switch races {
	case "goblin":
		group = goblinRaces
	case "elf":
		group = elfRaces
	case "dwarf":
		group = dwarfRaces
	case "vampire":
		group = vampireRaces
	case "undead":
		group = undeadRaces
	case "anthro":
		group = anthroRaces
	case "orc":
		group = orcRaces
	}
	chk := groupSearch(checklist, "race", group)

	rolesOfRace := roleData(chk)
	//valuekeySortedRolesOfRace := SortMapByValueThenKey(rolesOfRace)
	factionsOfRace := factionData(chk)
	//valuekeySortedFactionsOfRace := SortMapByValueThenKey(factionsOfRace)
	releasesOfRace := releaseData(chk)
	scalesOfRace := scaleData(chk)

	var pagedata DetailPageData
	pagedata.Title = strings.ToTitle(races) + " Race: " + strings.Join(group, ", ")
	pagedata.Type = "race"
	pagedata.Query = strings.Join(group, ", ")
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Checklist = chk
	pagedata.List1Title = "role"
	pagedata.List1 = rolesOfRace
	pagedata.List2Title = "faction"
	pagedata.List2 = factionsOfRace
	pagedata.List3Title = "release"
	pagedata.List3 = releasesOfRace
	pagedata.List4Title = "scale"
	pagedata.List4 = scalesOfRace

	if err := drilldowntpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// SECTION: FUNCTIONS BY FACTION
// Page listing Factions
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

// Page displaying data for Figures of a Faction
func factionHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	faction := reqvars["faction"]
	//Get factions from data
	chk := checklistByFaction(checklist, faction)

	rolesOfFaction := roleData(chk)
	racesOfFaction := raceData(chk)
	releasesOfFaction := releaseData(chk)
	scalesofFaction := scaleData(chk)

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
	pagedata.List4Title = "scale"
	pagedata.List4 = scalesofFaction

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Page displaying information of several pre-specified Factions
func factionsHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	factions := reqvars["factions"]
	//Get factions from data
	var group []string
	switch factions {
	case "light":
		group = lightFactions
	case "dark":
		group = darkFactions
	case "splinter":
		group = splinterFactions
	}
	chk := groupSearch(checklist, "faction", group)
	//chk := checklistByFaction(checklist, faction)

	rolesOfFaction := roleData(chk)
	racesOfFaction := raceData(chk)
	releasesOfFaction := releaseData(chk)
	scalesOfFaction := scaleData(chk)

	var pagedata DetailPageData
	pagedata.Title = strings.ToTitle(factions) + " Factions: " + strings.Join(group, ", ")
	pagedata.Type = "faction"
	pagedata.Query = strings.Join(group, ", ")
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Checklist = chk
	pagedata.List1Title = "role"
	pagedata.List1 = rolesOfFaction
	pagedata.List2Title = "race"
	pagedata.List2 = racesOfFaction
	pagedata.List3Title = "release"
	pagedata.List3 = releasesOfFaction
	pagedata.List4Title = "scale"
	pagedata.List4 = scalesOfFaction

	if err := drilldowntpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Page listing all Roles
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

// Page showing information about figure of a given Role
func roleHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	role := reqvars["role"]
	//Get factions from data
	chk := checklistByRole(checklist, role)

	factionsOfRole := factionData(chk)
	racesOfRole := raceData(chk)
	releasesOfRole := releaseData(chk)
	scalesOfRole := scaleData(chk)

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
	pagedata.List4Title = "scale"
	pagedata.List4 = scalesOfRole

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Page listing all Scales
func scaleDirHandler(w http.ResponseWriter, r *http.Request) {
	//Get scales from data
	scales := scaleData(checklist)
	//Sort
	valuekeySortedScales := SortMapByValueThenKey(scales)
	//Present page
	pagedata := &ListPageData{"scale", strconv.Itoa(len(scales)), scales, valuekeySortedScales}
	if err := tpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Page showing information about figure of a given Role
func scaleHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	scale := reqvars["scale"]
	//Get factions from data
	chk := checklistByScale(checklist, scale)

	factionsOfScale := factionData(chk)
	racesOfScale := raceData(chk)
	rolesOfScale := roleData(chk)
	releasesOfScale := releaseData(chk)

	var pagedata DetailPageData
	pagedata.Total = strconv.Itoa(len(chk.Figures))
	pagedata.Title = strings.ToTitle(scale) + " Scale"
	pagedata.Type = "scale"
	pagedata.Query = scale
	pagedata.Checklist = chk
	pagedata.List1Title = "race"
	pagedata.List1 = racesOfScale
	pagedata.List2Title = "role"
	pagedata.List2 = rolesOfScale
	pagedata.List3Title = "faction"
	pagedata.List3 = factionsOfScale
	pagedata.List4Title = "release"
	pagedata.List4 = releasesOfScale

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Page listing all Releases
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

// Page showing Figure data for a Release
func releaseHandler(w http.ResponseWriter, r *http.Request) {
	//parse request data
	reqvars := mux.Vars(r)
	release := reqvars["release"]
	//Get factions from data
	chk := checklistByRelease(checklist, release)

	//TODO: Get other data sets
	factionsOfRelease := factionData(chk)
	racesOfRelease := raceData(chk)
	rolesOfRelease := roleData(chk)
	scalesOfRelease := scaleData(chk)

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
	pagedata.List4Title = "scale"
	pagedata.List4 = scalesOfRelease

	if err := detailtpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// DRILLDOWN: Searching by 2 parameters
func drilldownHandler(w http.ResponseWriter, r *http.Request) {
	var remainingStats []string
	//parse request data
	reqvars := mux.Vars(r)
	faction := reqvars["faction"]
	race := reqvars["race"]
	release := reqvars["release"]
	role := reqvars["role"]
	scale := reqvars["scale"]

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
	if scale != "" {
		chk = checklistByScale(chk, scale)
		titlePart += strings.ToTitle(role) + " Scale; "
	} else {
		remainingStats = append(remainingStats, "scale")
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
	case "scale":
		pagedata.List2 = racesOf
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
	case "scale":
		pagedata.List3 = factionsOf
	}
	if err := drilldowntpl.Execute(w, pagedata); err != nil {
		fmt.Println(err)
	}
}

// Generic Checklist Search for a group of something
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

// GENERIC SUPPORT FUNCTIONS
// Sorting by keys, returning the ordered slice
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

// Sorting by value, returning an ordered slice
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

// Sorts a map by the value, then key, returning an ordered slice
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
