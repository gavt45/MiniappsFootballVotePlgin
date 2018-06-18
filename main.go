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

	//initialize_sheet()
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
	intDay, err := strconv.Atoi(day)
	if err!=nil{
		intDay = getNtp().Add(12*time.Hour).Day()
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
	dates := removeDublicates(selectAllLargerThen(getNtp().Day(), getDates()))
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
	//log.Println(fmt.Sprintf(string(voteInputXml), config.ServerRoot, r.URL.Query().Get("match")))
	fmt.Fprintf(w, string(voteInputXml), config.ServerRoot, r.URL.Query().Get("match"))
}

func voteHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	wnumber := r.URL.Query().Get("wnumber")
	matchNum := r.URL.Query().Get("match")
	scoreFromUri := r.URL.Query().Get("score")
	//DONE: check regex here
	if res, _ := regexp.MatchString("^([0-9]+):([0-9]+)$", scoreFromUri); res {
		//
		score := strings.Split(scoreFromUri, ":")
		go addVote(matchNum, wnumber, score[0], score[1])
		fmt.Fprintf(w,string(responseXml),"Thank you for participating!!")
	}else{
		fmt.Fprintf(w,string(responseXml), "Invalid score format!")
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request){
	log.Println("Got request:", r.URL.String(), "\nContent: ", r.Body)
	if len(r.URL.Query()) == 0 {
		fmt.Fprintf(w, string(errorXml), "Empty request!")
		return
	}
	wnumber := r.URL.Query().Get("wnumber")
	fmt.Fprintf(w, formResultXml(getWonResults(wnumber)))
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
	http.ListenAndServe(":"+config.Port, nil)
}