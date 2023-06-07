package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var config Config

func playwrightInit() (pw *playwright.Playwright) {
	opts := playwright.RunOptions{
		DriverDirectory:     "",
		SkipInstallBrowsers: false,
		Browsers:            []string{"chromium"},
		Verbose:             false,
	}
	err := playwright.Install(&opts)
	if err != nil {
		log.Fatalf("Could not install playwright: %v", err)
	}
	pw, err = playwright.Run()
	if err != nil {
		log.Fatalf("Could not start playwright: %v", err)
	}
	return pw
}

func playwrightGetPage(pw *playwright.Playwright) (page playwright.Page) {
	var headless bool
	if config.Headless != nil {
		headless = *config.Headless
	} else {
		headless = config.Totp
	}
	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(headless),
		},
	)
	if err != nil {
		log.Fatalf("Could not launch browser: %v", err)
	}
	page, err = browser.NewPage(
		playwright.BrowserNewContextOptions{
			IgnoreHttpsErrors: playwright.Bool(true),
		},
	)
	page.SetDefaultTimeout(60000)
	if err != nil {
		log.Fatalf("Could not create page: %v", err)
	}
	return page
}

func totpAsk() (totp int) {
	tries := 3
	var stdin string
	var err error
	for i := 1; i <= tries; i++ {
		log.Printf("Enter MFA TOTP (%d/3): ", i)
		_, err = fmt.Scanf("%s", &stdin)
		if err != nil {
			log.Printf("Error while reading input: %v", err)
			continue
		}
		totp, err = strconv.Atoi(stdin)
		if err != nil {
			log.Printf("Error while reading input: %v", err)
			continue
		}
		break
	}
	if err != nil {
		log.Fatalf("Could not get TOTP: %v", err)
	}
	return totp
}

func totpChoose(page playwright.Page) {
	if err := page.Click("div[data-se='google_otp'] a"); err != nil {
		log.Printf("Could not choose TOTP as authentication method. %v", err)
	}
}

func totpInput(page playwright.Page, totp int) {
	log.Println("Waiting for TOTP input...")
	if err := page.Type("input[type='text'][name='credentials.passcode']", strconv.Itoa(totp)); err != nil {
		log.Printf("Could not find TOTP input in time: %v.", err)
	}
	log.Println("Submitting TOTP...")
	if err := page.Click("input[type='submit']"); err != nil {
		log.Printf("Could not submit TOTP: %v", err)
	}
}

func navigateSAML(page playwright.Page) {
	log.Println("Navigating to https://vpn.deribit.com")
	if _, err := page.Goto("https://vpn.deribit.com"); err != nil {
		log.Fatalf("Could not navigate: %v", err)
	}
	log.Println("Clicking SAML login button...")
	if err := page.Click("button#saml-login-bn"); err != nil {
		log.Printf("Could not find the SAML login button on the FortiGate page in time. button ID: button#saml-login-bn. %v", err)
	}
}

func cookieSearch(page playwright.Page) {
	log.Println("Searching for FortiGate authentication cookie...")
	for {
		cookies, err := page.Context().Cookies()
		if err != nil {
			log.Printf("Could not get cookies from context: %v", err)
		}
		for _, cookie := range cookies {
			if cookie.Name == "SVPNCOOKIE" {
				log.Println("Found FortiGate authentication cookie: \n")
				fmt.Printf("SVPNCOOKIE=%s", cookie.Value)
				os.Exit(0)
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func usernameInput(page playwright.Page) {
	log.Println("Waiting for username input...")
	if err := page.Type("input[autocomplete='username']", config.Username); err != nil {
		log.Printf("Could not find username input in time: %v", err)
	}
	log.Println("Submitting username...")
	if err := page.Click("input[type='submit']"); err != nil {
		log.Printf("Could not submit username: %v", err)
	}
}

func passwordInput(page playwright.Page) {
	log.Println("Waiting for password input...")
	if err := page.Type("input[type='password'][name='credentials.passcode']", config.Password); err != nil {
		log.Printf("Could not find password input in time: %v.", err)
	}
	log.Println("Submitting password...")
	if err := page.Click("input[type='submit']"); err != nil {
		log.Printf("Could not submit password: %v", err)
	}
}

func screenshot(page playwright.Page) (screenshotPath string) {
	screenshotPath, _ = os.UserHomeDir()
	screenshotPath = filepath.Join(screenshotPath, "openfortivpn-saml.png")
	if _, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(screenshotPath),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}
	return
}

func main() {
	pw := playwrightInit()
	config = configRead()
	page := playwrightGetPage(pw)
	if config.Totp {
		totp := totpAsk()
		go totpChoose(page)
		go totpInput(page, totp)
	}
	go navigateSAML(page)
	go cookieSearch(page)
	go usernameInput(page)
	go passwordInput(page)
	time.Sleep(time.Second * 45)
	screenshotPath := screenshot(page)
	log.Fatalf("Timed out after 45s. See screenshot at %v. Exiting.", screenshotPath)
}
