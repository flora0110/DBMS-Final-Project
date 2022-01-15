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
	requestVoiceCompany()
}

func requestVoiceCompany() {
	// create a client
	client := http.Client{}

	file, err := os.Open("voice_company_link.txt")

	if err != nil {
		file, err := os.OpenFile("voice_company_link.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer file.Close()

		// set url and create the request
		url := "https://zh.moegirl.org.cn/index.php?title=Template:%E6%97%A5%E6%9C%AC%E5%A3%B0%E4%BC%98%E4%BA%8B%E5%8A%A1%E6%89%80&variant=zh-tw&mobileaction=toggle_view_desktop"
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

		companyNiHon_Regexp := regexp.MustCompile("<tr><td class=\"navbox-group\" style=\";padding:0 1em;;\">聲優事務所</td><td style=\"text-al(\n)?ign:left;border-left:2px solid #fdfdfd;width:100%;padding:0px;;;\" class=\"navbox-list navbox-odd\"><div style=\"padding:0 0.25em\">(.*)?</tr>")
		link_Regexp := regexp.MustCompile("href=\"/([A-Za-z0-9-%()_?.&;=]*)?\" title")

		// store the html into data
		data, err := ioutil.ReadAll(response.Body)
		str := string(data)

		str = companyNiHon_Regexp.FindString(str)
		links := link_Regexp.FindAllString(str, -1)

		for _, s := range links {
			s = "https://zh.moegirl.org.cn" + s[6:len(s)-7] + "\n"
			fmt.Print(s)
			file.Write([]byte(s))
		}
	} else {
		// file buffer
		bs := make([]byte, 8192*8, 8192*8)
		n := -1
		// read file voice_company_link.txt and store into bs
		n, err = file.Read(bs)
		// close file voice_company_link.txt
		file.Close()

		// get the array of voice company link
		voice_company_link_buff := strings.Split(string(bs[:n]), "\n")

		// RE
		data_Regexp := regexp.MustCompile("<td style=\".*?;\" bgcolor=\"#[A-Fa-f0-9]*?\">(\n)?(((名稱|名称)</td>(\n)?<td>.*?)|((網址|网址)</td>(\n)?<td><a target=\".*?\" rel=\".*?\" class=\".*?\" href=\".*?\">.*?</a>))(\n)?</td>")
		data_name_position_Regexp := regexp.MustCompile("<td>.*?(\n)?</td>")
		name_filter := regexp.MustCompile(">.*?<")
		data_link_Regexp := regexp.MustCompile("href=\".*?\">")

		// voice_company.txt detect
		file, blockErr := os.Open("voice_company.txt")
		// initialize have_done
		have_done := -1
		// forCount store which have done
		forCount := []string{}
		// if voice_company.txt exist
		if blockErr == nil {
			// relloc bs
			bs = make([]byte, 16384*8, 16384*8)
			// read file voice_company.txt and store into bs
			n, err = file.Read(bs)
			// close file voice_company.txt
			file.Close()
			// forCount store which have done
			forCount = strings.Split(string(bs[:n]), "\n")
			// set have_done
			have_done = have_done + len(forCount)
		}

		// open or create file voice_company.txt
		file, err = os.OpenFile("voice_company.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer file.Close()

		for i, link := range voice_company_link_buff {

			if i >= have_done && len(link) != 0 { /* haven't done */
				// set url and create the request
				url := link
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

				fmt.Printf("the %dth case, status code: %v\n", i, response.StatusCode)

				// store the html into data
				data, err := ioutil.ReadAll(response.Body)
				str := string(data)
				datas := data_Regexp.FindAllString(str, -1)

				data_str := ""
				for _, data := range datas {
					link_bool, _ := regexp.MatchString("(網址|网址)", data)
					if link_bool {
						data = data_link_Regexp.FindString(data)
						data = data[6 : len(data)-2]
						data_str += (",," + data)
					} else {
						data = data_name_position_Regexp.FindString(data)
						if data[len(data)-6] == '\n' {
							data = data[4 : len(data)-6]
						} else {
							data = data[4 : len(data)-5]
						}
						tag_bool, _ := regexp.MatchString(">(.*)?<", data)
						if tag_bool {
							data = name_filter.FindString(data)
							data = data[1 : len(data)-1]
						}
						data_str += ("\"" + data + "\"")
					}
				}
				file.Write([]byte(data_str + "\n"))
			} else { /* have done */
				fmt.Printf("the %dth case have done, print out\n", i)
				// write the data back into file voice_company.txt
				file.Write([]byte(forCount[i] + "\n"))
			}

			// sleep Zzzz
			if i >= have_done && (i%20 == 0 || i%27 == 0) {
				// sleep 1 seconds
				if i%7 == 0 {
					time.Sleep(7 * time.Second)
				} else if i%3 == 0 {
					time.Sleep(4 * time.Second)
				} else {
					time.Sleep(3 * time.Second)
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
