package dnslistener

import (
	"TukTuk/database"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"html"
	"log"
	"regexp"
	"sync"
	"time"
)

type DnsMsg struct {
	Timestamp       string
	SourceIP        string
	DestinationIP   string
	DnsQuery        string
	DnsAnswer       []string
	DnsAnswerTTL    []string
	NumberOfAnswers string
	DnsResponseCode string
	DnsOpCode       string
}

func StartDNS() {
	startServer()
}

func startServer() {
	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", HandlerTCP)

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", HandlerUDP)

	tcpServer := &dns.Server{Addr: "0.0.0.0:53",
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	udpServer := &dns.Server{Addr: "0.0.0.0:53",
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			log.Fatal("TCP-server start failed", err.Error())
		}
	}()
	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			log.Fatal("UDP-server start failed", err.Error())
		}
	}()
}

var records = map[string]string{
	"*.tt.pwn.bar.":  "127.0.0.1",
	"*.tt.pwn.bar.6": "[::]:1",
}

func HandlerTCP(w dns.ResponseWriter, req *dns.Msg) {
	Handler(w, req)
}

func HandlerUDP(w dns.ResponseWriter, req *dns.Msg) {
	Handler(w, req)
}

func logDNS(query string, sourceIp string) {
	_, err := database.DNSDB.Query("insert into dns (data, source_ip, time) values ($1, $2, $3)", html.EscapeString(query), sourceIp, time.Now().String())
	if err != nil {
		log.Println(err)
	}
}

func Handler(w dns.ResponseWriter, req *dns.Msg) {
	defer w.Close()
	question := req.Question[0]
	matched, err := regexp.MatchString(`^*.tt.pwn.bar.`, question.Name)
	if err != nil {
		log.Println(err)
	}
	if matched {
		m := new(dns.Msg)
		m.SetReply(req)
		m.Compress = false
		switch req.Opcode {
		case dns.OpcodeQuery:
			logDNS(req.String(), w.RemoteAddr().String())
			answerQuery(m)
		}
		w.WriteMsg(m)
	} else {
		resp, err := Lookup(req)
		if err != nil {
			resp = &dns.Msg{}
			resp.SetRcode(req, dns.RcodeServerFailure)
			log.Println("fail", question.Name)
		}
		w.WriteMsg(resp)
	}
}

func answerQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ip := records["*.tt.pwn.bar."]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeAAAA:
			log.Printf("ipv6 query for %s\n", q.Name)
			ip := records["*.tt.pwn.bar.6"]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func Lookup(req *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	qName := req.Question[0].Name

	res := make(chan *dns.Msg, 1)
	var wg sync.WaitGroup
	L := func(nameserver string) {
		defer wg.Done()
		r, _, err := c.Exchange(req, nameserver)
		if err != nil {
			log.Printf("%s socket error on %s", qName, nameserver)
			log.Printf("error:%s", err.Error())
			return
		}
		if r != nil && r.Rcode != dns.RcodeSuccess {
			if r.Rcode == dns.RcodeServerFailure {
				return
			}
		}
		select {
		case res <- r:
		default:
		}
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Start lookup on each nameserver top-down, in every second
	nameservers := []string{"8.8.8.8:53", "8.8.4.4:53"}
	for _, nameserver := range nameservers {
		wg.Add(1)
		go L(nameserver)
		// but exit early, if we have an answer
		select {
		case r := <-res:
			return r, nil
		case <-ticker.C:
			continue
		}
	}

	// wait for all the namservers to finish
	wg.Wait()
	select {
	case r := <-res:
		return r, nil
	default:
		return nil, errors.New("can't resolve ip for" + qName)
	}
}
