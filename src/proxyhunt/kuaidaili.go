package proxyhunt

import (
	"strconv"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

const (
	domainUrl = "http://www.kuaidaili.com/proxylist/"
)

func GetKuaiDaili() []IP {

	ips := make([]IP, 0)
	for i := 1; i <= 10; i++ {
		link := domainUrl + strconv.Itoa(i)
		content := GetContent(link)

		if content == "" {
			continue
		}

		reader := strings.NewReader(content)

		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Println(err)
			continue
		}

		doc.Find("div#list > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "高匿名") {
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
					addr := ip + ":" + strconv.Itoa(port)
					res := IP{Addr:addr, Ip:ip, Port:port}
					ips = append(ips, res)
				}
			}
		})
		time.Sleep(5 * time.Second)
	}

	return ips
}
