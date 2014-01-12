package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("search:", r.Form["search"])
		fmt.Println("password:", r.Form["password"])

		var search string
		for i, search_sub := range r.Form["search"] {

			if i == 0 {
				search = search_sub
			} else {
				search += "+" + search_sub
			}
		}
		search = strings.Replace(search, " ", "+", -1)
		url := "http://octopart.com/api/v3/parts/search"
		url += "?apikey=004834a2"
		url += "&q=" + search

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
				fmt.Fprintln(w, returnModel)
			}

		}
		numberofResult := getMiddleString(body, "\"hits\": ", ", \"stats_results\"")
		fmt.Fprintln(w, "\nResults Number: ", numberofResult)
		fmt.Fprintln(w, url)
	}
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

func main() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	http.HandleFunc("/search", search)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
