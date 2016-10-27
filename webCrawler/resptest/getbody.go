package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/PuerkitoBio/goquery"
)

func getLogger(x string) *log.Logger {
	logger := log.New(os.Stdout, x, log.Lshortfile)
	return logger
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}
func main() {
	logr := getLogger(" ")
	startUrl3:="http://studygolang.com/topics"
	//startUrl2:="https://www.oschina.net/question/tag/javascript"
	//startUrl1:="https://www.qcloud.com/doc/product/236/3188"
	req, err := http.NewRequest("GET", startUrl3, nil)
	if err != nil {
		logr.Fatalln(err)
	}
	cli := NewHttpClient()
	resp, err := cli.Do(req)
	if err != nil {
		logr.Println(err)
	}

	/*
	bod:=make([]byte,0)

	for  {
		b := make([]byte, 1)
		n,err:=resp.Body.Read(b)
		if n==0 {
			break
		}
		if err != nil {
			logr.Fatal(err)
		}
		bod=append(bod,b[0])
	}
	fmt.Println(len(bod))
	for i := 0; i < 1000; i++ {
		fmt.Println(bod[i])
	}
*/
	fmt.Printf("status is %v,status code is %v\n",resp.Status,resp.StatusCode)
	if resp.StatusCode!=200 {
		logr.Fatal("status error")
	}
	doc,err:=goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logr.Fatal(err.Error())
	}
	fmt.Println(doc.Find("*").Text())
	/*
	doc.Find(".topics .topic").Each(func(i int, contentSelection *goquery.Selection) {
		title := contentSelection.Find(".title a").Text()
		log.Println("第", i+1, "个帖子的标题：", title)
	})
	*/
	//fmt.Printf("doc is %v\n",t)
}
