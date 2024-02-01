package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain,hasMX,hasSPF,hasDMARC,spfRecord,dmarcRecord")
	// 一直读取下一行数据
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}
}

func checkDomain(domainStr string) {
	var domain, hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	// 查找指定域名的 DNS 用来记录可以交换电子邮件的服务器
	mxRecords, err := net.LookupMX(domainStr)
	if err != nil {
		log.Println(err)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}
	// TXT 记录存储有关 SPF 的信息，该信息可以识别授权服务器以代表您的组织发送电子邮件
	txtRecord, err := net.LookupTXT(domainStr)
	if err != nil {
		log.Println(err)
	}
	for _, v := range txtRecord {
		// 检测字符串是否以什么开头
		if strings.HasPrefix(v, "v=spf1") {
			hasSPF = true
			spfRecord = v
			break
		}
	}
	dmarcRecords, err := net.LookupTXT("_dmarc." + domainStr)
	if err != nil {
		log.Println(err)
	}
	for _, v := range dmarcRecords {
		if strings.HasPrefix(v, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = v
			break
		}
	}

	fmt.Println(domain, hasMX, hasSPF, hasDMARC, spfRecord, dmarcRecord)
}
