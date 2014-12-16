package proxyhunt

import (
	"log"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

const (
	domain = "http://free-proxy-list.net/"
)

func GetFreeProxyList() []IP {
	content := GetContent(domain)

	ips := make([]IP, 0)

	if content == "" {
		log.Println("freeproxylist empty")

		return ips
	}

	reader := strings.NewReader(content)

	doc , err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Println(err)
		return ips
	}

	doc.Find("table#proxylisttable > tbody > tr").Each(func(i int, s*goquery.Selection) {
		if strings.Contains(s.Text(), "elite") {
			var ip string
			var port int
			s.Find("td").Each(func(j int, d *goquery.Selection) {
				if j == 0 {
					ip = strings.TrimSpace(d.Text())
				}

				if j == 1 {
					portstr := strings.TrimSpace(d.Text())
					port, _ = strconv.Atoi(portstr)
				}
			})
			if port != 0 {
				key := ip + ":" + strconv.Itoa(port)
				res := IP{Addr:key, Ip:ip, Port:port}
				ips = append(ips, res)
			}
		}
	})
	return ips

}
