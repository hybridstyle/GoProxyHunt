package main

import "fmt"

import (
	"proxyhunt"
	"strings"
	"time"
	"net/url"
	"net/http"
	"net"
	"crypto/tls"
	"strconv"
	"os"
	"log"
	"bufio"
	"github.com/go-martini/martini"
	"encoding/json"
)

const (
	checkUrl string    = "http://www.baidu.com"
	bws string         = "BWS"
	filename string = "proxydata"
	cleaninternal int = 5 * 60 * 1000 // 5 minute
	saveinternal time.Duration = 8 * 60 * 1000 * time.Millisecond //5 min
	dataformat string = "%s\t%s\t%d\t%d\t%d"
)
var ipmap = make(map[string]proxyhunt.IP)
var queuemap = make(map[string]int)
var checkqueue = make(chan string,10)


func main() {

	loadProxy()//加载代理数据

	go pachongWorker()
	go checkWorker()
	go saveWorker()
	go cleanWorkder()

	m := martini.Classic()

	m.Get("/ips",IpJson)

	m.RunOnAddr(":10086")
}

func IpJson()string{
	res := make(map[string]interface{})
	arr := make([]proxyhunt.IP,0)
	count:=0
	for _,v := range ipmap {
		count++
		arr=append(arr,v)
	}
	res["count"]=count
	res["ips"]=arr

	tmp ,err :=json.Marshal(res)
	if err != nil {
		return "{error}"
	}
	return string(tmp)
}

func loadProxy(){
	currentPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	datafilepath := currentPath+"/"+filename

	datafile, err := os.Open(datafilepath)

	if err != nil {
		log.Println(err)
	}
	defer datafile.Close()

	scanner := bufio.NewScanner(datafile)

	for scanner.Scan() {
		var addr,ip string
		var port int
		var ctime,utime int64
		fmt.Sscanf(scanner.Text(),dataformat,&addr,&ip,&port,&ctime,&utime)
		p := proxyhunt.IP{Addr:addr,Ip:ip,Port:port,Ctime:ctime,Utime:utime}
		ipmap[p.Addr]=p
	}


}

/*
chan[] proxy to be checked
map[IP]int64 map saved avaiable proxy
10 mins save to file
 */
func saveWorker(){
	for {
		time.Sleep(saveinternal)
		saveToFile()

	}
}

func saveToFile(){
	currentPath ,err :=os.Getwd()

	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println(currentPath)
	datafilepath := currentPath+"/"+filename
	file, err := os.OpenFile(datafilepath,os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC,0666)

	if err != nil{
		log.Fatal(err)
	}
	defer file.Close();
	for _,v := range ipmap{
		val := fmt.Sprintf(dataformat,v.Addr,v.Ip,v.Port,v.Ctime,v.Utime)
		fmt.Println(val)
		file.Write([]byte(val+"\n"))
	}
	fmt.Println("\n")
}

/*
get the timeout proxy to check
 */
func cleanWorkder(){
	for {
		time.Sleep(1*time.Minute)
		now := Now()
		for _,v := range ipmap {
			internal := int(now - v.Utime)

			if internal > cleaninternal {
				_,ok := queuemap[v.Addr]
				if !ok {
					checkqueue<-v.Addr
				}
			}
		}
	}
}

/*
get proxy from chan and check
 */
func checkWorker(){
	for {
		select {
		case proxy:= <-checkqueue:
			if checkProxy(proxy,10){
				now := Now()
				v,ok := ipmap[proxy]

				if ok {
					v.Utime=now
				}else{
					ip,port:=getIpInfo(proxy)
					v=proxyhunt.IP{Addr:proxy,Ip:ip,Port:port,Ctime:now,Utime:now}
				}
				ipmap[proxy]=v

			}else{
				delete(ipmap,proxy)
			}
		}
	}
}

func pachongWorker(){
	for {
		ips := proxyhunt.GetPaChong()
//		ips := proxyhunt.GetList("http://pachong.org/area/city/name/上海/type/high.html")
		fmt.Println("pachong size:",len(ips))
		for _,ip := range ips {
			checkqueue<-ip.Addr
		}
		time.Sleep(1*time.Hour)
	}

}

func getIpInfo(proxy string)(string,int){
	arr := strings.Split(proxy,":")
	port,_ := strconv.Atoi(arr[1])
	return arr[0],port
}



func checkProxy(proxy string, timeout int) bool {
	proxyUrl, err := url.Parse("http://" + proxy)
	httpClient := &http.Client{
		Transport:&http.Transport{
			Proxy:http.ProxyURL(proxyUrl),
			Dial:func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Duration(timeout) * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout))
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
			DisableKeepAlives: true,
			//			DisableCompression: true,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	req, err := http.NewRequest("HEAD", checkUrl, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.94 Safari/537.36")

	req.Close = true

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("checkproxy error %s %s\n", proxy, err)
		return false
	}
	defer resp.Body.Close()

	//	ioutil.ReadAll(resp.Body)

	headers := resp.Header["Server"]
	if headers != nil {
		serverHeader := headers[0]
		if strings.Contains(serverHeader, bws) {
			return true
		}
	}

	return false
}


func Now() int64{
	return time.Now().UnixNano() / int64(time.Millisecond)
}

