package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip"
)

/*
1. Get country and city of each IP address.
2. Print counts per city showing also country.
*/

type Counts map[string]map[string][]net.IP

func main() {
	ips := parseIPs(getIPs())

	counts := make(map[string]map[string][]net.IP)

	for _, ip := range ips {
		g := &checkip.Geo{}
		_, err := g.Check(ip)
		if err != nil {
			log.Printf("while getting geolocation: %v", err)
			continue
		}
		if _, ok := counts[g.Country]; !ok {
			counts[g.Country] = make(map[string][]net.IP)
		}
		counts[g.Country][g.City] = append(counts[g.Country][g.City], ip)
	}

	for country, m := range counts {
		for city, ips := range m {
			fmt.Printf("%s;%s;%d\n", country, city, len(ips))
		}
	}
}

func getIPs() []string {
	var ips []string

	if len(os.Args[1:]) == 0 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			ips = append(ips, s.Text())
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		ips = os.Args[1:]
	}

	return ips
}

func parseIPs(ips []string) []net.IP {
	var ipsParsed []net.IP
	for _, ip := range ips {
		ipParsed := net.ParseIP(ip)
		if ipParsed == nil {
			log.Printf("Can't parse IP address: %s", ip)
			continue
		}
		ipsParsed = append(ipsParsed, ipParsed)
	}
	return ipsParsed
}
