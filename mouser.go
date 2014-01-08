package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//爬虫mouser网站，根据大分类，爬出所有型号，将型号和厂家存入txt文档
// .\mouser.exe
func main() {
	var base_url string = "http://www.mouser.com/Optoelectronics/_/N-5g5v/"
	var fileName string = "mouser/Optoelectronics.txt"
	var totalPage int64 = 4166
	//默认为0
	var i int64 = 571

	f, _ := os.OpenFile(fileName, os.O_APPEND, 0666)
	defer f.Close()
	w := bufio.NewWriter(f)

	for i <= totalPage {
		url_sub := strconv.FormatInt(i*25, 10)
		url := base_url + "?No=" + url_sub

		//var url string = "http://www.mouser.com/Optoelectronics/_/N-5g5v/?No=25"
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("connect error 1")
			//log.Fatal(err)
			time.Sleep(30000 * time.Millisecond)
			continue
		}
		bodyByte, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("connect error 2")
			time.Sleep(30000 * time.Millisecond)
			continue
			//log.Fatal(err)
		}
		res.Body.Close()
		body := string(bodyByte)
		if len(body) != 0 {
			//fmt.Println(body)
		}
		//fmt.Println(body)

		oneLineSlice := strings.Split(body, "SearchResultsRow")

		for _, oneLine := range oneLineSlice {
			//fmt.Println(onePart)
			onePartSlice := strings.Split(oneLine, "</a><br />")
			if len(onePartSlice) > 1 {
				//fmt.Println(onePartSlice[1])
				//fmt.Println("\n\n")

				model := getMiddleString(onePartSlice[1], `3d">`, "")
				//fmt.Println(model)

				brand := getMiddleString(onePartSlice[2], `">`, "</a>")
				fmt.Println(model, "\t", brand)
				w.WriteString(model + "\t" + brand + "\r\n")
			}

			//fmt.Println("\n")
		}
		fmt.Println(url)
		i++
		time.Sleep(5000 * time.Millisecond)
	}
	w.Flush()
	//fmt.Println(stocksSlice)
}

//给定一个字符串str0，提取两个字符串str1,str2之间的字符串
func getMiddleString(str0, str1, str2 string) string {
	myExp, _ := regexp.Compile(str1 + "(.*)" + str2)
	mySlice := myExp.FindStringSubmatch(str0)
	if len(mySlice) == 0 {
		return ""
	} else {
		return mySlice[1]
	}
}
