package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	requestAnimation()
}

func requestAnimation() {

	file, err := os.OpenFile("animation.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()

	if err == nil {
		// create a client
		client := http.Client{}
		years := [2]string{"2020", "2021"}
		seasons := [4]string{"冬", "春", "夏", "秋"}

		// count the number of animations
		num := 0

		// Animation_Regexp detect animation or director
		Animation_Regexp := regexp.MustCompile("(<span class=\"mw-headline\" id=\".*?\">.*?</span>)|(<li>(总|系列)?(监督|監督)[^輔佐]*?：.*?</li>)|(<li>(動畫製作|动画制作)：.*?</li>)")
		// analyze the content that filter with Animation_Regexp
		content_Regexp := regexp.MustCompile(">(.*)<")
		// get the name of director or company
		director_company_Regexp := regexp.MustCompile("：(.*)")

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

					// filter data with Animation_Regexp
					Animation_Arr := Animation_Regexp.FindAllString(string(data), -1)

					// print the data in Animation_Arr
					for _, str := range Animation_Arr {
						// process Animation_Arr[i]
						str = content_Regexp.FindString(str)
						str = str[1 : len(str)-1]

						// detect whether it's the data of an animation
						if !(str == "简介" || str == "簡介" || str == "STAFF" || str == "CAST" || str == "导航" || str == "参见" || str == "參見" || str == "導航") {
							director_bool, _ := regexp.MatchString("(总|系列)?(监督|監督)(.*)?：(.*)?", str)
							company_bool, _ := regexp.MatchString("(動畫製作|动画制作)：(.*)?", str)

							if director_bool {
								fmt.Println(str)
								str = director_company_Regexp.FindString(str)
								str = str[3:]
								file.Write([]byte(",,D:" + str))
							} else if company_bool {
								fmt.Println(str)
								str = director_company_Regexp.FindString(str)
								str = str[3:]
								file.Write([]byte(",,C:" + str))
							} else {
								fmt.Println(num, "----------")
								num++
								fmt.Println(year, season, "\n", str)
								file.Write([]byte("\n" + year + ",," + season + ",," + str))
							}
						}
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
