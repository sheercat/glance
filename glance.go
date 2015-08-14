package glance

import (
	"bufio"
	"fmt"
	"github.com/k0kubun/pp"
	"io"
	"log"
	"os/exec"
	"regexp"
	"time"
)

type Chans struct {
	Server GlanceServer
	Msg    chan string
	State  chan error
}

type Glancer struct {
	Setting       GlanceServerList
	ChatClientMap map[string]ChatServerClient
	Chans         []Chans
}

func (g *Glancer) Start() error {
	pp.Println(g.Setting)

	for _, item := range g.Setting.Servers {
		for _, host := range item.Hosts {
			statechan := make(chan error)
			msgchan := make(chan string)
			// gorutine で tail 実行
			go execTail(host, item, msgchan, statechan)
			g.Chans = append(g.Chans, Chans{Server: item, Msg: msgchan, State: statechan})
		}
	}

	pp.Println(g.Chans)

	g.PrepareChatServers()
	g.EventLoop()

	return nil
}

func (g *Glancer) EventLoop() {
	for {
		for _, item := range g.Chans {
			select {
			case msg := <-item.Msg:
				g.Say(&item.Server, msg)
			default:
			}
		}
	}
}

func (g *Glancer) Say(gsvr *GlanceServer, msg string) {

	var filterChannels []ChatChannel
	for _, filter := range gsvr.ServerFilters {
		if filter.CompiledRegexp == nil {
			compiled, err := regexp.Compile(filter.Regexp)
			if err != nil {
				log.Println("compile regexp error")
				continue
			}
			filter.CompiledRegexp = compiled
		}

		if filter.CompiledRegexp.MatchString(msg) {
			for _, item := range filter.ChatChannels {
				filterChannels = append(filterChannels, item)
			}
		}
	}

	var channels []ChatChannel
	if len(filterChannels) > 0 {
		channels = filterChannels
	} else {
		channels = gsvr.DefaultChatChannels
	}
	for _, item := range channels {
		cli := g.ChatClientMap[item.Name]
		cli.Say(item.Ch, msg)
		time.Sleep(time.Second)
	}
}

func execTail(host string, server GlanceServer, msgchan chan string, statechan chan error) {

	// ssh
	cmd := exec.Command("ssh", "-tt", host, fmt.Sprintf(`tail -s 2 -F %s | grep -P "%s"`, server.LogPattern, server.TailFilter))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	log.Println("tail:" + host)
	printToChan(stdout, msgchan, host)

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return
	}
}

func printToChan(r io.Reader, msgchan chan string, host string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		str := scanner.Text()
		msgchan <- host + ":" + str 
	}
}

func Init(settingFile string) (Glancer, error) {
	srvs, err := LoadSetting(settingFile)
	if err != nil {
		log.Printf("%#v", err)
		return Glancer{}, err
	}
	return Glancer{Setting: srvs}, nil
}
