package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	c "github.com/TreyBastian/colourize"
)

type (
	config struct {
		PingCount int
	}
	result struct {
		IP  string
		Ret string
	}
)

var (
	conf config
	wg   sync.WaitGroup
	r    *regexp.Regexp
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Please give the IP file as an argument")
		os.Exit(1)
	}
	conf.PingCount = 2
}

func main() {
	pattern := `---\n([\s|\S]+)$`
	r, _ = regexp.Compile(pattern)
	targets := getIPList()
	res := make(chan result, len(targets))
	wg.Add(len(targets))
	for _, ip := range targets {
		go processIP(ip, res, &wg, r)
	}
	wg.Wait()
	close(res)
	for elem := range res {
		var color = c.Green
		if !strings.Contains(elem.Ret, fmt.Sprintf("%d packets received", conf.PingCount)) {
			color = c.Red
		}
		fmt.Printf("%s\n\t%s", c.Colourize(elem.IP, c.Bold, c.Blue), c.Colourize(elem.Ret, color))
	}
}

func processIP(IP string, res chan<- result, wg *sync.WaitGroup, pattern *regexp.Regexp) {
	cmd := fmt.Sprintf("ping %s -c %d -s 1 -W 2", IP, conf.PingCount)
	if isIPv6(IP) {
		cmd = fmt.Sprintf("ping6 %s -c %d -s 1", IP, conf.PingCount)
	}
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	match := r.FindStringSubmatch(string(out))
	ret := result{
		IP: IP,
	}
	if len(match) > 1 {
		ret.Ret = match[1]
	}
	if err != nil {

		if err.Error() == "exit status 68" {
			ret.Ret = "Host name unknown"
		}
	}
	res <- ret
	wg.Done()
}

func getIPList() (targets []string) {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !contains(targets, line) {
			targets = append(targets, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func contains(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func isIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}
