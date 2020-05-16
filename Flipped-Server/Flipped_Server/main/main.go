package main

import (
	"Flipped_Server/dataBase"
	"Flipped_Server/sqlmapper"
)

func main() {
	//httpServer := network.HttpServer{"", 8080, iris.New()}
	//httpServer.Run()
	//dataBase.Init()
	res := [3]byte{1,2,3}
	rtb := dataBase.RegisterTable{
		Pid: "id",
		Username:"name",
		Password:"pwd",
		UserType:0 ,
		Email:"mail"   ,
		Photo:res[:],
		RealName:"realname",
		Profession:"pf",
		Age:1      ,
		Region:"region"  ,
		Hobby:"hobby",
	}
	sqlmapper.Insert(rtb, "im")
}

