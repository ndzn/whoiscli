package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/manifoldco/promptui"
	"github.com/tidwall/gjson"
)

// TODO:
// add selection for what data you want, if you only want registration data, website data, etc.
// fetch from real nameserver and parse

// prompts domain and handles basic ui for the
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

	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Returned Whois Data:\n")

	s.Start()

	whoisInfo, err := whois.Whois(result)
	if err != nil {
		fmt.Printf("That isnt a valid domain. Please try again\n")
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
		domainIP, _ := net.LookupIP(parsed.Domain.Domain)
		getLocation(domainIP[0].String())
	} else {
		fmt.Println("This isn't a registered domain. (Is the spelling correct?)")
	}

}

// function to get location from given ip from domain
func getLocation(ip string) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Website Provider Data:\n")
	s.Start()
	json := readRequest("https://ipinfo.io/" + ip)
	s.Stop()
	fmt.Println("Website IP", ip)
	fmt.Println("Region:", gjson.Get(json, "region"))
	fmt.Println("City:", gjson.Get(json, "city"))
	fmt.Println("Country:", gjson.Get(json, "country"))
	fmt.Println("Provider:", gjson.Get(json, "org"))
} // fix

// helper function to make a request to a web page
func readRequest(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
