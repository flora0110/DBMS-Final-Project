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
		Voice_Regexp := regexp.MustCompile("<td width=\"[0-9]*?px\" bgcolor=\"#[a-fA-F0-9]*?\">(\n)?所(属|屬)公司(\n)?</td>(\n)?<td>(\n)?(<a href=\"/.*?\" title=\".*?\">.*?</a>|.*?)(\n)?</td>")

		// analyze the content that filter with Voice_Regexp
		content_Regexp_1 := regexp.MustCompile("<a href=\"/.*?\" title=\".*?\">.*?</a>")
		sub_content_Regexp := regexp.MustCompile(">.*?</a>")

		content_Regexp_2 := regexp.MustCompile("<td>(\n)?.*?(\n)?</td>")

		// file buffer
		bs := make([]byte, 8192*8, 8192*8)
		n := -1
		// read file voice_buff.txt and store into bs
		n, err = file.Read(bs)
		// close file voice_buff.txt
		file.Close()


		// get the array of voice actor/actress
		voice_buff := strings.Split(string(bs[:n]), "\n")

		// voice.txt detect
		file, blockErr := os.Open("voice.txt")
		// initialize have_done
		have_done := -1
		// forCount store which have done
		forCount := []string{}
		// if voice.txt exist
		if blockErr == nil {
			// relloc bs
			bs = make([]byte, 16384*8, 16384*8)
			// read file voice.txt and store into bs
			n, err = file.Read(bs)
			// close file voice.txt
			file.Close()
			// forCount store which have done
			forCount = strings.Split(string(bs[:n]), "\n")
			// set have_done
			have_done = have_done + len(forCount)
		}

		// open or create file voice.txt
		file, err = os.OpenFile("voice.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer file.Close()

		// process datas
		for i, value := range voice_buff {

			if i >= have_done && len(value) != 0 { /* haven't done */
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
						if content_Regexp_1.FindString(str) != "" {
							// 2nd filter
							str = content_Regexp_1.FindString(str)
							// 3rd filter
							str = sub_content_Regexp.FindString(str)

							// modify str
							if str[len(str)-5] == '\n' {
								str = str[1:len(str)-5]
							} else {
								str = str[1:len(str)-4]
							}
						} else {
							// 2nd filter
							str = content_Regexp_2.FindString(str)
							// modify str
							if str[0] == '\n' {
								str = str[1:]
							}
							if str[len(str)-6] == '\n' {
								str = str[4:len(str)-6]
							} else {
								str = str[4:len(str)-5]
							}
						}
						str = value[0:len(value)-1] + "," + str + "\n"
					} else {
						str = value + "\n"
					}

					// write the data into file voice.txt
					file.Write([]byte(str))
				} else if response.StatusCode == 404 {
					// tag it's 404
					str := value[0: len(value)-1] + " -----------> 404 <-----------\n"
					// write the data into file voice.txt
					file.Write([]byte(str))
				} else {
					fmt.Println("Blocking QAQ or Robot detecting ...")
					break
				}

			} else { /* have done */
				fmt.Printf("the %dth case have done, print out\n", i)
				// write the data back into file voice.txt
				file.Write([]byte(forCount[i] + "\n"))
			}

			// sleep Zzzz
			if i >= have_done && (i % 20 == 0 || i % 27 == 0) {
				// sleep 1 seconds
				if i % 7 == 0 {
					time.Sleep(7 * time.Second)
				} else if i % 3 == 0 {
					time.Sleep(4 * time.Second)
				} else {
					time.Sleep(2 * time.Second)
				}
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
