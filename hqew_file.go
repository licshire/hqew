package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//爬虫hqew二级网站, 将型号，厂家，DC，封装，描述（是否现货）存入txt文件
//从file中读取二级域名
// .\hqew_file.exe
func main() {
	dat, _ := ioutil.ReadFile("hqew_list_file.txt")

	secondDomainSlice := strings.Split(string(dat), "\r\n")
	fmt.Print(secondDomainSlice)
	for _, secondDomain := range secondDomainSlice {

		base_url_sub := []string{"http://", secondDomain, ".hqew.com/ic/ic.html"}
		base_url := strings.Join(base_url_sub, "")

		body, isMatchStock, _ := getBody(strings.Join([]string{base_url, "?Page=1"}, ""))
		if isMatchStock {
			page := getPage(body)
			companyName, companyShortName := getCompanyName(body)
			fileName := strings.Join([]string{"result/", secondDomain, "_", companyShortName, "_",
				time.Now().Format("2006-01-02"), ".txt"}, "")
			os.Remove(fileName)
			fmt.Println(fileName)
			f, _ := os.Create(fileName)
			w := bufio.NewWriter(f)
			for i := 1; i <= page; i++ {
				if i == -1 {
					break
				}
				if i == -1 {
					continue
				}
				pageString := strconv.FormatInt(int64(i), 10)
				url_sub := []string{base_url, "?Page=", pageString}
				url := strings.Join(url_sub, "")

				body, _, _ := getBody(url)
				writeFile(body, w)

				fmt.Println(url)
				time.Sleep(5000 * time.Millisecond)
			}
			//将公司名和公司链接写入hqew_list.txt
			f_hqew_list, _ := os.OpenFile("hqew_list.txt", os.O_APPEND, 0666)
			w_hqew_list := bufio.NewWriter(f_hqew_list)
			line_hqew_list := strings.Join([]string{
				time.Now().Format("2006-01-02"),
				companyName, base_url, "\r\n"}, "\t")
			fmt.Println(line_hqew_list)
			w_hqew_list.WriteString(line_hqew_list)

			w.Flush()
			w_hqew_list.Flush()
			defer f_hqew_list.Close()
			defer f.Close()
		}
	}

}

// 根据网址 获取内容，确定是否需要进一步处理
func getBody(url string) (string, bool, error) {
	//根据url首先判断是否读取
	if strings.Contains(url, ".hqew.com/ic/ic.html") {
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
			return body, true, nil
		} else {
			return "NotList", false, nil
		}
	} else {
		return "", false, nil
	}
}

//根据body内容 获取页码数
func getPage(body string) int {
	stocksExp := regexp.MustCompile(
		`找到<i id="ctl00_cph_Content_pager2_spantotal">(.*)</i>条结果`)
	stocksSlice := stocksExp.FindStringSubmatch(body)
	//fmt.Println(stocksSlice[1])
	stocks, _ := strconv.ParseInt(stocksSlice[1], 10, 0)
	//fmt.Println(stocksSlice)
	if stocks%30 == 0 {
		return int(stocks / 30)
	} else {
		return int(stocks/30) + 1
	}
}

//根据body内容 获取公司名
func getCompanyName(body string) (string, string) {
	companyNameExp := regexp.MustCompile(
		`<li class="headtitle"><h1 title="(.*)">(.*)</h1></li>`)
	companyNameSlice := companyNameExp.FindStringSubmatch(body)
	companyName := companyNameSlice[2]
	companyName_short := strings.Replace(companyName, "有限公司", "", -1)
	companyName_short = strings.Replace(companyName_short, "深圳市", "", -1)
	companyName_short = strings.Replace(companyName_short, "经营部", "", -1)
	//companyName_short = strings.Replace(companyName_short, "香港", "", -1)
	//fmt.Println(companyName, "\t", companyName_short)
	//fmt.Println(stocks)
	return companyName, companyName_short
}

//根据网址 获取 华强网注册的唯一二级域名
func getSecondDomain(url string) string {
	if strings.Contains(url, "http://") {
		s := regexp.MustCompile("//").Split(url, 4)
		s2 := regexp.MustCompile(".hqew.com").Split(s[1], 5)
		return s2[0]
	} else {
		s2 := regexp.MustCompile(".hqew.com").Split(url, 5)
		return s2[0]
	}
}

//根据网页源文件 获取 型号表的内容。以<tr>开头,以</tr>结尾
func writeFile(str string, w *bufio.Writer) {
	_, companyName := getCompanyName(str)
	//fmt.Println(companyName)
	myExp, _ := regexp.Compile(`<tr class="tr0"><td class="c1".*</td></tr>`)
	trsSlice := myExp.FindAllString(str, -1)

	if len(trsSlice) != 0 {

		tdsSlice := strings.Split(trsSlice[0], `</tr>`)
		for _, tdString := range tdsSlice {
			var modelExp *regexp.Regexp
			var modelSlice []string

			var stockStatus bool = false
			//获取型号 不带现货标志
			modelExp = regexp.MustCompile(`<td class="c1"><span class="g_fl">(.*)</span></td>`)
			modelSlice = modelExp.FindStringSubmatch(tdString)
			if len(modelSlice) == 0 {
				//获取型号 带现货标志
				modelExp = regexp.MustCompile(`<td class="c1"><span class="g_fl">(.*)</span><a href='http://www.hqew.com/about/service_03.html' target='_blank' class='g_fl g_xh'></a></td>`)
				modelSlice = modelExp.FindStringSubmatch(tdString)
				if len(modelSlice) != 0 {
					stockStatus = true
				}
			}
			//获取厂家 brand
			brandExp, _ := regexp.Compile(`<td class="c2">(.*)</td><td class="c3">`)
			brandSlice := brandExp.FindStringSubmatch(tdString)
			//获取 DC
			dcExp, _ := regexp.Compile(`<td class="c3">(.*)</td><td class="c4">`)
			dcSlice := dcExp.FindStringSubmatch(tdString)
			//获取 数量
			amountExp, _ := regexp.Compile(`<td class="c4">(.*)</td><td class="c5">`)
			amountSlice := amountExp.FindStringSubmatch(tdString)
			//获取 封装
			packExp, _ := regexp.Compile(`<td class="c5">(.*)</td><td class="c6">`)
			packSlice := packExp.FindStringSubmatch(tdString)
			//获取 交易说明
			//tradeExp, _ := regexp.Compile(`<td class="c6">(.*)</td><td class="c7">`)
			//tradeSlice := tradeExp.FindStringSubmatch(tdString)
			//获取 仓库
			//warehouseExp, _ := regexp.Compile(`<td class="c7">(.*)</td><td class="c8">`)
			//warehouseSlice := warehouseExp.FindStringSubmatch(tdString)

			if len(modelSlice) != 0 {
				modelString := modelSlice[1]
				brandString := brandSlice[1]
				brandString = changeBrand(brandString)
				dcString := dcSlice[1]
				amountString := amountSlice[1]
				packString := packSlice[1]
				if packString == "&nbsp;" {
					packString = ""
				}
				//tradeString := tradeSlice[1]
				//warehouseString := warehouseSlice[1]
				var desc string
				if stockStatus == false {
					desc = strings.Join([]string{"from hqew, ", companyName}, "")
				} else {
					desc = strings.Join([]string{"from hqew, 现货 ", companyName}, "")
				}

				line := strings.Join([]string{
					modelString, "\t", desc, "\t",
					brandString, "\t", dcString, "\t",
					amountString, "\t", packString, "\r\n"}, "")
				fmt.Println(line)
				w.WriteString(line)
			}

		}

	}

	//return "test"
}

func changeBrand(brand string) string {
	switch brand {
	case "":
		brand = "Other"
	case "&nbsp;":
		brand = "Other"

	case "AIC/沛亨":
		brand = "AIC"
	case "ANALOGICTECH":
		brand = "ANALOGIC"

	case "BELLING/上海":
		brand = "BELLING"
	case "cypress":
		brand = "CYPRESS"
	case "CYPRESS专业":
		brand = "CYPRESS"

	case "DIODES":
		brand = "DIODES INC"
	case "FSC":
		brand = "FAIRCHILD"
	case "FSC原装":
		brand = "FAIRCHILD"
	case "FUJ":
		brand = "FUJITSU"
	case "FUJI":
		brand = "FUJITSU"
	case "HIT":
		brand = "HITACHI"
	case "INF英飞凌":
		brand = "INFINEON"
	case "LN南麟":
		brand = "NATLINEAR"

	case "Microne南京":
		brand = "MICRONE"
	case "MOT":
		brand = "MOTOROLA"
	case "NXP/PHI":
		brand = "NXP"
	case "NS国半":
		brand = "NXP"

	case "OB台湾昂宝":
		brand = "OB"

	case "ON":
		brand = "ON SEMI"
	case "ONSEMI":
		brand = "ON SEMI"
	case "PAN松下":
		brand = "PANASONIC"
	case "PAN":
		brand = "PANASONIC"
	case "PHI":
		brand = "PHILIPS"
	case "PHI飞利蒲":
		brand = "PHILIPS"
	case "RICHTEK/立锜":
		brand = "RICHTEK"
	case "SAM":
		brand = "SAMSUNG"
	case "SEIKO/精工":
		brand = "SEIKO"
	case "SHAPR/夏普":
		brand = "SHAPR"
	case "TI德州":
		brand = "SHAPR"

	case "TOS东芝":
		brand = "TOSHIBA"
	case "TOS":
		brand = "TOSHIBA"
	case "佰鸿":
		brand = "BRIGHT"
	case "长电":
		brand = "CJ"
	case "东芝":
		brand = "TOSHIBA"
	case "三星":
		brand = "SAMSUNG"
	case "三洋":
		brand = "SANYO"
	case "松下":
		brand = "PANASONIC"
	case "台湾远翔":
		brand = "FEELING-TECH"

	default:
		brand = brand
	}
	return brand

}
