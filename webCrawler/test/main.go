package main

import (
	"time"
	"fmt"
	//"bytes"
	"log"
	"os"
)

var chan1 chan int
var chanlen int=9
var interval time.Duration=1500*time.Millisecond

func getchan() chan int{
	return chan1
}

func receive(ch chan int)  {
	fmt.Printf("begin to recieve element from chan...\n")
	timer:=time.After(16*time.Second)
	//var f bool=false
	loop:
	for {
		select {
		case e,ok:=<-getchan():
			if !ok {
				fmt.Println("chan closed in receive\n")
				//f=true
				break loop
			}
			fmt.Printf("recieve element is %v\n",e)
			time.Sleep(interval)
		case <-timer:
			fmt.Printf("timeout!\n")
			//f=true
			break loop
		}

	}
}
func ExampleLogger() {
	//var buf bytes.Buffer
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	logger.Print("Hello, log file!")

	//fmt.Println(&buf)

	// Output:
	// logger: example_test.go:16: Hello, log file!
}

func main() {
	chan1=make(chan int,chanlen)
	ExampleLogger()
	go func() {
		for i:=0;i<chanlen;i++ {
			if i > 0&&i % 3 == 0 {
				fmt.Printf("Reset chan1\n")
				chan1=make(chan int,chanlen)
			}
			fmt.Printf("send element...%d\n",i)
			chan1<-i
			time.Sleep(interval)
		}
		fmt.Printf("close chan...\n")
		close(chan1)
	}()

	receive(chan1)
	//fmt.Printf("now end is %s!\n",tm)
}
