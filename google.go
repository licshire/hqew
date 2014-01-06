package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	//var url string = "https://www.google.com.hk/search?num=100&newwindow=1&safe=strict&site=&source=hp&q=site%3Ahqew.com+ic.html+%E8%AF%A2%E4%BB%B7&oq=site%3Ahqew.com+ic.html+%E8%AF%A2%E4%BB%B7&gs_l=hp.3...4851.13830.0.14017.31.29.2.0.0.0.124.1515.28j1.29.0....0...1c.1j4.32.hp..26.5.391.PwcfajwGLH8"
	//var url string = "https://www.google.com.hk/search?q=site:hqew.com+ic.html+%E8%AF%A2%E4%BB%B7&num=100&newwindow=1&safe=strict&ei=O-fIUs7UCMTJiAeD-4DADQ&start=100&sa=N&biw=1096&bih=576"
	var url string = "https://www.google.com.hk/search?q=site:hqew.com+ic.html+%E8%AF%A2%E4%BB%B7&num=100&newwindow=1&safe=strict&ei=lQLJUujhH4TsiAfgkYHADA&start=200&sa=N&biw=1096&bih=576"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println("connect error")
	}
	bodyByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Println("connect error")
	}
	res.Body.Close()
	body := string(bodyByte)
	if len(body) != 0 {
		//fmt.Println(body)
	}
	//fmt.Println(body)

	f_hqew_list, _ := os.OpenFile("hqew_list_all.txt", os.O_APPEND, 0666)
	w_hqew_list := bufio.NewWriter(f_hqew_list)

	onePartSlice := strings.Split(body, `<div class="kv" style="margin-bottom:2px">`)
	//onePartSlice := strings.Split(body, "ic.html")
	for _, onePart := range onePartSlice {
		fmt.Println("\r\n\r\n\r\n\r\n\r\n", onePart)

		secondDomainExp := regexp.MustCompile(`<cite>(.*).hqew.com/(.*)</cite><`)
		secondDomainSlice := secondDomainExp.FindStringSubmatch(onePart)
		if len(secondDomainSlice) >= 3 {
			secondDomain := strings.Replace(secondDomainSlice[1], "www.", "", -1)
			fmt.Println(secondDomain)
			if isGet(secondDomain) {
				line_hqew_list := strings.Join([]string{secondDomain, time.Now().Format("2006-01-02 23:12"), "\r\n"}, "\t")
				w_hqew_list.WriteString(line_hqew_list)
			}

		}
		time.Sleep(5000 * time.Millisecond)
	}

	w_hqew_list.Flush()
	defer f_hqew_list.Close()

	//fmt.Println(stocksSlice)
}

func isGet(secondDomain string) bool {
	//根据url首先判断是否读取
	url := strings.Join([]string{"http://", secondDomain, ".hqew.com/ic/ic.html?Page=1"}, "")
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println("connect error")
	}
	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Println("connect error")
	}
	res.Body.Close()
	body := string(robots)
	if len(body) == 0 {
		fmt.Println("null body")
	}

	//判断该页面是否包含型号列表
	isMatch, _ := regexp.MatchString(
		`<tr class="tr0"><td class="c1".*</td></tr>`, body)
	if isMatch {
		return true
	} else {
		return false
	}

}
