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
	whoisparser "github.com/likexian/whois-parser"
	"github.com/manifoldco/promptui"
	"github.com/tidwall/gjson"
	whois "github.com/undiabler/golang-whois"
)

// TODO:
// add selection for what data you want, if you only want registration data, website data, etc.
// fetch from real nameserver and parse

// prompts domain and handles basic ui for the program
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

	whoisInfo, err := whois.GetWhois(result)
	if err != nil {
		fmt.Printf("That isnt a valid domain. Please try again\n")
		return
	} else {
		fmt.Println("Getting whois info...")
	}
	parsed, err := whoisparser.Parse(whoisInfo)
	domainIP, _ := net.LookupIP(parsed.Domain.Domain)
	if err == nil {
		fmt.Println("\033[H\033[2J")
		s.Stop()
		fmt.Println(color.GreenString("Domain: ") + parsed.Domain.Domain)
		fmt.Println(color.GreenString("Registered @: "), parsed.Registrar.Name)
		fmt.Println(color.GreenString("Created: "), parsed.Domain.CreatedDate)
		fmt.Println(color.GreenString("Expires: "), parsed.Domain.ExpirationDate)
		fmt.Println(color.GreenString("Nameservers: "))

		for _, v := range parsed.Domain.NameServers {
			fmt.Println(v)
		}

	} else {
		fmt.Println("This isn't a registered domain. (Is the spelling correct?)")
	}

	confirm := promptui.Prompt{
		Label:     "Show Webserver Data",
		IsConfirm: true,
	}

	confirmResult, err := confirm.Run()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	if confirmResult == "y" {
		getLocation(domainIP[0].String())
	}

}

// function to get location from given ip from domain
func getLocation(ip string) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Website Provider Data:\n")
	s.Start()
	if ip == "" {
		fmt.Println("No IP found")
		return
	}
	json := request("https://ipinfo.io/" + ip + "/json")
	if gjson.Get(json, "status").String() == "404" {
		fmt.Println("No IP found")
	} else {
		s.Stop()
		fmt.Println(color.GreenString("IP: "), gjson.Get(json, "ip").String())
		fmt.Println(color.GreenString("City: "), gjson.Get(json, "city").String())
		fmt.Println(color.GreenString("Region: "), gjson.Get(json, "region").String())
		fmt.Println(color.GreenString("Country: "), gjson.Get(json, "country").String())
		fmt.Println(color.GreenString("Location: "), gjson.Get(json, "loc").String())
		if gjson.Get(json, "hostname").String() == "" {
			fmt.Println(color.GreenString("Hostname: "), "No hostname found")
		} else {
			fmt.Println(color.GreenString("Hostname: "), gjson.Get(json, "hostname").String())
		}
		fmt.Println(color.GreenString("Org: "), gjson.Get(json, "org").String())
	}
} // fix

// helper function to make a request to a web page
func request(link string) string {
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
