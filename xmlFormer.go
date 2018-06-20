package main

import (
	"fmt"
	"strconv"
	"time"
	"log"
)

func formVoteXml(dates []time.Time)(string){
	navigation := ""
	for _,day := range dates{
		if int(day.Month()) < 10 {
			navigation += fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot+"matches?day="+dateToString(day), strconv.Itoa(day.Day())+".0"+strconv.Itoa(int(day.Month())))
		}else {
			navigation += fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot+"matches?day="+dateToString(day), strconv.Itoa(day.Day())+"."+strconv.Itoa(int(day.Month())))
		}
	}
	out := fmt.Sprintf(string(daysXml), navigation)
	log.Println("Vote xml: "+out)
	return out
}

func formMatchesXml(matches []Match)(string){
	navigation := ""
	for _,match := range matches{
		navigation+=fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot+"voteInput?match="+strconv.Itoa(match.idx), match.team1 + " vs "+ match.team2)
	}

	if len(matches) == 0 {
		navigation=fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot, "No matches to vote. Go to start.")
	}else {
		navigation+=fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot, "Main menu")
	}
	out := fmt.Sprintf(string(daysXml), navigation)
	//log.Println("Vote xml: "+out)
	return out
}
func formResultXml(results []Result, voted int)(string){
	resp := "<page version=\"2.0\"><div>%s</div><navigation><link pageId=\"%sdays\">vote</link><link pageId=\"%sresult\">view results</link><link pageId=\"%sdetailedResult\">view detailed results</link></navigation></page>"
	//if len(results) == 0{
	return fmt.Sprintf(string(resp), "❓you voted: "+strconv.Itoa(voted)+"<br/>👍 you guessed: "+strconv.Itoa(len(results)), config.ServerRoot, config.ServerRoot, config.ServerRoot)
	//}
	/*
	resp="You was right at:<br/>"
	for _,result := range results{
		resp+=result.teamA+" vs "+result.teamB+"<br/>"
	}*/
	//out := fmt.Sprintf(string(responseXml), resp)
	//log.Println("Result xml: "+out)
	//return out
}

func formDetailedResultXml(results []DetailedResult)(string){
	resp := "<page version=\"2.0\"><div>Your votes:<br/>"
	for _, result := range results{
		if result.realKnown {
			resp += result.team1 + " vs " + result.team2 + " " + strconv.Itoa(result.vscoreA) + ":" + strconv.Itoa(result.vscoreB) + " real result " + strconv.Itoa(result.rscoreA) + ":" + strconv.Itoa(result.rscoreB)+"<br/>"
		}else{
			resp += result.team1 + " vs " + result.team2 + " " + strconv.Itoa(result.vscoreA) + ":" + strconv.Itoa(result.vscoreB) + " real result unknown<br/>"
		}
	}
	resp+="</div><navigation><link pageId=\"%sdays\">vote</link><link pageId=\"%sr\">view results</link></navigation></page>"
	return fmt.Sprintf(resp, config.ServerRoot, config.ServerRoot)
}


func formVoteRespXml(match string, score string)(string){
	resp := "<page version=\"2.0\"><div>%s</div><navigation><link pageId=\"%sdays\">vote</link><link pageId=\"%sresult\">view results</link></navigation></page>"
	return fmt.Sprintf(string(resp), "Thank you! Your forecast is accepted. "+match+" "+score, config.ServerRoot, config.ServerRoot)
}
