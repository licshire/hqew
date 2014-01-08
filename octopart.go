package main

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	//"time"
)

//搜索octopart网站，根据部分型号和厂家。补充库存型号列表。
// .\octopart.exe everlight 19-21*
func main() {
	var search string
	for i, search_sub := range os.Args {
		if i == 0 {
			continue
		}
		if i == 1 {
			search = search_sub
		} else {
			search += "+" + search_sub
		}
	}
	//model = strings.Replace(model, "+", ",", -1)
	url := "http://octopart.com/api/v3/parts/search"
	//url := "http://octopart.com/api/v3/parts/match"
	url += "?apikey=004834a2"
	url += "&q=" + search
	//url += "&queries=[{\"mpn\":\"PY1111C\"}]"
	//url += "&pretty_print=true"
	url += "&start=1&limit=999"

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
	if len(body) == 0 {
		fmt.Println("null body")
	}
	//fmt.Println(body)

	onelineSlice := strings.Split(body, "http://octopart.com")
	//fmt.Println(onelineSlice[0])
	for _, oneline := range onelineSlice {
		returnModel := getMiddleString(oneline, "mpn\": \"", "\", \"brand")
		if returnModel != "" {
			fmt.Println(returnModel)
		}

	}
	numberofResult := getMiddleString(body, "\"hits\": ", ", \"stats_results\"")
	fmt.Println("\nResults Number: ", numberofResult)
	fmt.Println(url)
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
