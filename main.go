package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/tealeg/xlsx"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"github.com/badoux/goscraper"
)

var(
	output = "output.xlsx"
	rowsize int
	resultData []TableInfoStruct
	thread int
	fileName string
)
type TableInfoStruct struct {
	DomainName string
	StatusCode string
	WebTitle string
	ServerString string
	IpAddr string
	//其他类型
}

func init()  {
	flag.IntVar(&thread,"t",20,"探测线程")
	flag.StringVar(&fileName,"f","domain.txt","域名列表文件名")
}


func main() {
	flag.Parse()
	defer ants.Release()
	domainArr := readDomainListFromTxt(fileName)
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(thread, func(i interface{}) {
		getHttpInfo(i)
		wg.Done()
	})
	defer p.Release()
	for i := 0; i < len(domainArr); i++ {
		wg.Add(1)
		_ = p.Invoke(domainArr[i])
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	createExcel(resultData)
}

func createExcel(exceldata []TableInfoStruct){
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("result")
	if err != nil {
		fmt.Printf("createExcel函数出错："+err.Error())
	}
	{
		//表头
		row := sheet.AddRow()
		domainCell := row.AddCell()
		domainCell.Value = "域名"
		statuscodeCell := row.AddCell()
		statuscodeCell.Value = "状态码"
		webtitleCell := row.AddCell()
		webtitleCell.Value = "标题"
		serverCell := row.AddCell()
		serverCell.Value = "服务器信息"
		ipaddrCell := row.AddCell()
		ipaddrCell.Value = "IP"
	}
	for i:=0;i<len(exceldata);i++{
		row := sheet.AddRow()
		domainCell := row.AddCell()
		domainCell.Value = exceldata[i].DomainName
		statuscodeCell := row.AddCell()
		statuscodeCell.Value = exceldata[i].StatusCode
		webtitleCell := row.AddCell()
		webtitleCell.Value = exceldata[i].WebTitle
		serverCell := row.AddCell()
		serverCell.Value = exceldata[i].ServerString
		ipaddrCell := row.AddCell()
		ipaddrCell.Value = exceldata[i].IpAddr
	}
	err = file.Save(output)
	if err != nil {
		fmt.Printf("createExcel函数出错："+err.Error())
	}

}

func readDomainListFromTxt(f string) []string {
	var domains []string
	r, _ := os.Open(f)
	defer r.Close()
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		domains = append(domains, line)
	}
	domains = RemoveDuplicatesAndEmpty(domains)
	return domains
}

func getHttpInfo(domain interface{})  {
	url := "http://"+domain.(string)
	tab := TableInfoStruct{}
	tab.DomainName = domain.(string)	//统一原始域名，部分301跳转后与原始不同
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	tab.ServerString = resp.Header.Get("Server")
	tab.StatusCode = strconv.Itoa(resp.StatusCode)
	tab.WebTitle = getwebtitle(url)
	ip, _ := net.ResolveIPAddr("ip",tab.DomainName)
	tab.IpAddr = ip.String()
	fmt.Println(tab.DomainName)
	resultData = append(resultData, tab)

}

//去重函数，对domains.txt操作
func RemoveDuplicatesAndEmpty(a []string) (ret []string){
	a_len := len(a)
	for i:=0; i < a_len; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue;
		}
		ret = append(ret, a[i])
	}
	return
}

func getwebtitle(weburl string) string {
	s, err := goscraper.Scrape(weburl, 5)
	if err != nil {
		return ""
	}
	return s.Preview.Title

}
