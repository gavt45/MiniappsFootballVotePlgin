package main

func selectAllLargerThen(val int, arr []int)([]int){
	o:=[]int{}
	for _,e := range arr{
		if e > val {
			o=append(o, e)
		}
	}
	return o
}

func intContains(val int, arr []int)(bool){
	for _,e:=range arr{
		if val==e{return true}
	}
	return false
}

func removeDublicates(in []int)([]int){
	o:=[]int{}
	for _,e := range in{
		if !intContains(e,o){
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
