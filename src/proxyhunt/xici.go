package proxyhunt

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strconv"
	"log"
	"time"
)

var (
	xiciUrls = []string{"http://www.xici.net.co/nn/1", "http://www.xici.net.co/nn/2", "http://www.xici.net.co/nn/3", "http://www.xici.net.co/wn/1", "http://www.xici.net.co/wn/2"}
)

func GetXiCi() []IP {
	ips := make([]IP, 0)
	for _, link := range xiciUrls {
		log.Println("crawler:", link)
		content := GetContent(link)

		if content == "" {
			continue
		}

		reader := strings.NewReader(content)

		doc , err := goquery.NewDocumentFromReader(reader)

		if err != nil {
			fmt.Println(err)
		}


		doc.Find("table#ip_list > tbody > tr").Each(func(i int, s*goquery.Selection) {
			var ip string
			var port int
			s.Find("td").Each(func(j int, d *goquery.Selection) {
				if j == 2 {
					ip = strings.TrimSpace(d.Text())
				}

				if j == 3 {
					portstr := strings.TrimSpace(d.Text())
					port, _ = strconv.Atoi(portstr)
				}
			})
			key := ip + ":" + strconv.Itoa(port)
			if port != 0 {
				d := IP{Addr:key, Ip:ip, Port:port}
				ips = append(ips, d)
			}
		});
		time.Sleep(5 * time.Second)
	}
	return ips;
}
