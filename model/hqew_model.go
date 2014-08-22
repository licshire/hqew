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

//爬虫hqew网站, 根据型号判断是否有现货
//从file中读取型号列表models.txt
// .\hqew_model.exe
func main() {
	dat, _ := ioutil.ReadFile("models.txt")

	modelSlice := strings.Split(string(dat), "\r\n")
	//fmt.Print(secondDomainSlice)
	old_data, _ := ioutil.ReadFile("model.html")
	f_hqew_list, _ := os.OpenFile("model.html", os.O_RDWR, 0666)
	w_hqew_list := bufio.NewWriter(f_hqew_list)
	w_hqew_list.WriteString("<table width='70%'>")
	for _, model := range modelSlice {

		base_url_sub := []string{"http://www.hqew.com/ic/", model, ".html"}
		base_url := strings.Join(base_url_sub, "")

		_, isStock, _ := getBody(base_url)
		fmt.Println(model, "\t", isStock)
		time.Sleep(1000 * time.Millisecond)

		if isStock {

			line_hqew_list := "<tr><td width='10%'>" + time.Now().Format("2006-01-02") + "</td><td width='20%'>" + model + "</td><td><a href='" + base_url + "' target='_blank'>" + base_url + "</a></td><tr>\r\n"
			w_hqew_list.WriteString(line_hqew_list)
		}
	}

	w_hqew_list.WriteString("</table><hr>\r\n")
	w_hqew_list.WriteString(string(old_data))
	w_hqew_list.Flush()
	defer f_hqew_list.Close()

}

// 根据网址 获取内容，确定是否需要进一步处理
func getBody(url string) (string, bool, error) {
	//根据url首先判断是否读取
	if strings.Contains(url, ".hqew.com/ic") {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("connect error")
			log.Fatal(err)
		}
		robots, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("connect error")
			log.Fatal(err)
		}
		res.Body.Close()
		body := string(robots)
		if len(body) == 0 {
			fmt.Println("null body")
		}

		//
		isMatch_stock, _ := regexp.MatchString(
			`xh ic_ico1`, body)
		isMatch_black, _ := regexp.MatchString(
			`target="_blank" class="g_fb"`, body)
		isMatch := isMatch_stock || isMatch_black
		if isMatch {
			return body, true, nil
		} else {
			return "NotStock", false, nil
		}
	} else {
		return "", false, nil
	}
}
