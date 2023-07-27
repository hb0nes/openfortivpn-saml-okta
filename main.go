package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var config *Config

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
	} else if config.Totp || config.Verify || config.Webauthn || config.FastPass {
		headless = true
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
	if config.Totp {
		page.SetDefaultTimeout(7500)
	} else {
		page.SetDefaultTimeout(60000)
	}
	if err != nil {
		log.Fatalf("Could not create page: %v", err)
	}
	return page
}

func totpAsk() (totp string) {
	tries := 3
	var err error
	for i := 1; i <= tries; i++ {
		log.Printf("Enter MFA TOTP (%d/3): ", i)
		_, err = fmt.Scanf("%s", &totp)
		if err != nil {
			log.Printf("Error while reading input: %v", err)
			continue
		}
		match, _ := regexp.Match("^[0-9]{6}$", []byte(totp))
		if !match {
			log.Println("Not a correct TOTP. Should be 6 numbers.")
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
	// Sometimes when using TOTP, it asks for second authentication method anyway.
	time.Sleep(time.Second * 3)
	if err := page.Click("div[data-se='okta_password'] a"); err != nil {
		log.Printf("Could not choose TOTP as authentication method. %v", err)
	}
}

func webauthnChoose(page playwright.Page) {
	if err := page.Click("div[data-se='webauthn'] a"); err != nil {
		log.Printf("Could not choose webauthn as authentication method. %v", err)
	}
	log.Println("Chose webauthn authentication. Press your authenticator device.")
}

func totpInput(page playwright.Page, totp string) {
	log.Println("Waiting for TOTP input...")
	if err := page.Type("input[type='text'][name='credentials.passcode']", totp); err != nil {
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
				log.Printf("Found FortiGate authentication cookie: SVPNCOOKIE=%v", cookie.Value)
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

func verifyChoose(page playwright.Page) {
	if err := page.Click("div[data-se='okta_verify-push'] a"); err != nil {
		log.Printf("Could not choose Okta Verify as authentication method. %v", err)
	}
	log.Printf("Okta Verify selected. Check your Okta Verify app.")
	// Sometimes when using okta verify, it asks for second authentication method anyway.
	// Wait some time for the user to verify access and then enter the password
	time.Sleep(time.Second * 3)
	if err := page.Click("div[data-se='okta_password'] a"); err != nil {
		log.Printf("Could not choose Okta Verify as authentication method. %v", err)
	}
}

func verifySearchChallenge(page playwright.Page) {
	challenge, err := page.TextContent("div[class='number-challenge-section'] span")
	if err != nil {
		log.Println(err)
	}
	log.Printf("Found Okta Verify challenge number: %s", challenge)
}

func fastPassChoose(page playwright.Page) {
	if err := page.Click("div[class='okta-verify-container'] a"); err != nil {
		log.Printf("Could not choose Okta FastPass as authentication method. %v", err)
	}
	if err := page.Click("div[data-se='okta_verify-signed_nonce'] a"); err != nil {
		log.Printf("Could not choose Okta FastPass as authentication method. %v", err)
	}
}

func main() {
	pw := playwrightInit()
	config = configRead()
	page := playwrightGetPage(pw)
	go navigateSAML(page)
	go cookieSearch(page)
	if config.FastPass {
		go fastPassChoose(page)
	} else {
		go usernameInput(page)
		go passwordInput(page)
	}
	if config.Verify {
		go verifyChoose(page)
		go verifySearchChallenge(page)
	}
	if config.Totp {
		totp := totpAsk()
		go totpChoose(page)
		go totpInput(page, totp)
	}
	if config.Webauthn {
		go webauthnChoose(page)
	}
	time.Sleep(time.Second * 60)
	screenshotPath := screenshot(page)
	log.Fatalf("Timed out. See screenshot at %v. Exiting.", screenshotPath)
}
