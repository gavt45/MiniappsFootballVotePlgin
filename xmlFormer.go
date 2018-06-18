package main

import (
	"fmt"
	"strconv"
)

func formVoteXml(dates []int)(string){
	navigation := ""
	for _,day := range dates{
		navigation+=fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot+"matches?day="+strconv.Itoa(day), "june "+strconv.Itoa(day))
	}
	out := fmt.Sprintf(string(daysXml), navigation)
	//log.Println("Vote xml: "+out
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
	//resp := ""
	//if len(results) == 0{
	return fmt.Sprintf(string(responseXml), "❓you voted: "+strconv.Itoa(voted)+"<br/>👍 you guessed: "+strconv.Itoa(len(results)))
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
