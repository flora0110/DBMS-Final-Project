package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	requestVoice()
}

func requestVoice() {

	file, err := os.Open("voice_buff.txt")

	if err != nil {

		file, err := os.OpenFile("voice_buff.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer file.Close()

		// create a client
		client := http.Client{}

		// Voice_Regexp detect gender or voice actor/actress
		Voice_Regexp := regexp.MustCompile("(<b>男</b>)|(<b>女</b>)|(\n<a href=\"/([A-Z0-9%]*)?\" title=\"(.*)?\">(.*)?</a>)")
		// analyze the content that filter with Voice_Regexp
		content_Regexp := regexp.MustCompile(">(.*)<")

		// set url and create the request
		url := "https://zh.moegirl.org.cn/%E5%A3%B0%E4%BC%98"
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

		// if request success
		if response.StatusCode == 200 {
			// store the html into data
			data, err := ioutil.ReadAll(response.Body)
			fmt.Println("request success")
			CheckErr(err)

			// filter data with voice_arr
			voice_arr := Voice_Regexp.FindAllString(string(data), -1)

			// initialize gender
			gender := "女"
			for _, str := range voice_arr {
				// process voice_arr[i]
				str = content_Regexp.FindString(str)
				str = str[1 : len(str)-1]

				if str == "男" {
					gender = "男"
				} else if str == "女" {
					gender = "女"
				} else { // if it's voice
					str = str + "," + gender + "\n"
					file.Write([]byte(str))
				}
			}
		} else { // if request fail
			fmt.Println("request fail", response.Status)
		}

	} else {

		bs := make([]byte, 8192*8, 8192*8)
		n := -1
		n, err = file.Read(bs)
		voice_buff := strings.Split(string(bs[:n]), "\n")

		for _, value := range voice_buff {
			if len(value) != 0 {
				values := strings.Split(value, ",")
				url := "https://zh.moegirl.org.cn/" + values[0]
				fmt.Println(url)
			}
		}

		file.Close()
	}
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
