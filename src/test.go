package main

import (
	"proxyhunt"
	"fmt"
)

func main() {
	ips := proxyhunt.GetCheckerProxyNet()

	fmt.Println(len(ips))

//	for _,ip := range ips {
//		fmt.Printf("%s\n",ip.Addr)
//	}
}

