package main

import (
    "github.com/tebeka/selenium"
    "github.com/tebeka/selenium/chrome"
    "github.com/tebeka/selenium/log"
    "strings"
    "time"
    "fmt"
)

func main() {
    // Run Chrome browser
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
    //    "--headless",  // comment out this line to see the browser
    }})
    caps.SetLogLevel(log.Performance, log.All)

    driver, err := selenium.NewRemote(caps, "")
    if err != nil {
        panic(err)
    }
    defer driver.Quit()

    driver.Get("https://worldofwarships.com/en/content/contents-and-drop-rates-of-containers/")
	// Find and click the GDPR accept button
	gdprAcceptButton, err := driver.FindElement(selenium.ByID, "onetrust-accept-btn-handler")
	if err != nil {
		panic(err)
	}
	gdprAcceptButton.Click()
	time.Sleep(time.Second * 2) // Give some time for the page to load
	// Scroll down the page

	// Scroll down the page progressively
	//pageHeightScript := "return Math.max( document.body.scrollHeight, document.body.offsetHeight, document.documentElement.clientHeight, document.documentElement.scrollHeight, document.documentElement.offsetHeight );"
	//pageHeight, err := driver.ExecuteScript(pageHeightScript, nil)
	//if err != nil {
	//	panic(err)
	//}

	scrollStep := 500.0 // Scroll down by 100 pixels
	for scrollTop := 0.0; scrollTop < 100000.0; scrollTop += scrollStep {
		scrollScript := fmt.Sprintf("window.scrollTo(0, %f);", scrollTop)
		_, err := driver.ExecuteScript(scrollScript, nil)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 500) // Give some time for the page to settle
	}
	time.Sleep(time.Second * 1) // Give some time for the page to load


	// Get browser logs including network events
	logs, err := driver.Log(log.Performance)
	if err != nil {
		panic(err)
	}

	// Find URLs of loaded resources
//	var resourceURLs []string
//	for _, l := range logs {
//        fmt.Printf("%s\n", l.Message)
//		if strings.Contains(l.Message, "\"method\":\"Network.responseReceived\"") {
//			// Parse JSON response to extract URL
//			urlStart := strings.Index(l.Message, "\"url\":\"") + len("\"url\":\"")
//			urlEnd := strings.Index(l.Message[urlStart:], "\"")
//			resourceURLs = append(resourceURLs, l.Message[urlStart:urlStart+urlEnd])
//		}
//	}
//
//	// Print loaded resource URLs
//	fmt.Println("Loaded resource URLs:")
//	for _, url := range resourceURLs {
//		fmt.Println(url)
//	}

//	// Capture browser logs
//	logs, err := driver.Log(log.Browser)
//	if err != nil {
//		panic(err)
//	}
//
//	// Find URLs containing "vortex" in logs
	var vortexURLs []string
	for _, l := range logs {
        fmt.Printf("%s\n", l.Message)
		if containsVortex := strings.Contains(l.Message, "vortex.worldofwarships.com"); containsVortex {
			vortexURLs = append(vortexURLs, l.Message)
		}
	}

	// Print vortex URLs
	fmt.Println("URLs containing 'vortex':")
	for _, url := range vortexURLs {
		fmt.Println(url)
	}

	time.Sleep(time.Second * 1)
}

