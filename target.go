package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
)

type Target struct {
	Ip        string
	Port      string
	HostNames []string
}

func (target *Target) url() string {
	return fmt.Sprintf("%s:%s", target.Ip, target.Port)
}

func ParserTarget(text string) Target {

	parts := strings.Split(text, ":")
	ip := parts[0]
	port := "443"
	if len(parts) == 2 {
		port = parts[1]
	}

	return Target{
		Ip:   ip,
		Port: port,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getCertHost(ch chan Target, target Target) {
	conf := &tls.Config{InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", target.url(), conf)
	if err != nil {
		log.Fatalln("Error in Dial", err)
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
			return
		}
	}()
	certs := conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		if len(cert.DNSNames) > 0 {
			target.HostNames = cert.DNSNames
			ch <- target
		}
	}
}
