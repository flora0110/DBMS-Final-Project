package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	// create a client
	client := http.Client{}

	// set url and create the request
	url := "https://zh.moegirl.org.cn/%E5%8A%A8%E7%94%BB"
	request, err := http.NewRequest("GET", url, nil)
	CheckErr(err)

	// set header
	cookName := &http.Cookie{Name: "stopMobileRedirect", Value: "true"}
	request.AddCookie(cookName)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
	request.Header.Add("sec-ch-ua-platform", "Windows")

	// client sends request
	response, err := client.Do(request)
	CheckErr(err)
	defer response.Body.Close()

	// print the status code
	fmt.Printf("status code: %v\n", response.StatusCode)

	// store the html into data
	data, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))

}

// check error
func CheckErr(err error) {
	defer func() {
		if ins, ok := recover().(error); ok {
			fmt.Println("ERROR: ", ins.Error())
		}
	}()
	if err != nil {
		panic(err)
	}
}