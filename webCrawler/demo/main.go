package main

import (
	"fmt"
	"webCrawler/base"
	"webCrawler/scheduler"
)

func startSchdule() {
	channelArgs:=base.NewChannelArgs(10,10,10,10)
	poolArgs:=base.NewPoolBaseArgs(3,3)
	crawDepth:=uint32(2)

}

func main() {
	fmt.Printf("this is main\n")
}
