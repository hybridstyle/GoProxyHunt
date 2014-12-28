package proxyhunt

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

const (
	url82Format = "http://www.82ip.com/index.asp?stype=1&page=%d"
)

func Get82IP()[]IP{
	ips := make([]IP,0)

	for i:=0;i<5;i++{
		link := fmt.Sprintf(url82Format,i)

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

		docs.Find("div#list > table > tbody > tr").Each(func (i int, s *goquery.Selection){
			var ip string
			var port int

			s.Find("td").Each(func (j int, d *goquery.Selection){
				if j==0 {
					ip=strings.TrimSpace(d.Text())
				}

				if j==1 {
					portstr := strings.TrimSpace(d.Text())

					port,_ = strconv.Atoi(portstr)
				}
			})

			if port != 0 {
				addr := ip +":"+strconv.Itoa(port)
				res := IP{Addr:addr,Ip:ip,Port:port}
				ips=append(ips,res)
			}
		})
	}
	return ips
}
