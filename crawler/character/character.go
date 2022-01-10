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
	requestAnimation()
}

func requestAnimation() {

	// to read the name of animations
	file, err := os.Open("animation.txt")
	// file buffer
	bs := make([]byte, 8192*8, 8192*8)
	n := -1
	// read file animation.txt and store into bs
	n, err = file.Read(bs)
	// close file animation.txt
	file.Close()

	// get the array of animation name
	animation_name_buff := strings.Split(string(bs[:n]), "\n")
	// to count the num of animation index
	num := 0

	file, err = os.OpenFile("character.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()

	if err == nil {
		// create a client
		client := http.Client{}
		years := [2]string{"2020", "2021"}
		seasons := [4]string{"冬", "春", "夏", "秋"}

		// Cast_Regexp detect casts
		Cast_Regexp := regexp.MustCompile("<span class=\"mw-headline\" id=\"CAST[_0-9]*?\">.*?</span></h3>(\n)?<div class=\"columns-list\" style=\"column-count:2;;;;column-rule-style:none;;\"> \n<ul><li>.*?</li>(\n<li>.*?</li>)*?</ul>")
		// analyze the content that filter with Cast_Regexp
		content_Regexp := regexp.MustCompile("<li>.*?：.*?</li>")
		// analyze the character that filter with content_Regexp
		character_Regexp := regexp.MustCompile(".*：")
		// analyze the cast that filter with content_Regexp
		cast_Regexp := regexp.MustCompile("：.*")

		for _, year := range years {
			for _, season := range seasons {
				// set url and create the request
				url := "https://zh.moegirl.org.cn/日本" + year + "年" + season + "季动画"
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

					// filter data with Cast_Regexp
					Cast_Arr := Cast_Regexp.FindAllString(string(data), -1)


					// print the data in Cast_Arr
					for _, str := range Cast_Arr {
						// process Cast_Arr[i]
						casts := content_Regexp.FindAllString(str, -1)
						for _, data := range casts {
							data = data[4:len(data)-5]
							character := character_Regexp.FindString(data)
							character = character[:len(character)-3]
							cast := cast_Regexp.FindString(data)
							cast = cast[3:]
							fmt.Println(character + ",," + animation_name_buff[num] + ",," + cast)
							file.Write([]byte(character + ",," + animation_name_buff[num] + ",," + cast + "\n"))
						}
						fmt.Println()
						num++
					}
				} else { // if request fail
					fmt.Println("request fail", response.Status)
				}

				// sleep 4 seconds
				time.Sleep(4 * time.Second)
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
