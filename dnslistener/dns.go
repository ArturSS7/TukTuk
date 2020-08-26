package dnslistener

import (
	"TukTuk/config"
	"TukTuk/database"
	"TukTuk/emailalert"
	"TukTuk/telegrambot"
	"errors"
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
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

var domain string
var records map[string]string

func StartDNS(Domain string) {
	records = make(map[string]string)
	domain = Domain
	records["*."+domain] = config.Settings.DomainConfig.NonExistingIPV4
	records["*."+domain+"6"] = config.Settings.DomainConfig.NonExistingIPV6
	records["existing."+domain] = config.Settings.DomainConfig.IPV4
	records["existing."+domain+"6"] = config.Settings.DomainConfig.IPV6
	records["acme."+domain] = config.Settings.DomainConfig.AcmeTxtChallenge
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

func HandlerTCP(w dns.ResponseWriter, req *dns.Msg) {
	Handler(w, req)
}

func HandlerUDP(w dns.ResponseWriter, req *dns.Msg) {
	Handler(w, req)
}

func logDNS(query string, sourceIp string) {
	var lastInsertId int64 = 0
	err := database.DNSDB.QueryRow("insert into dns (data, source_ip, time) values ($1, $2, $3)  RETURNING id", html.EscapeString(query), sourceIp, time.Now().String()).Scan(&lastInsertId)
	if err != nil {
		log.Println(err)
	}

	//Send Alert to telegram
	telegrambot.BotSendAlert(html.EscapeString(query), sourceIp, time.Now().String(), "DNS", lastInsertId)
	//Send Alert to email
	emailalert.SendEmailAlert("DNS Alert", "Remoute Address: "+sourceIp+"\n+"+html.EscapeString(query)+"\n"+time.Now().String())

}

func Handler(w dns.ResponseWriter, req *dns.Msg) {
	defer w.Close()
	fmt.Println(req)
	question := req.Question[0]
	matched, err := regexp.MatchString(`^*.`+domain, question.Name)
	fmt.Println(matched)
	if err != nil {
		log.Println(err)
	}
	m := new(dns.Msg)
	m.SetReply(req)
	m.Compress = false
	switch req.Opcode {
	case dns.OpcodeQuery:
		if question.Qtype == dns.TypeCAA {
			if question.Name == domain {
				answerAcmeCAA(m)
				log.Println("answered acme caa query")
				w.WriteMsg(m)
			}
		}
	}
	if matched {
		switch req.Opcode {
		case dns.OpcodeQuery:
			if question.Qtype == dns.TypeTXT {
				if question.Name == "_acme-challenge."+domain {
					answerAcme(m)
					log.Println("answered acme query")
				}
			}
		}
		var result bool
		re := regexp.MustCompile(`([a-z0-9\-]+\.)` + domain)
		d := re.Find([]byte(strings.ToLower(question.Name)))
		fmt.Println(d)
		rows, err := database.DNSDB.Query("select exists(select domain from dns_domains where (domain = $1 and delete_time > $2) or (domain=$1 and delete_time is null))", d, time.Now().Unix())
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			err = rows.Scan(&result)
			if err != nil {
				log.Println(err)
			}
		}
		if result {
			logDNS(req.String(), w.RemoteAddr().String())
			answerQuery(m, true)
		} else {
			answerQuery(m, false)
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

func answerAcme(m *dns.Msg) {
	for _, q := range m.Question {
		ip := records["acme."+domain]
		if ip != "" {
			rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, ip))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func answerAcmeCAA(m *dns.Msg) {
	for _, q := range m.Question {
		rr, err := dns.NewRR(fmt.Sprintf("%s CAA %s", q.Name, "0 issuewild \"letsencrypt.org\""))
		if err == nil {
			m.Answer = append(m.Answer, rr)
		}
	}
}

func answerQuery(m *dns.Msg, resolveIP bool) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)

			ip := ""
			if resolveIP {
				ip = records["existing."+domain]
			} else {
				ip = records["*."+domain]
			}

			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeAAAA:
			log.Printf("ipv6 query for %s\n", q.Name)

			ip := ""
			if resolveIP {
				ip = records["existing."+domain+"6"]
			} else {
				ip = records["*."+domain]
			}

			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
				if err != nil {
					log.Println(err)
				}
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}

		case dns.TypeMX:
			log.Printf("mx query for %s\n", q.Name)

			ip := ""
			if resolveIP {
				ip = records["existing."+domain+"6"]
			} else {
				ip = records["*."+domain]
			}

			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s MX 10 %s", q.Name, "pwn.bar"))

				if err != nil {
					log.Println(err)
				}
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
