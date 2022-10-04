package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/manifoldco/promptui"
	"github.com/tidwall/gjson"
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

	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Returned Whois Data:\n")

	s.Start()

	whoisInfo, err := whois.Whois(result)
	if err != nil {
		fmt.Printf("That isnt a valid domain. Please try again\n")
		return
	} else {
		fmt.Println("Getting whois info...\n")
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
		//domainIP, _ := net.LookupIP(parsed.Domain.Domain)
		// fmt.Println(color.YellowString("Data from Domain's website"))
		getLocation("8.8.8.8") // testing ip for now
	} else {
		fmt.Println("This isn't a registered domain. (Is the spelling correct?)")
	}

}

// function to get location from given ip from domain
func getLocation(ip string) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Website provider data:\n")
	s.Start()
	json := readRequest("https://ipinfo.io/" + ip)
	s.Stop()
	// fmt.Println(json)
	// region := gjson.Get(json, "region")
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
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
