package main

import "time"

func selectAllLargerThen(val time.Time, arr []time.Time)([]time.Time){
	o:=[]time.Time{}
	for _,e := range arr{
		if e.After(val) {
			o=append(o, e)
		}
	}
	return o
}

func timeContains(val time.Time, arr []time.Time)(bool){
	for _,e:=range arr{
		if val==e{return true}
	}
	return false
}

func removeDublicates(in []time.Time)([]time.Time){
	o:=[]time.Time{}
	for _,e := range in{
		if !timeContains(e,o){
			o=append(o,e)
		}
	}
	return o
}

func allEquals(in []int)(bool){
	if len(in) == 0 {return false}
	std:=in[0]
	for _,a := range in{
		if a!=std{
			return false
		}else{
			std=a
		}
	}
	return true
}

func dateToString(t time.Time)(string) {
	return t.Format("2006-01-_2@15:04:05")
}

func stringToDate(date string)(time.Time, error){
	//layout := "DD-MM-YYYY hh:mm:ss"
	//return time.Parse(layout, date)
	return time.Parse("2006-01-_2@15:04:05", date)
}
