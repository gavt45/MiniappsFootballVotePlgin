package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
	"strings"
	"regexp"
)

type Config struct {
	ServerRoot          string
	Port                string
	MainPageXML string
	ResponseXML       string
	ErrorXML string
	DaysXML string
	VoteInputXML string
	DbPath string
	UpdateResultsKey string
	PathToGoogleKeyJson string
}


var config = new(Config)
//var outputFile = new(os.File)
var responseXml = []byte{}
var errorXml = []byte{}
var mainPageXml = []byte{}
var daysXml = []byte{}
var voteInputXml = []byte{}

// added for test commit

//var knownKeys = []string{"ref_sid", "event.id", "event.order", "subscriber", "abonent", "protocol", "user_id", "service", "event.text", "event.referer", "event", "lang", "serviceId", "wnumber"}

func init_system() (*Config, []byte, []byte, []byte, []byte, []byte, error) {
	cfg_bytes, err := ioutil.ReadFile(os.Args[1])
	json.Unmarshal(cfg_bytes, config)
	//log.Println("config: ",config)
	/*
	if !exists("out.csv") {
		ioutil.WriteFile("out.csv", []byte("page,button,user_id,wnumber,protocol\n"), 0644)
	}
	*/
	//f, err := os.OpenFile("out.csv", os.O_APPEND|os.O_WRONLY, 0600)
	resp_xml, err := ioutil.ReadFile(config.ResponseXML)
	main_page_xml, err := ioutil.ReadFile(config.MainPageXML)
	errXml, err := ioutil.ReadFile(config.ErrorXML)
	days_xml, err := ioutil.ReadFile(config.DaysXML)
	vote_input_xml, err := ioutil.ReadFile(config.VoteInputXML)

	if err != nil {
		log.Fatal("Error reading from response files: ", err.Error())
	}

	initialize_sheet()
	initDb()
	return config, resp_xml, errXml, main_page_xml, days_xml, vote_input_xml, err
}

func getMatchesHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	day:=r.URL.Query().Get("day")
	wnumber := r.URL.Query().Get("wnumber")
	intDay, err := stringToDate(day)
	if err!=nil{
		intDay = getNtp().Add(12*time.Hour)
	}
	matches := getMatches(wnumber, intDay)
	//log.Println("Matches: ",matches)
	fmt.Fprintf(w, formMatchesXml(matches))
}

func getDatesHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	//wnumber := r.URL.Query().Get("wnumber")
	//log.Println("Days: ",getNtp().Day())
	dates := removeDublicates(selectAllLargerThen(getNtp().Add(12*time.Hour), getDates()))
	fmt.Fprintf(w, formVoteXml(dates))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	//callback := r.URL.Query().Get("callback") // request should have "callback" parameter
	//if callback == "" {
	fmt.Fprintf(w, string(mainPageXml), config.ServerRoot, config.ServerRoot)
	//} else {
	//	http.Redirect(w, r, callback, 302)
	//}
}

func voteInputHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	log.Println(fmt.Sprintf(string(voteInputXml), config.ServerRoot, r.URL.Query().Get("match"), config.ServerRoot))
	fmt.Fprintf(w, string(voteInputXml), config.ServerRoot, r.URL.Query().Get("match"), config.ServerRoot)
}

func voteHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	wnumber := r.URL.Query().Get("wnumber")
	matchNum := r.URL.Query().Get("match")
	matchNumInt, err := strconv.Atoi(matchNum)
	if err != nil {
		fmt.Fprintf(w, string(errorXml), "Internal error: invalid match num "+matchNum)
		return
	}
	scoreFromUri := r.URL.Query().Get("score")
	//DONE: check regex here
	log.Println("Score: "+scoreFromUri)
	var scoreTemplate = regexp.MustCompile("^\\d{1}\\s{0,1}(:|\\-)\\s{0,1}\\d{1}$")
	log.Println("find string output: "+scoreTemplate.FindString(scoreFromUri))
	if scoreTemplate.FindString(scoreFromUri) != "" {
		//
		scoreFromUri = scoreTemplate.FindString(scoreFromUri)
		score:=[]string{}
		if strings.Contains(scoreFromUri, ":") {
			score = strings.Split(scoreFromUri, ":")
		}else {
			score = strings.Split(scoreFromUri, "-")
		}
		go addVote(matchNum, wnumber, score[0], score[1])
		//fmt.Fprintf(w,string(responseXml),"Thank you for participating!!")
		match:=getMatchById(matchNumInt)
		fmt.Fprintf(w, formVoteRespXml(match.team1+" vs "+match.team2, score[0]+":"+score[1]))
	}else{
		//fmt.Fprintf(w,string(responseXml), "Invalid score format!")
		fmt.Fprintf(w, string(voteInputXml), config.ServerRoot, matchNum, config.ServerRoot)
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	wnumber := r.URL.Query().Get("wnumber")
	fmt.Fprintf(w, formResultXml(getWonResults(wnumber), countUserVotes(wnumber)))
}
func parseUpdErr(err string) (string) {
	out := "Google sheet plugin error: "
	if strings.Contains(err, "404") {
		out += "google sheet not found(404)"
	} else if strings.Contains(err, "403") {
		out += "access denied!(403) You should share your sheet to miniapps@miniappstesterbot.iam.gserviceaccount.com"
	} else {
		out += "unknown google sheets error: " + err
	}
	return out
}
func updateResultHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, "ERROR: %s", "Empty request!")
		return
	}
	key:=r.URL.Query().Get("key")
	if key != config.UpdateResultsKey{
		fmt.Fprintf(w, "ERROR: %s", "Invalid key!")
		return
	}
	spreadsheetId := r.URL.Query().Get("spreadsheetId")
	log.Println("Id: ",spreadsheetId)
	updErr := updSheet(spreadsheetId) // Id of spreadsheet should be passed in "spreadsheetId" parameter
	if updErr != nil {
		fmt.Fprintf(w, parseUpdErr(string(updErr.Error())))
		return
	}
	matches, err := getMatchesFromSheet()
	if err != nil {
		fmt.Fprintf(w, "ERROR: %s", "Error getting matches from sheet:",err)
		return
	}else if len(matches) == 0{
		fmt.Fprintf(w, "ERROR: %s", "Matches length is empty!")
		return
	}
	go updateMatches(matches)
	fmt.Fprintf(w, "Updating...")
}
func main() {
	log.Println("Starting...")
	if len(os.Args) < 2 {
		log.Fatal("You should pass me a config name like: ", os.Args[0], " <json config name>")
	}
	cfg, respXml, errXml, main_page_xml, days_xml, vote_input_xml, err := init_system()
	config = cfg
	errorXml=errXml
	//outputFile = f
	responseXml = respXml
	mainPageXml = main_page_xml
	daysXml=days_xml
	voteInputXml = vote_input_xml
	//log.Println(string(response_xml))
	log.Println("Config: ", config)
	if err != nil {
		//outputFile.Close()
		panic(err)
	}
	log.Println("Done! Listening...")
	http.HandleFunc(config.ServerRoot, mainHandler)
	http.HandleFunc(config.ServerRoot+"days", getDatesHandler)
	http.HandleFunc(config.ServerRoot+"matches", getMatchesHandler)
	http.HandleFunc(config.ServerRoot+"voteInput", voteInputHandler)
	http.HandleFunc(config.ServerRoot+"vote", voteHandler)
	http.HandleFunc(config.ServerRoot+"result", resultHandler)
	http.HandleFunc(config.ServerRoot+"updateResults", updateResultHandler)
	http.ListenAndServe(":"+config.Port, nil)
}