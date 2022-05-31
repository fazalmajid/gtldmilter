// based on https://github.com/phalaaxx/pf-milters bogomilter
package main

import (
	"bufio"
	"flag"
	//"fmt"
	"github.com/phalaaxx/milter"
	"log"
	"net"
	"os"
	"strings"
)

/* global variables */
var SuspiciousGTLD map[string]bool
var SuspiciousDests map[string]bool

/* GtldMilter object */
type GtldMilter struct {
	milter.Milter
	from string
}

/* MailFrom is called on envelope from address */
func (b *GtldMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	// save from address for later reference
	b.from = from
	log.Println("received email from ", b.from)
	return milter.RespContinue, nil
}

/* RcptTo is called on envelope from address */
func (b *GtldMilter) RcptTo(rcptTo string, m *milter.Modifier) (milter.Response, error) {
	log.Println("received email from ", b.from, " to ", rcptTo)
	components := strings.Split(b.from, ".")
	clen := len(components)
	if clen < 1 {
		return milter.RespContinue, nil
	}
	tld := components[clen-1]
	if SuspiciousGTLD[tld] && SuspiciousDests[rcptTo] {
		log.Println("reject email from ", b.from, " to ", rcptTo)
		return milter.RespReject, nil
	}
	return milter.RespContinue, nil
}

/* NewObject creates new GtldMilter instance */
func RunServer(socket net.Listener) {
	// declare milter init function
	init := func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
		return &GtldMilter{},
			milter.OptAddHeader | milter.OptChangeHeader,
			milter.OptNoConnect | milter.OptNoHelo | milter.OptNoBody | milter.OptNoHeaders | milter.OptNoEOH
	}
	// start server
	if err := milter.RunServer(socket, init); err != nil {
		log.Fatal(err)
	}
}

/* main program */
func main() {
	// parse commandline arguments
	var protocol, address, gtlds, dests string
	flag.StringVar(&protocol,
		"proto",
		"unix",
		"Protocol family (unix or tcp)")
	flag.StringVar(&address,
		"addr",
		"/var/spool/postfix/milter/gtld",
		"Bind to address or unix domain socket")
	flag.StringVar(&gtlds,
		"gtlds",
		"/etc/postfix/gtlds.bad",
		"Filename of suspicious GTLDs, one per line")
	flag.StringVar(&dests,
		"dests",
		"/etc/postfix/dests.bad",
		"Filename of destination emails to filter suspicious GTLDs on, one per line")
	flag.Parse()

	// read in bad GTLDs
	SuspiciousGTLD = make(map[string]bool, 100)
	f, err := os.Open(gtlds)
	if err != nil {
		log.Fatal("could not open gtlds file: ", err)
	}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		SuspiciousGTLD[line] = true
	}
	//log.Println("gtlds: ", SuspiciousGTLD)

	// read in bad dests
	SuspiciousDests = make(map[string]bool, 10)
	f, err = os.Open(dests)
	if err != nil {
		log.Fatal("could not open dests file: ", err)
	}
	scan = bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		SuspiciousDests[line] = true
	}
	//log.Println("gtlds: ", SuspiciousDests)

	// make sure the specified protocol is either unix or tcp
	if protocol != "unix" && protocol != "tcp" {
		log.Fatal("invalid protocol name")
	}

	// make sure socket does not exist
	if protocol == "unix" {
		// ignore os.Remove errors
		os.Remove(address)
	}

	// bind to listening address
	socket, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	if protocol == "unix" {
		// set mode 0660 for unix domain sockets
		if err := os.Chmod(address, 0660); err != nil {
			log.Fatal(err)
		}
		// remove socket on exit
		defer os.Remove(address)
	}

	log.Println("starting GTLD milter")

	// run server
	go RunServer(socket)

	// sleep forever
	select {}
}
