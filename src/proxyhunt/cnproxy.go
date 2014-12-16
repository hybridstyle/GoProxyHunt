package proxyhunt

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
)

const (
	cnProxyUrl = "http://cn-proxy.com/"
)

func GetCnProxy() []IP {
	ips := make([]IP, 0)

	content := GetContent(cnProxyUrl)

	if content == "" {
		return ips
	}

	reader := strings.NewReader(content)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return ips
	}

	doc.Find("table.sortable > tbody > tr").Each(func(i int, s *goquery.Selection) {
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

			if port != 0 {
				addr := ip + ":" + strconv.Itoa(port)
				res := IP{Addr:addr, Ip:ip, Port:port}
				ips = append(ips, res)
			}

		})
	})

	return ips
}
