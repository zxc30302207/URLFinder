package result

import (
	"bufio"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/pingc0y/URLFinder/cmd"
	"github.com/pingc0y/URLFinder/mode"
	"github.com/pingc0y/URLFinder/util"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:embed report.html
var html string

var (
	ResultJs  []mode.Link
	ResultUrl []mode.Link
	Fuzzs     []mode.Link
	Infos     []mode.Info

	EndUrl   []string
	Jsinurl  map[string]string
	Jstourl  map[string]string
	Urltourl map[string]string
	Domains  []string
	Redirect map[string]bool
)

func outHtmlString(link mode.Link) string {
	ht := `<tr class="ant-table-row ant-table-row-level-0" data-row-key="0">
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					<a href="` + link.Url + `" target="_blank" >
						` + link.Url + ` </a>
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					` + link.Status + `
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					` + link.Size + `
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					` + link.Title + `
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					<a href="` + link.Redirect + `" target="_blank" style="display:inline-bconfig.Lock">
						` + link.Redirect + ` </a>
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					<a href="` + link.Source + `" target="_blank" style="display:inline-bconfig.Lock">
						` + link.Source + ` </a>
				</td>
			</tr>`
	return ht
}

func outHtmlInfoString(ty, val, sou string) string {
	ht := `<tr class="ant-table-row ant-table-row-level-0" data-row-key="0">
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					` + ty + `
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					` + val + `
				</td>
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					<a href="` + sou + `" target="_blank" >
						` + sou + ` </a>
				</td>
			</tr>`
	return ht
}

func outHtmlDomainString(domain string) string {
	ht := `<tr class="ant-table-row ant-table-row-level-0" data-row-key="0">
				<td class="ant-table-column-has-actions ant-table-column-has-sorters">
					` + domain + `
				</td>
			</tr>`
	return ht
}

// 匯出csv
func OutFileCsv(out string) {
	//獲取域名
	var host string
	re := regexp.MustCompile("([a-z0-9\\-]+\\.)*([a-z0-9\\-]+\\.[a-z0-9\\-]+)(:[0-9]+)?")
	hosts := re.FindAllString(cmd.U, 1)
	if len(hosts) == 0 {
		host = cmd.U
	} else {
		host = hosts[0]
	}

	//抓取的域名優先排序
	if cmd.S != "" {
		ResultUrl = util.SelectSort(ResultUrl)
		ResultJs = util.SelectSort(ResultJs)
	}
	ResultJsHost, ResultJsOther := util.UrlDispose(ResultJs, host, util.GetHost(cmd.U))
	ResultUrlHost, ResultUrlOther := util.UrlDispose(ResultUrl, host, util.GetHost(cmd.U))
	Domains = util.GetDomains(util.MergeArray(ResultJs, ResultUrl))
	var fileName string
	if out != "" {
		fileName = out
	} else {
		//輸出到檔案
		if strings.Contains(host, ":") {
			host = strings.Replace(host, ":", "：", -1)
		}
		//在當前資料夾建立資料夾
		err := os.MkdirAll(cmd.O+"/"+host, 0755)
		if err != nil {
			fmt.Printf(cmd.O+"/"+host+" 目錄建立失敗 ：%s", err)
			return
		}
		//多相同url處理
		fileName = cmd.O + "/" + host + "/" + host + ".csv"
		for fileNum := 1; util.Exists(fileName); fileNum++ {
			fileName = cmd.O + "/" + host + "/" + host + "(" + strconv.Itoa(fileNum) + ").csv"
		}
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)

	resultWriter := csv.NewWriter(file)
	// 寫資料到檔案
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	if cmd.S == "" {
		resultWriter.Write([]string{"url", "Source"})
	} else {
		resultWriter.Write([]string{"url", "Status", "Size", "Title", "Redirect", "Source"})
	}
	if cmd.D == "" {
		resultWriter.Write([]string{strconv.Itoa(len(ResultJsHost)) + " JS to " + util.GetHost(cmd.U)})
	} else {
		resultWriter.Write([]string{strconv.Itoa(len(ResultJsHost)+len(ResultJsOther)) + " JS to " + cmd.D})
	}

	for _, j := range ResultJsHost {
		if cmd.S != "" {
			resultWriter.Write([]string{j.Url, j.Status, j.Size, "", j.Redirect, j.Source})
		} else {
			resultWriter.Write([]string{j.Url, j.Source})
		}
	}

	if cmd.D == "" {
		resultWriter.Write([]string{""})
		resultWriter.Write([]string{strconv.Itoa(len(ResultJsOther)) + " JS to Other"})
	}
	for _, j := range ResultJsOther {
		if cmd.S != "" {
			resultWriter.Write([]string{j.Url, j.Status, j.Size, "", j.Redirect, j.Source})
		} else {
			resultWriter.Write([]string{j.Url, j.Source})
		}
	}

	resultWriter.Write([]string{""})
	if cmd.D == "" {
		resultWriter.Write([]string{strconv.Itoa(len(ResultUrlHost)) + " URL to " + util.GetHost(cmd.U)})
	} else {
		resultWriter.Write([]string{strconv.Itoa(len(ResultUrlHost)+len(ResultUrlOther)) + " URL to " + cmd.D})
	}

	for _, u := range ResultUrlHost {
		if cmd.S != "" {
			resultWriter.Write([]string{u.Url, u.Status, u.Size, u.Title, u.Redirect, u.Source})
		} else {
			resultWriter.Write([]string{u.Url, u.Source})
		}
	}
	if cmd.D == "" {
		resultWriter.Write([]string{""})
		resultWriter.Write([]string{strconv.Itoa(len(ResultUrlOther)) + " URL to Other"})
	}
	for _, u := range ResultUrlOther {
		if cmd.S != "" {
			resultWriter.Write([]string{u.Url, u.Status, u.Size, u.Title, u.Redirect, u.Source})
		} else {
			resultWriter.Write([]string{u.Url, u.Source})
		}
	}
	if cmd.S != "" && cmd.Z != 0 {
		resultWriter.Write([]string{""})
		resultWriter.Write([]string{strconv.Itoa(len(Fuzzs)) + " URL to Fuzz"})
		Fuzzs = util.SelectSort(Fuzzs)
		for _, u := range Fuzzs {
			resultWriter.Write([]string{u.Url, u.Status, u.Size, u.Title, u.Redirect, "Fuzz"})
		}
	}
	resultWriter.Write([]string{""})
	resultWriter.Write([]string{strconv.Itoa(len(Domains)) + " Domain"})
	for _, u := range Domains {
		resultWriter.Write([]string{u})
	}
	resultWriter.Write([]string{""})
	resultWriter.Write([]string{"Phone"})
	for i := range Infos {
		for i2 := range Infos[i].Phone {
			resultWriter.Write([]string{Infos[i].Phone[i2], "", "", "", Infos[i].Source})
		}
	}
	resultWriter.Write([]string{""})
	resultWriter.Write([]string{"Email"})
	for i := range Infos {
		for i2 := range Infos[i].Email {
			resultWriter.Write([]string{Infos[i].Email[i2], "", "", "", Infos[i].Source})
		}
	}
	resultWriter.Write([]string{""})
	resultWriter.Write([]string{"Email"})
	for i := range Infos {
		for i2 := range Infos[i].IDcard {
			resultWriter.Write([]string{Infos[i].IDcard[i2], "", "", "", Infos[i].Source})
		}
	}
	resultWriter.Write([]string{""})
	resultWriter.Write([]string{"JWT"})
	for i := range Infos {
		for i2 := range Infos[i].JWT {
			resultWriter.Write([]string{Infos[i].JWT[i2], "", "", "", Infos[i].Source})
		}
	}
	resultWriter.Write([]string{""})
	resultWriter.Write([]string{"Other"})
	for i := range Infos {
		for i2 := range Infos[i].Other {
			resultWriter.Write([]string{Infos[i].Other[i2], "", "", "", Infos[i].Source})
		}
	}

	resultWriter.Flush()

	fmt.Println(strconv.Itoa(len(ResultJsHost)+len(ResultJsOther))+"JS + "+strconv.Itoa(len(ResultUrlHost)+len(ResultUrlOther))+"URL --> ", file.Name())

	return
}

// 匯出json
func OutFileJson(out string) {
	jsons := make(map[string]interface{})
	var info map[string][]map[string]string
	//獲取域名
	var host string
	re := regexp.MustCompile("([a-z0-9\\-]+\\.)*([a-z0-9\\-]+\\.[a-z0-9\\-]+)(:[0-9]+)?")
	hosts := re.FindAllString(cmd.U, 1)
	if len(hosts) == 0 {
		host = cmd.U
	} else {
		host = hosts[0]
	}
	//抓取的域名優先排序
	if cmd.S != "" {
		ResultUrl = util.SelectSort(ResultUrl)
		ResultJs = util.SelectSort(ResultJs)
	}
	ResultJsHost, ResultJsOther := util.UrlDispose(ResultJs, host, util.GetHost(cmd.U))
	ResultUrlHost, ResultUrlOther := util.UrlDispose(ResultUrl, host, util.GetHost(cmd.U))
	Domains = util.GetDomains(util.MergeArray(ResultJs, ResultUrl))

	if len(Infos) > 0 {
		info = make(map[string][]map[string]string)
		info["IDcard"] = nil
		info["JWT"] = nil
		info["Email"] = nil
		info["Phone"] = nil
		info["Other"] = nil
	}

	for i := range Infos {
		for i2 := range Infos[i].IDcard {
			info["IDcard"] = append(info["IDcard"], map[string]string{"IDcard": Infos[i].IDcard[i2], "Source": Infos[i].Source})
		}
		for i2 := range Infos[i].JWT {
			info["JWT"] = append(info["JWT"], map[string]string{"JWT": Infos[i].JWT[i2], "Source": Infos[i].Source})
		}
		for i2 := range Infos[i].Email {
			info["Email"] = append(info["Email"], map[string]string{"Email": Infos[i].Email[i2], "Source": Infos[i].Source})
		}
		for i2 := range Infos[i].Phone {
			info["Phone"] = append(info["Phone"], map[string]string{"Phone": Infos[i].Phone[i2], "Source": Infos[i].Source})
		}
		for i2 := range Infos[i].Other {
			info["Other"] = append(info["Other"], map[string]string{"Other": Infos[i].Other[i2], "Source": Infos[i].Source})
		}
	}

	var fileName string
	if out != "" {
		fileName = out
	} else {
		//輸出到檔案
		if strings.Contains(host, ":") {
			host = strings.Replace(host, ":", "：", -1)
		}
		//在當前資料夾建立資料夾
		err := os.MkdirAll(cmd.O+"/"+host, 0755)
		if err != nil {
			fmt.Printf(cmd.O+"/"+host+" 目錄建立失敗 ：%s", err)
			return
		}
		//多相同url處理
		fileName = cmd.O + "/" + host + "/" + host + ".json"
		for fileNum := 1; util.Exists(fileName); fileNum++ {
			fileName = cmd.O + "/" + host + "/" + host + "(" + strconv.Itoa(fileNum) + ").json"
		}
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("建立失敗：%s", err)
		return
	}
	if cmd.D == "" {
		jsons["jsOther"] = ResultJsOther
		jsons["urlOther"] = ResultUrlOther
	}
	jsons["js"] = ResultJsHost
	jsons["url"] = ResultUrlHost
	jsons["info"] = info
	jsons["fuzz"] = Fuzzs
	jsons["domain"] = Domains
	if cmd.S != "" && cmd.Z != 0 {
		Fuzzs = util.SelectSort(Fuzzs)
		if len(Fuzzs) > 0 {
			jsons["fuzz"] = Fuzzs
		} else {
			jsons["fuzz"] = nil
		}

	}

	defer file.Close()

	data, err := json.Marshal(jsons)
	if err != nil {
		fmt.Printf("json化失敗：%s", err)
		return
	}
	buf := bufio.NewWriter(file)
	// 位元組寫入
	buf.Write(data)
	// 將緩衝中的資料寫入
	err = buf.Flush()
	if err != nil {
		fmt.Println("json儲存失敗:", err)
	}
	fmt.Println(strconv.Itoa(len(ResultJsHost)+len(ResultJsOther))+"JS + "+strconv.Itoa(len(ResultUrlHost)+len(ResultUrlOther))+"URL --> ", file.Name())
	return
}

// 匯出html
func OutFileHtml(out string) {
	htmlTemp := html
	//獲取域名
	var host string
	re := regexp.MustCompile("([a-z0-9\\-]+\\.)*([a-z0-9\\-]+\\.[a-z0-9\\-]+)(:[0-9]+)?")
	hosts := re.FindAllString(cmd.U, 1)
	if len(hosts) == 0 {
		host = cmd.U
	} else {
		host = hosts[0]
	}

	//抓取的域名優先排序
	if cmd.S != "" {
		ResultUrl = util.SelectSort(ResultUrl)
		ResultJs = util.SelectSort(ResultJs)
	}
	ResultJsHost, ResultJsOther := util.UrlDispose(ResultJs, host, util.GetHost(cmd.U))
	ResultUrlHost, ResultUrlOther := util.UrlDispose(ResultUrl, host, util.GetHost(cmd.U))
	Domains = util.GetDomains(util.MergeArray(ResultJs, ResultUrl))

	var fileName string
	if out != "" {
		fileName = out
	} else {
		//輸出到檔案
		if strings.Contains(host, ":") {
			host = strings.Replace(host, ":", "：", -1)
		}
		//在當前資料夾建立資料夾
		err := os.MkdirAll(cmd.O+"/"+host, 0755)
		if err != nil {
			fmt.Printf(cmd.O+"/"+host+" 目錄建立失敗 ：%s", err)
			return
		}
		//多相同url處理
		fileName = cmd.O + "/" + host + "/" + host + ".html"
		for fileNum := 1; util.Exists(fileName); fileNum++ {
			fileName = cmd.O + "/" + host + "/" + host + "(" + strconv.Itoa(fileNum) + ").html"
		}
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)

	file.WriteString("\xEF\xBB\xBF") // 寫入UTF-8 BOM,防止中文亂碼
	// 寫資料到檔案
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if cmd.D == "" {
		htmlTemp = strings.Replace(htmlTemp, "{urlHost}", util.GetHost(cmd.U), -1)
	} else {
		htmlTemp = strings.Replace(htmlTemp, "{urlHost}", cmd.D, -1)
	}
	var ResultJsHostStr string
	for _, j := range ResultJsHost {
		ResultJsHostStr += outHtmlString(j)
	}
	htmlTemp = strings.Replace(htmlTemp, "{JS}", ResultJsHostStr, -1)

	var ResultJsOtherStr string
	for _, j := range ResultJsOther {
		ResultJsOtherStr += outHtmlString(j)
	}
	htmlTemp = strings.Replace(htmlTemp, "{JSOther}", ResultJsOtherStr, -1)

	var ResultUrlHostStr string
	for _, u := range ResultUrlHost {
		ResultUrlHostStr += outHtmlString(u)
	}
	htmlTemp = strings.Replace(htmlTemp, "{URL}", ResultUrlHostStr, -1)

	var ResultUrlOtherStr string
	for _, u := range ResultUrlOther {
		ResultUrlOtherStr += outHtmlString(u)
	}
	htmlTemp = strings.Replace(htmlTemp, "{URLOther}", ResultUrlOtherStr, -1)

	var FuzzsStr string
	if cmd.S != "" && cmd.Z != 0 {
		Fuzzs = util.SelectSort(Fuzzs)
		for _, u := range Fuzzs {
			FuzzsStr += outHtmlString(u)
		}
	}
	htmlTemp = strings.Replace(htmlTemp, "{Fuzz}", FuzzsStr, -1)

	var DomainsStr string
	for _, u := range Domains {
		DomainsStr += outHtmlDomainString(u)
	}
	htmlTemp = strings.Replace(htmlTemp, "{Domains}", DomainsStr, -1)

	var Infostr string
	for i := range Infos {
		for i2 := range Infos[i].Phone {
			Infostr += outHtmlInfoString("Phone", Infos[i].Phone[i2], Infos[i].Source)
		}
	}
	for i := range Infos {
		for i2 := range Infos[i].Email {
			Infostr += outHtmlInfoString("Email", Infos[i].Email[i2], Infos[i].Source)
		}
	}
	for i := range Infos {
		for i2 := range Infos[i].IDcard {
			Infostr += outHtmlInfoString("IDcard", Infos[i].IDcard[i2], Infos[i].Source)
		}
	}
	for i := range Infos {
		for i2 := range Infos[i].JWT {
			Infostr += outHtmlInfoString("JWT", Infos[i].JWT[i2], Infos[i].Source)
		}
	}
	for i := range Infos {
		for i2 := range Infos[i].Other {
			Infostr += outHtmlInfoString("Other", Infos[i].Other[i2], Infos[i].Source)
		}
	}
	htmlTemp = strings.Replace(htmlTemp, "{Info}", Infostr, -1)
	writer.WriteString(htmlTemp)
	writer.Flush() //內容是先寫到快取對,所以需要呼叫flush將快取對資料真正寫到檔案中
	fmt.Println(strconv.Itoa(len(ResultJsHost)+len(ResultJsOther))+"JS + "+strconv.Itoa(len(ResultUrlHost)+len(ResultUrlOther))+"URL --> ", file.Name())
	return
}

// 列印
func Print() {
	//獲取域名
	var host string
	re := regexp.MustCompile("([a-z0-9\\-]+\\.)*([a-z0-9\\-]+\\.[a-z0-9\\-]+)(:[0-9]+)?")
	hosts := re.FindAllString(cmd.U, 1)
	if len(hosts) == 0 {
		host = cmd.U
	} else {
		host = hosts[0]
	}
	//列印JS
	if cmd.S != "" {
		ResultJs = util.SelectSort(ResultJs)
		ResultUrl = util.SelectSort(ResultUrl)

	}
	//抓取的域名優先排序
	ResultJsHost, ResultJsOther := util.UrlDispose(ResultJs, host, util.GetHost(cmd.U))
	ResultUrlHost, ResultUrlOther := util.UrlDispose(ResultUrl, host, util.GetHost(cmd.U))
	Domains = util.GetDomains(util.MergeArray(ResultJs, ResultUrl))
	var ulen string
	if len(ResultUrl) != 0 {
		uleni := 0
		for _, u := range ResultUrl {
			uleni += len(u.Url)
		}
		ulen = strconv.Itoa(uleni/len(ResultUrl) + 10)
	}
	var jlen string
	if len(ResultJs) != 0 {
		jleni := 0
		for _, j := range ResultJs {
			jleni += len(j.Url)
		}
		jlen = strconv.Itoa(jleni/len(ResultJs) + 10)
	}
	if cmd.D == "" {
		fmt.Println(strconv.Itoa(len(ResultJsHost)) + " JS to " + util.GetHost(cmd.U))
	} else {
		fmt.Println(strconv.Itoa(len(ResultJsHost)+len(ResultJsOther)) + " JS to " + cmd.D)
	}
	for _, j := range ResultJsHost {
		if cmd.S != "" {
			if strings.HasPrefix(j.Status, "2") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+jlen+"s", j.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", j.Status, j.Size, j.Source))
			} else if strings.HasPrefix(j.Status, "3") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+jlen+"s", j.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", j.Status, j.Size, j.Source))
			} else {
				fmt.Printf(color.LightBlue.Sprintf("%-"+jlen+"s", j.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", j.Status, j.Size, j.Source))
			}
		} else {
			fmt.Println(j.Url)
		}
	}
	if cmd.D == "" {
		fmt.Println("\n" + strconv.Itoa(len(ResultJsOther)) + " JS to Other")
	}
	for _, j := range ResultJsOther {
		if cmd.S != "" {
			if strings.HasPrefix(j.Status, "2") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+jlen+"s", j.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", j.Status, j.Size, j.Source))
			} else if strings.HasPrefix(j.Status, "3") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+jlen+"s", j.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", j.Status, j.Size, j.Source))
			} else {
				fmt.Printf(color.LightBlue.Sprintf("%-"+jlen+"s", j.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", j.Status, j.Size, j.Source))
			}
		} else {
			fmt.Println(j.Url)
		}
	}

	fmt.Println("\n  ")

	if cmd.D == "" {
		fmt.Println(strconv.Itoa(len(ResultUrlHost)) + " URL to " + util.GetHost(cmd.U))
	} else {
		fmt.Println(strconv.Itoa(len(ResultUrlHost)+len(ResultUrlOther)) + " URL to " + cmd.D)
	}

	for _, u := range ResultUrlHost {
		urlx, err := url.QueryUnescape(u.Url)
		if err == nil {
			u.Url = urlx
		}
		if cmd.S != "" && len(u.Title) != 0 {
			if u.Status == "疑似危險路由" {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Source: %s ]\n", u.Status, u.Source))
			} else if strings.HasPrefix(u.Status, "2") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
			} else if strings.HasPrefix(u.Status, "3") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
			} else {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
			}
		} else if cmd.S != "" {
			if strings.HasPrefix(u.Status, "2") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
			} else if strings.HasPrefix(u.Status, "3") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
			} else {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
			}
		} else {
			fmt.Println(u.Url)
		}
	}
	if cmd.D == "" {
		fmt.Println("\n" + strconv.Itoa(len(ResultUrlOther)) + " URL to Other")
	}
	for _, u := range ResultUrlOther {
		urlx, err := url.QueryUnescape(u.Url)
		if err == nil {
			u.Url = urlx
		}
		if cmd.S != "" && len(u.Title) != 0 {
			if u.Status == "疑似危險路由" {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Source: %s ]\n", u.Status, u.Source))
			} else if strings.HasPrefix(u.Status, "2") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
			} else if strings.HasPrefix(u.Status, "3") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
			} else {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
			}
		} else if cmd.S != "" {
			if strings.HasPrefix(u.Status, "2") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
			} else if strings.HasPrefix(u.Status, "3") {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
			} else {
				fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
			}
		} else {
			fmt.Println(u.Url)
		}
	}

	if cmd.S != "" && cmd.Z != 0 {
		fmt.Println("\n" + strconv.Itoa(len(Fuzzs)) + " URL to Fuzz")
		Fuzzs = util.SelectSort(Fuzzs)
		for _, u := range Fuzzs {
			if len(u.Title) != 0 {
				if u.Status == "疑似危險路由" {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Source: %s ]\n", u.Status, u.Source))
				} else if strings.HasPrefix(u.Status, "2") {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
				} else if strings.HasPrefix(u.Status, "3") {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
				} else {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Title: %s, Source: %s ]\n", u.Status, u.Size, u.Title, u.Source))
				}
			} else {
				if strings.HasPrefix(u.Status, "2") {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightGreen.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
				} else if strings.HasPrefix(u.Status, "3") {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightYellow.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
				} else {
					fmt.Printf(color.LightBlue.Sprintf("%-"+ulen+"s", u.Url) + color.LightRed.Sprintf(" [ Status: %s, Size: %s, Source: %s ]\n", u.Status, u.Size, u.Source))
				}
			}
		}
	}
	fmt.Println("\n" + strconv.Itoa(len(Domains)) + " Domain")
	for _, u := range Domains {
		fmt.Printf(color.LightBlue.Sprintf("%s \n", u))

	}

	if len(Infos) > 0 {
		fmt.Println("\n Phone ")
		for i := range Infos {
			for i2 := range Infos[i].Phone {
				fmt.Printf(color.LightBlue.Sprintf("%-10s", Infos[i].Phone[i2]) + color.LightGreen.Sprintf(" [ Source: %s ]\n", Infos[i].Source))
			}
		}
		fmt.Println("\n Email ")
		for i := range Infos {
			for i2 := range Infos[i].Email {
				fmt.Printf(color.LightBlue.Sprintf("%-10s", Infos[i].Email[i2]) + color.LightGreen.Sprintf(" [ Source: %s ]\n", Infos[i].Source))
			}
		}
		fmt.Println("\n IDcard ")
		for i := range Infos {
			for i2 := range Infos[i].IDcard {
				fmt.Printf(color.LightBlue.Sprintf("%-10s", Infos[i].IDcard[i2]) + color.LightGreen.Sprintf(" [ Source: %s ]\n", Infos[i].Source))
			}
		}
		fmt.Println("\n JWT ")
		for i := range Infos {
			for i2 := range Infos[i].JWT {
				fmt.Printf(color.LightBlue.Sprintf("%-10s", Infos[i].JWT[i2]) + color.LightGreen.Sprintf(" [ Source: %s ]\n", Infos[i].Source))
			}
		}

		fmt.Println("\n Other ")
		for i := range Infos {
			for i2 := range Infos[i].Other {
				fmt.Printf(color.LightBlue.Sprintf("%-10s", Infos[i].Other[i2]) + color.LightGreen.Sprintf(" [ Source: %s ]\n", Infos[i].Source))
			}
		}

	}

}
