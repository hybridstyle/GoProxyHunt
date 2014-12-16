package proxyhunt

import (
	"fmt"
	"time"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

const (
	urlFormat = "http://proxy.com.ru/gaoni/list_%d.html"
	count     = 40
)

func GetProxyCom() []IP {
	ips := make([]IP, 0)
	for i := 1; i <= count ; i++ {

		link := fmt.Sprintf(urlFormat, i)
		content := GetContent(link);

		if content == "" {
			continue
		}
		if content == "test" {
			continue
		}

		reader := strings.NewReader(content)
		docs, err := goquery.NewDocumentFromReader(reader)

		if err != nil {
			continue
		}

		docs.Find("center > font > table > tbody > tr > td > font > table > tbody > tr").Each(func(i int, s*goquery.Selection) {
			var ip string
			var port int
			s.Find("td").Each(func(j int, d *goquery.Selection) {
				if j == 1 {
					ip = strings.TrimSpace(d.Text())
				}

				if j == 2 {
					portstr := strings.TrimSpace(d.Text())
					port, _ = strconv.Atoi(portstr)
				}
			})
			if port != 0 {
				addr := ip + ":" + strconv.Itoa(port)
				res := IP{Addr:addr, Ip:ip, Port:port}
				ips = append(ips, res)
			}

		})

		time.Sleep(1 * time.Second)
	}
	return ips
}
