package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"strconv"
)
type Result struct {
	teamA string
	teamB string
}

var database=new(sql.DB)
var initialized = false
func init(){
	database, _ = sql.Open("sqlite3", "C:\\Users\\gav\\go\\src\\MiniappsFootballVotePlugin\\football.db")
	initialized=true
}

func addVote(matchNum string, wnumber string, scoreA string, scoreB string) {
	if !initialized {return}
	log.Println("adding vote: "+"INSERT INTO votes VALUES("+matchNum+", '"+wnumber+"', "+scoreA+", "+scoreB+")")
	stmt, _ := database.Prepare("INSERT INTO votes VALUES("+matchNum+", '"+wnumber+"', "+scoreA+", "+scoreB+")")
	stmt.Exec()
	//statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	//statement.Exec()
	//statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	//statement.Exec("Nic", "Raboy")
	/*
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}*/
}
func getDates()([]int){
	out := []int{}
	if !initialized {return out}
	rows, _ := database.Query("select matches.date from matches")
	var date int
	for rows.Next() {
		rows.Scan(&date)
		log.Println("Date: ",date)
		out=append(out, date)
	}
	return out
}
func getMatches(wnumber string, date int)([]string){
	out := []string{}
	if !initialized {return out}
	log.Println("Query: "+"select com1,com2 from matches,votes where votes.wnumber='"+wnumber+"' and matches.number<>votes.number and matches.date="+strconv.Itoa(date))
	rows, _ := database.Query("select com1,com2 from matches,votes where votes.wnumber='"+wnumber+"' and matches.number<>votes.number and matches.date="+strconv.Itoa(date))
	var com1 string
	var com2 string
	for rows.Next() {
		rows.Scan(&com1, &com2)
		//log.Println("Teams for user "+wnumber+" and date ",date,": ",com1," ",com2)
		out=append(out, com1+" vs "+com2)
	}
	out=removeDouble(out)
	return out
}
func count(val string, arr []string)(int){
	c:=0
	for _,el := range arr{
		if el == val{c+=1}
	}
	return c
}
func contains(val string, arr []string)(bool){
	for _,e:=range arr{
		if val==e{return true}
	}
	return false
}
func removeDouble(in []string)([]string){
	//mapping := new(map[string]int)
	out:=[]string{}
	//get max count
	max:=0
	for _, val := range in{
		//log.Println("Val: ",val,"; Count: ",count(val, in), " In: ",in)
		if count(val,in) > max {
			max=count(val,in)
		}
		/*
		if count(val, in) > 2 && !contains(val,out){
			out=append(out, val)
		}*/
	}
	for _,val := range in {
		if count(val, in) >= max && !contains(val,out){
			out=append(out, val)
		}
	}
	return out
}
func getWonResults(wnumber string)([]Result){
	out := []Result{}
	if !initialized {return out}
	rows, _ := database.Query("select matches.com1,matches.com2,results.scoreA,results.scoreB from matches,votes,results where votes.wnumber='"+wnumber+"' and votes.scoreA=results.scoreA and votes.scoreB=results.scoreB and matches.number=votes.number and matches.number=results.number")
	var com1 string
	var com2 string
	var scoreA int
	var scoreB int
	for rows.Next() {
		rows.Scan(&com1, &com2, &scoreA, &scoreB)
		//log.Println("User "+wnumber+" won at: ",com1," ",com2, " with score ",scoreA,":",scoreB)
		out=append(out, Result{teamA:com1, teamB:com2})
	}
	return out
}


func main() {
	//log.Println(getNtp())
	//addVote("1", "def", "3", "2")
	log.Println(getWonResults("def"))
	/*
	database, _ := sql.Open("sqlite3", "./nraboy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Nic", "Raboy")
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}*/
}

func tests(){main()}
