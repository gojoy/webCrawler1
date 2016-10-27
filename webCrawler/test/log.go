package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	//"io"
	//"bytes"
	//"github.com/PuerkitoBio/goquery"
	"webCrawler/base"
)

func getLogger(x string) *log.Logger {
	logger:=log.New(os.Stdout,x,log.Lshortfile)
	return logger
}


func NewHttpClient() *http.Client {
	return &http.Client{}
}
func xmain() {
	logr:=getLogger(" ")
	srartUrl:="http://www.qq.com"
	req,err:=http.NewRequest("GET",srartUrl,nil)
	if err != nil {
		logr.Fatalln(err)
	}
	cli:=NewHttpClient()
	resp,err:=cli.Do(req)
	if err != nil {
		logr.Println(err)
	}

	para:=base.GenParase()
	datalist,errs:=para(resp,2)
	if errs != nil {
		logr.Println(errs)
	}
	fmt.Println("begin data\n")
	for _,v:=range datalist{
		fmt.Print(v.Valid())
	}




}
