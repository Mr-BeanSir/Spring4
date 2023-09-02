package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"regexp"
)

func main() {
	c := colly.NewCollector()
	//c.SetProxy("http://127.0.0.1:10087")
	c.Limit(&colly.LimitRule{
		Parallelism: 15,
	})
	log.Println("here")
	page := c.Clone()
	page.Async = true
	re, _ := regexp.Compile("tid=(\\d+)&")
	table, _ := os.OpenFile("F:\\Spring4.txt", os.O_WRONLY|os.O_APPEND, 0644)
	defer table.Close()
	page.OnHTML("#postlist", func(e *colly.HTMLElement) {
		tid := re.FindStringSubmatch(e.Request.URL.String())
		title := e.ChildText("#thread_subject")
		tag := e.ChildText("table:nth-child(1) > tbody > tr > td.plc.ptm.pbn.vwthd > h1 > a")
		content := e.ChildText("tr:nth-child(1) > td.plc > div.pct > div > div.t_fsz > table > tbody > tr > td")
		if len(tid) < 1 {
			log.Println(e.Request.URL.String())
			return
		}
		log.Println("visit page title:" + title + "tid:" + tid[1])
		filename := fmt.Sprintf("F:\\Sprint4Text\\%s.txt", tid[1])
		_, err := os.Stat(filename)
		var file *os.File
		if err != nil {
			file, _ = os.Create(filename)
		} else {
			file, _ = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
		}
		file.WriteString(content)
		file.Close()
		log.Println(title, tag, tid[1])
		table.WriteString(fmt.Sprintf("%s,%s,%s\n", tid[1], title, tag))
	})

	c.OnHTML("tr > th > div.tl_tit.cl > a.s.xst", func(e *colly.HTMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))
		page.Visit(url)
	})

	c.OnHTML("#fd_page_bottom > div > a.nxt", func(e *colly.HTMLElement) {
		log.Println("visit url:" + e.Request.AbsoluteURL(e.Attr("href")))
		c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.Visit("https://chunman4.com/forum.php?mod=forumdisplay&fid=47&orderby=dateline&orderby=dateline&filter=author&page=940")

	log.Println("执行完成？")
	fmt.Scanln()
}
