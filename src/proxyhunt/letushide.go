package proxyhunt

import (
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

const (
	letusurlFormat = "http://letushide.com/filter/http,hap,all/%d/list_of_free_HTTP_High_Anonymity_proxy_servers"
	letuscount     = 2
)

func GetLetUsShide() []IP {
	ips := make([]IP, 0)

	for i := 1; i <= letuscount; i++ {
		link := fmt.Sprintf(letusurlFormat, i)

		content := GetContent(link)

		if content == "" {
			continue
		}

		reader := strings.NewReader(content)
		doc , err := goquery.NewDocumentFromReader(reader)

		if err != nil {
			continue
		}

		doc.Find("table#basic > tbody > tr").Each(func(i int, s *goquery.Selection) {
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
	}
	return ips;
}


