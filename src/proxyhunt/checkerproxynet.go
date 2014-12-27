package proxyhunt

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"log"
)

const (
	checkerProxyNetUrl = "http://checkerproxy.net/all_proxy"
)

func GetCheckerProxyNet() []IP {
	ips := make([]IP, 0)

	content := GetContent(checkerProxyNetUrl)

	if content == "" {
		return ips
	}

	reader := strings.NewReader(content)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return ips
	}

	size :=0
	doc.Find("td.proxy-ipport").Each(func(i int, s *goquery.Selection) {
		size++
		addr := strings.TrimSpace(s.Text())

		ipinfo := strings.Split(addr, ":")

		ipstr := ipinfo[0]
		portstr := ipinfo[1]

		port, err := strconv.Atoi(portstr)
		if err == nil {
			res := IP{Addr:addr, Ip:ipstr, Port:port}
			ips = append(ips, res)
		}else{
			log.Println(addr,err)
			//TODO:解析此种情况 117.165.167.135:8123 (China)
		}
	})

	log.Println("size:",size)

	return ips

}
