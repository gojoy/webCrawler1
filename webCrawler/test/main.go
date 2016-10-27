package main

import (
	"time"
	"fmt"

	"log"
	"os"
	"flag"

	"runtime"
	"strings"
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

func ltmain() {
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
func put(s... string)  {
	fmt.Printf("in put this is %s\n",s)
}
var infile *string = flag.String("i", "input11", "input context")
var outfile *string=flag.String("o","output22","output context")
var name string
var fileok=flag.Bool("v",false,"test bool")
const vc="vxc"
/*
func init() {
	//s:=[]string{"t","-","t","x"}
	cmd:=flag.NewFlagSet("t",flag.ExitOnError)
	var testcmd string
	cmd.StringVar(&testcmd,"t","default","we test")
	flag.StringVar(&name,"n","gg","this is name")
	cmd.Parse(os.Args[:])
	fmt.Printf("testcmd is %v\n",testcmd)
	fmt.Printf("rest is %v\n",cmd.Args())
}
*/
func main()  {
	var poto []string
	poto=[]string{"tcp://127.0.0.1:3300","unix:///var/run/docker.sock"}
	for _,v:=range poto{
		sp:=strings.SplitN(v,"://",2)
		fmt.Println(len(sp),sp[0],sp[1])
	}

	fmt.Printf("os is %v,%v\n",runtime.GOOS,runtime.GOARCH)
	fmt.Printf("type is %T\n",vc)
	ch:=make(chan int,2)
	go func(ch chan int) {
		time.Sleep(1*time.Second)
		fmt.Println("in gor")
		ch<-10
		time.Sleep(1*time.Second)
		fmt.Printf("ater 1s\n")
		ch<-20
		close(ch)
	}(ch)
	for x:=range ch{
		fmt.Printf("rece a x %d\n",x)
	}
	if v, ok := <-ch; ok {
		fmt.Printf("v is %v\n",v)
	}else {
		fmt.Printf("ch is close")
	}
	/*
	fmt.Printf("the value infile is %T\n",name)
	flag.Parse()
	fmt.Printf("after inout name is %v\n",name)
	if infile!=nil {
		fmt.Printf("the inout is %s\n",*infile)
	}
	if outfile!=nil {
		fmt.Printf("the output is %v\n", *outfile)
	}
	if *fileok {
		fmt.Printf("we test ok\n")
	}else {
		fmt.Printf("no bool\n")
	}
	fmt.Printf("the words is %v\n",len(flag.Args()))
*/

}
