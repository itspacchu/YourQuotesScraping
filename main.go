package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var DownloadedFiles = make(map[string]bool)

func buildFileName(fullUrlFile string) string {
	fileUrl, e := url.Parse(fullUrlFile)
	if e != nil {
		panic(e)
	}

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	return fileName
}

func downloadImages(Httpurl string) {
	bname := buildFileName(Httpurl)
	if _, ok := DownloadedFiles[bname]; ok {
		return
	}
	resp, err := http.Get(Httpurl)
	if err != nil {
		fmt.Printf("[ERR] Failed to download %s ... Skipping\n", Httpurl)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[INFO] Downloading image %s\n", bname)
	f, _ := os.Create("./images/" + bname)
	defer f.Close()
	f.ReadFrom(resp.Body)
	DownloadedFiles[bname] = true
}

func main() {
	// https://stackoverflow.com/questions/70867224/selenium-go-how-to
	service, err := selenium.NewChromeDriverService("/usr/bin/chromedriver", 4444)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		"--headless",
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[INFO] Getting profile %s\n", os.Getenv("QUOTE_URL"))
	driver.Get(os.Getenv("QUOTE_URL"))
	time.Sleep(5)
	currUrl, _ := driver.CurrentURL()
	fmt.Printf("[Chrome] Opened url for %s", currUrl)
	grid_button, err := driver.FindElement(selenium.ByXPATH, "/html/body/div[1]/div/div/div[3]/div/div[1]/main/div/div/div/div[3]/div[1]/div/div[2]/div[1]")
	if err != nil {
		fmt.Printf("Unable to find the grid button")
	}
	grid_button.Click()

	thegrid, _ := driver.FindElement(selenium.ByXPATH, "/html/body/div[1]/div/div/div[3]/div/div[1]/main/div/div/div/div[3]/div[2]")
	if err != nil {
		fmt.Printf("Unable to find the grid of images")
	}

	// start scrolling
	for i := 0; ; i += 20 {
		pageHeight, err := driver.ExecuteScript("return document.body.scrollHeight", nil)
		driver.ExecuteScript(fmt.Sprintf("window.scrollTo(0,%d)", i), nil)
		time.Sleep(10)
		newHeight, err := driver.ExecuteScript("return document.body.scrollHeight", nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[INFO] ------- Scrolling %f -- %f ... %d\n", pageHeight.(float64), newHeight.(float64), i)
		if pageHeight.(float64) >= newHeight.(float64) {
			elements, _ := thegrid.FindElements(selenium.ByCSSSelector, "img")
			for _, element := range elements {
				src, _ := element.GetAttribute("src")
				if len(src) > 1 {
					downloadImages(src)
				}
			}
		}
	}
}
