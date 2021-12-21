package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	requestVoice()
}

func requestVoice() {

	// create a client
	client := http.Client{}

	file, err := os.Open("voice_buff.txt")

	if err != nil {

		file, err := os.OpenFile("voice_buff.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer file.Close()

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

		// Voice_Regexp detect gender or voice actor/actress
		Voice_Regexp := regexp.MustCompile("<td width=\"([0-9]*)?px\" bgcolor=\"#([A-F0-9]*)?\">所属公司</td>\n<td><a href=\"/(.*)?\" title=\"(.*)?\">(.*)?</a>\n</td>")

		// analyze the content that filter with Voice_Regexp
		content_Regexp := regexp.MustCompile("<a href=\"/(.*)?\" title=\"(.*)?\">(.*)?</a>")
		sub_content_Regexp := regexp.MustCompile(">(.*)?</a>")

		// file buffer
		bs := make([]byte, 8192*8, 8192*8)
		n := -1
		// read file voice_buff.txt and store into bs
		n, err = file.Read(bs)
		// close file voice_buff.txt
		file.Close()

		// open or create file voice.txt
		file, err = os.OpenFile("voice.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer file.Close()

		// get the array of voice actor/actress
		voice_buff := strings.Split(string(bs[:n]), "\n")


		for i, value := range voice_buff {
			if len(value) != 0 {
				// get name(values[0]) and gender(values[1])
				values := strings.Split(value, ",")
				// set url of voice actor/actress 's page
				url := "https://zh.moegirl.org.cn/" + values[0]
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
				fmt.Printf("the %dth case status code: %v\n", i, response.StatusCode)

				// if request success
				if response.StatusCode == 200 {
					// store the html into data
					data, err := ioutil.ReadAll(response.Body)
					CheckErr(err)
					// 1st filter of company
					str := Voice_Regexp.FindString(string(data))
					if(len(str) != 0) {
						// 2nd filter
						str = content_Regexp.FindString(str)
						// 3rd filter
						str = sub_content_Regexp.FindString(str)

						// modify str
						str = str[1:len(str)-4]
						str = value[0:len(value)-1] + "," + str + "\n"
					} else {
						// TODO: detect whether add "\n"
						str = value + "\n"
					}

					// write the data into file voice.txt
					file.Write([]byte(str))
				}
			}
			// sleep
			if(i % 20 == 0 || i % 27 == 0) {
				// sleep 1 seconds
				time.Sleep(1 * time.Second)
				fmt.Println("Sleep Over.....")
			}
		}
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
