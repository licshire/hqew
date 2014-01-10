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
	var fileName string = "Optoelectronics.txt"
	var fileNameLog string = "log.txt"
	var totalPage int64 = 4166
	//默认为0
	var i int64 = 1576

	args := os.Args
	if len(args) == 2 {
		noString := args[1]
		noInt64, _ := strconv.ParseInt(noString, 10, 0)
		i = noInt64 / 25
	} else {
		fmt.Println("please input .\\mouser.exe noNumber ")
		os.Exit(1)
		//fmt.Println(base_url)
	}
	fmt.Println(i)
	//os.Exit(1)

	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	w := bufio.NewWriter(f)
	fLog, _ := os.OpenFile(fileNameLog, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	wLog := bufio.NewWriter(fLog)

	for i <= totalPage {
		url_sub := strconv.FormatInt(i*25, 10)
		url := base_url + "?No=" + url_sub

		//var url string = "http://www.mouser.com/Optoelectronics/_/N-5g5v/?No=25"
		fmt.Println("before http.get(url)")
		res, err := http.Get(url)
		fmt.Println("after http.get(url)")
		fmt.Println(res.StatusCode)

		if res.StatusCode != 200 {
			fmt.Println("connect error 0")
			wLog.WriteString(res.Status + "\t" + "connect error 0\t" + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
			wLog.Flush()
			//log.Fatal(err)
			time.Sleep(30000 * time.Millisecond)
			continue
		}
		//os.Exit(1)
		if err != nil {
			fmt.Println("connect error 1")
			wLog.WriteString("connect error 1\t" + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
			wLog.Flush()
			//log.Fatal(err)
			time.Sleep(30000 * time.Millisecond)
			continue
		}
		fmt.Println("before ioutil.ReadAll(res.Body)")
		bodyByte, err := ioutil.ReadAll(res.Body)
		fmt.Println("after ioutil.ReadAll(res.Body)")
		if err != nil {
			for {
				fmt.Println("connect error 2")
				wLog.WriteString("connect error 2\t" + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
				wLog.Flush()
				time.Sleep(30000 * time.Millisecond)
				bodyByte, err = ioutil.ReadAll(res.Body)
				if err == nil {
					break
				}
				//log.Fatal(err)
			}

		}
		fmt.Println("res.Body.Close")
		res.Body.Close()
		body := string(bodyByte)
		if len(body) != 0 {
			//fmt.Println(body)
		}
		//fmt.Println(body)

		oneLineSlice := strings.Split(body, "SearchResultsRow")
		var number int64 = 0
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
				_, err := w.WriteString(model + "\t" + brand + "\r\n")
				if err != nil {
					fmt.Println("write to file error 1")
					wLog.WriteString("write to file error 1\t" + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
					wLog.Flush()
					continue
				}
				number++
			}
			//fmt.Println("\n")
		}
		fmt.Println(url, "\t", time.Now().Format("2006-01-02 15:04:05"))
		wLog.WriteString(url + "\t" + time.Now().Format("2006-01-02 15:04:05") + "\t" + "获取型号" + strconv.FormatInt(number, 10) + "\r\n")
		i++

		err2 := w.Flush()
		if err2 != nil {
			fmt.Println("write to file error 2")
			wLog.WriteString("write to file error 2\t" + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
			continue
		}
		wLog.Flush()
		time.Sleep(5000 * time.Millisecond)
	}
	//fmt.Println(stocksSlice)
	defer f.Close()
	defer fLog.Close()
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
	return ""
}
