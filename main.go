package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/jreisinger/checkip/check"
)

func init() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)
}

func main() {
	ips := parseIPs(getIPs())

	var locations []*Location

	for _, ip := range ips {
		res, err := check.DBip(ip)
		if err != nil {
			log.Printf("while getting geolocation of %s: %v", ip, err)
			continue
		}
		js, err := res.Info.JsonString()
		if err != nil {
			log.Printf("while getting info for %s: %v", ip, err)
			continue
		}
		l := Location{IP: ip}
		if err := json.Unmarshal([]byte(js), &l); err != nil {
			log.Printf("while unmarshalling info for %s: %v", ip, err)
			continue
		}
		locations = append(locations, &l)
	}

	sort.Sort(customSort{locations, func(x, y *Location) bool {
		if x.Country != y.Country {
			return x.Country < y.Country
		}
		if x.City != y.City {
			return x.City < y.City
		}
		return false
	}})

	printLocations(locations)
}

type Location struct {
	IP      net.IP
	Country string
	City    string
}

type customSort struct {
	l    []*Location
	less func(x, y *Location) bool
}

func (x customSort) Len() int           { return len(x.l) }
func (x customSort) Less(i, j int) bool { return x.less(x.l[i], x.l[j]) }
func (x customSort) Swap(i, j int)      { x.l[i], x.l[j] = x.l[j], x.l[i] }

func printLocations(locations []*Location) {
	const format = "%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, l := range locations {
		// you don't have to derefence here like (*l).Country
		fmt.Fprintf(tw, format, l.Country, l.City, l.IP)
	}
	tw.Flush() // calculate column widths and print table
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
			log.Printf("can't parse IP address: %s", ip)
			continue
		}
		ipsParsed = append(ipsParsed, ipParsed)
	}
	return ipsParsed
}
