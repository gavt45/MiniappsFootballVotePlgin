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
	//log.Println("Vote xml: "+out)
	return out
}

func formMatchesXml(matches []Match)(string){
	navigation := ""
	for _,match := range matches{
		navigation+=fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot+"voteInput?match="+strconv.Itoa(match.idx), match.name)
	}
	if len(matches) == 0 {
		navigation=fmt.Sprintf("<link pageId=\"%s\">%s</link>\n", config.ServerRoot, "No matches to vote. Go to start.")
	}
	out := fmt.Sprintf(string(daysXml), navigation)
	//log.Println("Vote xml: "+out)
	return out
}
func formResultXml(results []Result)(string){
	resp := ""
	if len(results) == 0{
		return fmt.Sprintf(string(responseXml), "You was not right in any of matches!")
	}
	resp="You was right at:<br/>"
	for _,result := range results{
		resp+=result.teamA+" vs "+result.teamB+"<br/>"
	}
	out := fmt.Sprintf(string(responseXml), resp)
	//log.Println("Result xml: "+out)
	return out
}
