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

type Match struct {
	idx int
	name string
}

var database=new(sql.DB)
var initialized = false
func initDb(){
	database, _ = sql.Open("sqlite3", config.DbPath)
	initialized=true
}

func addVote(matchNum string, wnumber string, scoreA string, scoreB string) {
	if !initialized {return}
	//log.Println("adding vote: "+"INSERT INTO votes VALUES("+matchNum+", '"+wnumber+"', "+scoreA+", "+scoreB+")")
	stmt, err := database.Prepare("INSERT INTO votes VALUES("+matchNum+", '"+wnumber+"', "+scoreA+", "+scoreB+")")
	log.Println("Db api err: ",err)
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
	rows, err := database.Query("select matches.date from matches")
	log.Println("Db api err: ",err)
	var date int
	for rows.Next() {
		rows.Scan(&date)
		//log.Println("Date: ",date)
		out=append(out, date)
	}
	return out
}
func countVotedMatches(wnumber string, date int)(int){
	rows, err := database.Query("select matches.number,matches.com1,matches.com2 from matches,votes where  matches.date="+strconv.Itoa(date)+" and votes.wnumber='"+wnumber+"' and votes.number=matches.number")
	log.Println("Db api err: ",err)
	var i = 0
	var idx int
	var com1 string
	var com2 string
	for rows.Next() {
		rows.Scan(&idx, &com1, &com2)
		i+=1
	}
	//log.Println("voted matches: ",i)
	return i
	}

func getAllMatchesByDay(day int)([]Match){
	out:=[]Match{}
	rows, err := database.Query("select matches.number,matches.com1,matches.com2 from matches where date="+strconv.Itoa(day))//"select distinct matches.number,matches.com1,matches.com2 from matches,votes where  matches.date="+strconv.Itoa(date)+" and votes.wnumber='"+wnumber+"' except select matches.number,matches.com1,matches.com2 from matches,votes where  matches.date="+strconv.Itoa(date)+" and votes.wnumber='"+wnumber+"' and votes.number=matches.number")
	log.Println("Db api err: ",err)
	var idx int
	var com1 string
	var com2 string
	for rows.Next() {
		rows.Scan(&idx, &com1, &com2)
		//log.Println("Teams for user "+wnumber+" and date ",date,": ",com1," ",com2)
		out=append(out, Match{idx:idx, name:com1+" vs "+com2})
	}
	return out
}

func getMatches(wnumber string, date int)([]Match){
	out := []Match{}
	if !initialized {return out}
	if countVotedMatches(wnumber,date) == 0{
		return getAllMatchesByDay(date)
	}
	//log.Println("Query: "+"select matches.number,com1,com2 from matches,votes where votes.wnumber='"+wnumber+"' and matches.number<>votes.number and matches.date="+strconv.Itoa(date))
	rows, err := database.Query("select distinct matches.number,matches.com1,matches.com2 from matches,votes where  matches.date="+strconv.Itoa(date)+" and votes.wnumber='"+wnumber+"' except select matches.number,matches.com1,matches.com2 from matches,votes where  matches.date="+strconv.Itoa(date)+" and votes.wnumber='"+wnumber+"' and votes.number=matches.number")
	log.Println("Db api err: ",err)
	var idx int
	var com1 string
	var com2 string
	for rows.Next() {
		rows.Scan(&idx, &com1, &com2)
		//log.Println("Teams for user "+wnumber+" and date ",date,": ",com1," ",com2)
		out=append(out, Match{idx:idx, name:com1+" vs "+com2})
	}
	//log.Println("Out before remove double: ",out)
	return out
	/*
	}else {
		out=[]Match{}
		rows, _  = database.Query("SELECT matches.number,com1,com2 from matches where matches.date="+strconv.Itoa(date))
		var idx int
		var com1 string
		var com2 string
		for rows.Next() {
			rows.Scan(&idx, &com1, &com2)
			//log.Println("Teams for user "+wnumber+" and date ",date,": ",com1," ",com2)
			out=append(out, Match{idx:idx, name:com1+" vs "+com2})
		}
		log.Println("out: ",out)
		return out
	}*/
}
/*
func count(val string, arr []Match)(int){
	c:=0
	for _,el := range arr{
		if el.name == val{c+=1}
	}
	return c
}
func contains(val string, arr []Match)(bool){
	for _,e:=range arr{
		if val==e.name{return true}
	}
	return false
}

func removeDouble(in []Match)([]Match){
	//mapping := new(map[string]int)
	out:=[]Match{}
	//get max count
	max:=0
	counts := []int{}
	log.Println("remove double in: ",in)
	for _, val := range in{
		//log.Println("Val: ",val,"; Count: ",count(val, in), " In: ",in)
		counts=append(counts, count(val.name,in))
		if count(val.name,in) > max {
			max=count(val.name,in)

		}
		/*
		if count(val, in) > 2 && !contains(val,out){
			out=append(out, val)
		}*
	}
	log.Println(counts)
	if !allEquals(counts) {
		for _, val := range in {
			if count(val.name, in) >= max && !contains(val.name, out) {
				out = append(out, val)
			}
		}
		return out
	}else {
		return []Match{}
	}
}*/
func getWonResults(wnumber string)([]Result){
	out := []Result{}
	if !initialized {return out}
	rows, err := database.Query("select matches.com1,matches.com2,results.scoreA,results.scoreB from matches,votes,results where votes.wnumber='"+wnumber+"' and votes.scoreA=results.scoreA and votes.scoreB=results.scoreB and matches.number=votes.number and matches.number=results.number")
	log.Println("Db api err: ",err)
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