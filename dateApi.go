package main

import (
//	"github.com/bt51/ntpclient"
	"time"
//	"log"
	"strconv"
)

func getNtp()(time.Time){
	d:=time.Now().UTC()
	if int(time.Now().UTC().Month()) < 10 {
		d, _ = stringToDate(strconv.Itoa(time.Now().UTC().Day()) + ".0" + strconv.Itoa(int(time.Now().UTC().Month())))
	}else {
		d, _ = stringToDate(strconv.Itoa(time.Now().UTC().Day()) + "." + strconv.Itoa(int(time.Now().UTC().Month())))
	}
	return d
	/*
	t1, err := ntpclient.GetNetworkTime("0.pool.ntp.org", 123)
	if err != nil {
		log.Println(err)
	}
	//log.Println("t1: ",t1)
	t2, err := ntpclient.GetNetworkTime("time.nist.gov", 123)
	if err != nil {
		log.Println(err)
	}
	//log.Println("t2: ",t2)
	if t1 == nil && t2==nil{
		return time.Now().UTC()
	}
	if time.Now().UTC().After(t1.Add(-12*time.Hour)) && time.Now().UTC().Before(t1.Add(12*time.Hour)) && time.Now().UTC().After(t2.Add(-12*time.Hour)) && time.Now().UTC().Before(t2.Add(12*time.Hour)){
		return ()
	} else {
		return *t1
	}*/
}
