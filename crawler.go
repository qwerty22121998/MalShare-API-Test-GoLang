package main

import (
	"net/http"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"bufio"
	"fmt"
)


func check(err error) {
	if err != nil {
		panic(err)
	}
}


var API_KEY string


func getAllLink() []string {
	url := "https://malshare.com/daily/"
	resp, err := http.Get(url)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	resp.Body.Close()

	var content = string(body)
	plainLinksRegexp, err := regexp.Compile("<a href=\"\\d*-\\d*-\\d*")

	plainLinks := plainLinksRegexp.FindAllString(content, -1)

	for i := range plainLinks {
		dateLink := strings.Replace(plainLinks[i], "<a href=\"", "", -1)
		plainLinks[i] = url  + dateLink + "/malshare_fileList." + dateLink + ".all.txt"
	}
	return plainLinks

}

func getDetail(md5 string) string {
	resp, err := http.Get("https://malshare.com/api.php?api_key=" + API_KEY + "&action=getfile&hash=" + md5)
	fmt.Println("https://malshare.com/api.php?api_key=" + API_KEY + "&action=getfile&hash=" + md5)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	resp.Body.Close()
	return string(body)
}


func getMd5(url string) []string {

	resp, err := http.Get(url)
	check(err)
	scanner :=  bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)
		// ioutil.ReadAll(resp.Body)

	check(err)

	md5s := make([]string, 0)

	for scanner.Scan(){
		md5Current := strings.Split(scanner.Text(), "\t")[0]
		if md5Current == "NULL" {
			continue
		}
		//fmt.Println(md5Current)
		md5s = append(md5s, md5Current)
	}

	resp.Body.Close()
	return md5s

}

func main() {
	API_KEY = os.Args[1]
	tmp := getAllLink()
	// test get the list all md5
	for i, v := range getMd5(tmp[1]) {
		fmt.Printf("Index %d : %s\n", i, v)
	}

}
