package main

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/log"
	"strings"
	"time"
)

type Headers struct {
	AccessControlAllowCredentials string `json:"access-control-allow-credentials"`
	AccessControlAllowOrigin      string `json:"access-control-allow-origin"`
	AccessControlExposeHeaders    string `json:"access-control-expose-headers"`
	ContentEncoding               string `json:"content-encoding"`
	ContentType                   string `json:"content-type"`
	Date                          string `json:"date"`
	Server                        string `json:"server"`
	Vary                          string `json:"vary"`
}

type SecurityDetails struct {
	CertificateId                     int      `json:"certificateId"`
	CertificateTransparencyCompliance string   `json:"certificateTransparencyCompliance"`
	Cipher                            string   `json:"cipher"`
	EncryptedClientHello              bool     `json:"encryptedClientHello"`
	Issuer                            string   `json:"issuer"`
	KeyExchange                       string   `json:"keyExchange"`
	KeyExchangeGroup                  string   `json:"keyExchangeGroup"`
	Protocol                          string   `json:"protocol"`
	SanList                           []string `json:"sanList"`
	ServerSignatureAlgorithm          int      `json:"serverSignatureAlgorithm"`
	SignedCertificateTimestampList    []string `json:"signedCertificateTimestampList"`
	SubjectName                       string   `json:"subjectName"`
	ValidFrom                         int64    `json:"validFrom"`
	ValidTo                           int64    `json:"validTo"`
}

type Timing struct {
	ConnectEnd               float64 `json:"connectEnd"`
	ConnectStart             float64 `json:"connectStart"`
	DnsEnd                   float64 `json:"dnsEnd"`
	DnsStart                 float64 `json:"dnsStart"`
	ProxyEnd                 float64 `json:"proxyEnd"`
	ProxyStart               float64 `json:"proxyStart"`
	PushEnd                  float64 `json:"pushEnd"`
	PushStart                float64 `json:"pushStart"`
	ReceiveHeadersEnd        float64 `json:"receiveHeadersEnd"`
	ReceiveHeadersStart      float64 `json:"receiveHeadersStart"`
	RequestTime              float64 `json:"requestTime"`
	SendEnd                  float64 `json:"sendEnd"`
	SendStart                float64 `json:"sendStart"`
	SslEnd                   float64 `json:"sslEnd"`
	SslStart                 float64 `json:"sslStart"`
	WorkerFetchStart         float64 `json:"workerFetchStart"`
	WorkerReady              float64 `json:"workerReady"`
	WorkerRespondWithSettled float64 `json:"workerRespondWithSettled"`
	WorkerStart              float64 `json:"workerStart"`
}

type Response struct {
	AlternateProtocolUsage string          `json:"alternateProtocolUsage"`
	ConnectionId           int             `json:"connectionId"`
	ConnectionReused       bool            `json:"connectionReused"`
	EncodedDataLength      int             `json:"encodedDataLength"`
	FromDiskCache          bool            `json:"fromDiskCache"`
	FromPrefetchCache      bool            `json:"fromPrefetchCache"`
	FromServiceWorker      bool            `json:"fromServiceWorker"`
	Headers                Headers         `json:"headers"`
	MimeType               string          `json:"mimeType"`
	Protocol               string          `json:"protocol"`
	RemoteIPAddress        string          `json:"remoteIPAddress"`
	RemotePort             int             `json:"remotePort"`
	ResponseTime           float64         `json:"responseTime"`
	SecurityDetails        SecurityDetails `json:"securityDetails"`
	SecurityState          string          `json:"securityState"`
	Status                 int             `json:"status"`
	StatusText             string          `json:"statusText"`
	Timing                 Timing          `json:"timing"`
	URL                    string          `json:"url"`
}

type Params struct {
	FrameID      string   `json:"frameId"`
	HasExtraInfo bool     `json:"hasExtraInfo"`
	LoaderID     string   `json:"loaderId"`
	RequestID    string   `json:"requestId"`
	Response     Response `json:"response"`
	Timestamp    float64  `json:"timestamp"`
	Type         string   `json:"type"`
}

type Message struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Root struct {
	Message Message `json:"message"`
	WebView string  `json:"webview"`
}

func CollectLootboxURLs() []string {
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
	scrollStep := 500.0 // Scroll down by 100 pixels
	for scrollTop := 0.0; scrollTop < 100000.0; scrollTop += scrollStep {
		scrollScript := fmt.Sprintf("window.scrollTo(0, %f);", scrollTop)
		_, err := driver.ExecuteScript(scrollScript, nil)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 1000) // Give some time for the page to settle
	}
	time.Sleep(time.Second * 1) // Give some time for the page to load

	// Get browser logs including network events
	logs, err := driver.Log(log.Performance)
	if err != nil {
		panic(err)
	}

	// Find URLs containing "vortex" in logs
	var urls []string
	for _, l := range logs {
		var message Root
		err := json.Unmarshal([]byte(l.Message), &message)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
		if message.Message.Method != "Network.responseReceived" {
			continue
		}
		//fmt.Printf("%s\n", l.Message)
		url := message.Message.Params.Response.URL

		if contains := strings.Contains(url, "https://vortex.worldofwarships.com/api/get_lootbox/"); contains {
			urls = append(urls, url)
		}
	}

	// Print vortex URLs
	fmt.Println("URLs containing 'vortex':")
	for _, url := range urls {
		fmt.Println(url)
	}
    return urls

}
