package proxyhunt

import (
	"net"
	"net/http"
	"time"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func Get(url string) *http.Response {
	timeout := 20
	httpClient := &http.Client {
		Transport:&http.Transport{
		Dial:func(netw, addr string) (net.Conn, error) {
			deadline := time.Now().Add(time.Duration(timeout) * time.Second)
			c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout))
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
			},
		ResponseHeaderTimeout:time.Duration(timeout) * time.Second,
			// DisableKeepAlives:true,
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},
		TLSHandshakeTimeout:10 * time.Second,
		},
	}

	req, err := http.NewRequest("GET",url,nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Mozilla/5.0ss (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.94 Safari/537.36")
	req.Close = true

	resp, err := httpClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp
}

func ParseContent(url string) string {
	resp := Get(url)

	if resp == nil {
		return ""
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	jsContent,_:=doc.Find("#js_content").Html()

	jsContent = strings.Replace(jsContent,"data-src","src",-1)
	return jsContent
}

func GetContent(url string)string {
	resp := Get(url)

	if resp == nil {
		return ""
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)

}
