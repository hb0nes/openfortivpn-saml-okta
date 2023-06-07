package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"strconv"
	"time"
)

var config Config

func searchCookie(page playwright.Page) {
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

func enterUsername(page playwright.Page) {
	log.Println("Waiting for username input...")
	if err := page.Type("input[autocomplete='username']", config.Username); err != nil {
		log.Printf("Could not find username input in time: %v", err)
	}
	log.Println("Submitting username...")
	if err := page.Click("input[type='submit']"); err != nil {
		log.Printf("Could not submit username: %v", err)
	}
}

func enterPassword(page playwright.Page) {
	log.Println("Waiting for password input...")
	if err := page.Type("input[type='password']", config.Password); err != nil {
		log.Printf("Could not find password input in time: %v.", err)
	}
	log.Println("Submitting password...")
	if err := page.Click("input[type='submit']"); err != nil {
		log.Printf("Could not submit password: %v", err)
	}
}

func askTOTP() (totp int, err error) {
	tries := 3
	var stdin string
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
	return totp, err
}

func chooseTOTP(page playwright.Page) {
	if err := page.Click("div[data-se='google_otp'] a"); err != nil {
		log.Printf("Could not choose TOTP as authentication method. %v", err)
	}
}

func enterTOTP(page playwright.Page, totp int) {
	log.Println("Waiting for TOTP input...")
	if err := page.Type("input[type='text'][name='credentials.passcode']", strconv.Itoa(totp)); err != nil {
		log.Printf("Could not find TOTP input in time: %v.", err)
	}
	log.Println("Submitting TOTP...")
	if err := page.Click("input[type='submit']"); err != nil {
		log.Printf("Could not submit TOTP: %v", err)
	}
}

func main() {
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
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Could not start playwright: %v", err)
	}
	config = readConfig()
	var totp int
	if config.Totp {
		totp, err = askTOTP()
		if err != nil {
			log.Fatalf("Could not get TOTP: %v", err)
		}
	}
	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(config.Totp),
		},
	)
	if err != nil {
		log.Fatalf("Could not launch browser: %v", err)
	}
	page, err := browser.NewPage(
		playwright.BrowserNewContextOptions{
			IgnoreHttpsErrors: playwright.Bool(true),
		},
	)
	if err != nil {
		log.Fatalf("Could not create page: %v", err)
	}
	if config.Totp {
		go chooseTOTP(page)
		go enterTOTP(page, totp)
	}
	log.Println("Navigating to https://vpn.deribit.com")
	if _, err = page.Goto("https://vpn.deribit.com"); err != nil {
		log.Fatalf("Could not navigate: %v", err)
	}
	log.Println("Clicking SAML login button...")
	if err = page.Click("button#saml-login-bn"); err != nil {
		log.Printf("Could not find the SAML login button on the FortiGate page in time. button ID: button#saml-login-bn. %v", err)
	}

	go searchCookie(page)
	go enterUsername(page)
	go enterPassword(page)
	time.Sleep(time.Second * 30)
	log.Fatal("Timed out after 30s. Exiting.")
}
