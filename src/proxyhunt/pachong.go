package proxyhunt

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"strconv"
	"time"
)
/*

proxy data from http://pachong.org/

 */
const (
	PACHONG="http://pachong.org"
	START_URL="http://pachong.org/area/short/name/cn/type/high.html"
)

func PaChongWorker(){
//	ips := GetPaChong()
//
//	for _,ip := range ips {
//
//	}
}

func GetPaChong()[]IP{

	content := GetContent(START_URL)

	reader := strings.NewReader(content)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		fmt.Println(err)
	}
	links := make([]string,0)
	doc.Find("div.natWap > table > tbody > tr > td > a").Each(func(i int,s *goquery.Selection){
		val,exist :=s.Attr("href")
		if exist {
			if !strings.Contains(val,"http") && val != ""{
//				fmt.Println(s.Text(),":",PACHONG+val)
				links=append(links,PACHONG+val)
			}

		}
	});
	ips := make([]IP,0)
	for _,url := range links {
		tmp := GetList(url)
		time.Sleep(5*time.Second)
		ips=append(ips,tmp...)
	}
	return ips
}

func GetList(url string) []IP{
//	fmt.Println("crawler ",url)
	content :=GetContent(url)

	start := strings.Index(content,"<script type=\"text/javascript\">")


	if start == -1 {
		log.Println("get content err")
		return nil
	}

	scriptlen := len("<script type=\"text/javascript\">")
	tmpcontent := content[start+scriptlen:]
	end := strings.Index(tmpcontent,"</script>")

	calVar := tmpcontent[0:end]

	calcDict := calcVariable(calVar)

	reader := strings.NewReader(content)

	doc,err :=goquery.NewDocumentFromReader(reader)

	if err != nil {
		fmt.Println(err)
	}

	ips := make([]IP,0)

	doc.Find("table.tb > tbody > tr").Each(func(i int,s *goquery.Selection){
		var ip string
		var port int
		s.Find("td").Each(func(j int,d *goquery.Selection){

			if j==1 {
				ip = strings.TrimSpace(d.Text());
			}
			if j==2 {
				portstr := strings.TrimSpace(d.Text())
				port = calcPort(portstr,calcDict)
			}
		})
		key := ip+":"+strconv.Itoa(port)
		if ip != "" {
			d := IP{Addr:key,Ip:ip,Port:port}
			ips = append(ips,d)
		}
	})

	return ips
}

func calcPort(portstr string,dict map[string]int)int {

	portstr=strings.Replace(portstr,"(","",-1)
	portstr=strings.Replace(portstr,")","",-1)
	portstr=strings.Replace(portstr,"document.write","",-1)
	portstr=strings.Replace(portstr,";","",-1)

	for k,v := range dict {
		portstr=strings.Replace(portstr,k,strconv.Itoa(v),-1)
	}

	index := strings.Index(portstr,"^")

	one,_ := strconv.Atoi(portstr[0:index])

	portstr = portstr[index+1:]

	index = strings.Index(portstr,"+")

	two,_ := strconv.Atoi(portstr[0:index])

	three,_:=strconv.Atoi(portstr[index+1:])

	return (one^two)+three
}

func calcVariable(val string)map[string]int{

	val =strings.Replace(val,"var ","",-1)
	arr := strings.Split(val,";")

	dict := make(map[string]int)

	for i:=0;i<len(arr);i++{
		tmp := arr[i]
		if tmp == "" {
			continue;
		}
		tmparr := strings.Split(tmp,"=")
		key := tmparr[0]
		val := tmparr[1]
		for k,v := range dict {
			vstr := strconv.Itoa(v)
			val=strings.Replace(val,k,vstr,-1)
		}
		calVal :=add(val)
		dict[key]=calVal
	}
	return dict;
}
func add(expr string)int {
	addIndex := strings.Index(expr,"+")
	prefix := expr[0:addIndex]
	prefixInt,_ := strconv.Atoi(prefix)
	suffix := expr[addIndex+1:]
	xorindex := strings.Index(suffix,"^")

	if xorindex != -1 {
		a,b := xor(suffix)
		return prefixInt+a^b
	}else{
		suffixInt,_ := strconv.Atoi(suffix)
		return prefixInt+suffixInt
	}
}

func xor(expr string)(int,int) {
	xorindex := strings.Index(expr,"^")
	prefix := expr[0:xorindex]
	prefixInt,_ := strconv.Atoi(prefix)
	suffix := expr[xorindex+1:]
	suffixInt,_ := strconv.Atoi(suffix)
	return prefixInt,suffixInt
}
