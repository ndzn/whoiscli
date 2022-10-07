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
// fix some cctld's whois servers returning unparseable data

// main function to run the program, loops and asks for repeat
func main() {
	for {
		getWhois()
		prompt := promptui.Prompt{
			Label:     "Lookup another domain",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			break
		}
	}
}

// prompts domain and handles basic ui for the program
func getWhois() {
	domain := func(input string) error {
		return nil
	}
	// build the prompt
	prompt := promptui.Prompt{
		Label:    "Domain to lookup",
		Validate: domain,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed")
		return
	}
	// build loading icon
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Returned Whois Data:\n")
	s.Start()

	// makes whoisInfo the result of the whois lookup
	whoisInfo, err := whois.Whois(result)
	if err != nil {
		s.Stop()
		fmt.Println(color.RedString("Invalid domain"))
		return
	}
	// parses the whois data
	parsed, err := whoisparser.Parse(whoisInfo)

	// get the ip from the local dns using the built-in net package
	domainIP, _ := net.LookupIP(parsed.Domain.Domain)
	// return parsed whois data
	s.Stop()
	if err == nil {
		fmt.Println("\033[H\033[2J")
		fmt.Println(color.GreenString("Domain: ") + parsed.Domain.Domain)
		fmt.Println(color.GreenString("Registered @: "), parsed.Registrar.Name)
		fmt.Println(color.GreenString("Created: "), parsed.Domain.CreatedDate)
		fmt.Println(color.GreenString("Expires: "), parsed.Domain.ExpirationDate)
		fmt.Println(color.GreenString("Last Updated: "), parsed.Domain.UpdatedDate)
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
	// build confirmation prompt
	confirmResult, err := confirm.Run()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	if confirmResult == "y" {
		getLocation(domainIP[0].String())
	}
}

// function to get location and information about the webserver from the resolved ip for the domain
func getLocation(ip string) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.FinalMSG = color.GreenString("Website Provider Data:\n")
	s.Start()
	if ip == "" {
		fmt.Println("No IP found")
		return
	}
	// uses the request function to get the json data and prints it out
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
}

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
