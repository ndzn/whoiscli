package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/manifoldco/promptui"
)

func main() {
	domain := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Domain to lookup",
		Validate: domain,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed")
		return
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Returned Whois Data:\n")

	s.Start()

	whoisInfo, err := whois.Whois(result)
	if err != nil {
		fmt.Printf("That isnt a valid domain. Please try again")
		return
	} else {
		fmt.Println("Getting whois info...")
	}
	parsed, err := whoisparser.Parse(whoisInfo)
	if err == nil {
		fmt.Print("\033[H\033[2J")
		s.Stop()
		fmt.Println(color.GreenString("Domain: ") + parsed.Domain.Domain)
		fmt.Println(color.GreenString("Registered @: "), parsed.Registrar.Name)
		fmt.Println(color.GreenString("Created: "), parsed.Domain.CreatedDate)
		fmt.Println(color.GreenString("Expires: "), parsed.Domain.ExpirationDate)
		// fmt.Println("Website IP", domainIP)
		// json, _ := http.Get("https://ipinfo.io/", net.LookupIP(parsed.Domain.Domain))
		// fmt.Println(json)
		domainIP, _ := net.LookupIP(parsed.Domain.Domain)
		getLocation("20.69.96.99") // testing ip for now
		println(domainIP)
	} else {
		fmt.Println("This isnt a registered domain")
	}

}

// function to get location from given ip
func getLocation(domain string) {
	fmt.Println(readRequest("https://ipinfo.io/" + domain))
} // fix

// helper function to make a request to a web page
func readRequest(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
