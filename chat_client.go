package glance

import (
	"github.com/k0kubun/pp"
	"log"
)

type ChatServerClient interface {
	Say(string, string)
}

func (g *Glancer) PrepareChatServers() {
	pp.Println(g.Setting.Config.ChatServers)

	g.ChatClientMap = make(map[string]ChatServerClient)
	for _, item := range g.Setting.Config.ChatServers {
		switch item.Type {
		case "irc":
			g.ChatClientMap[item.Name] = g.ConnectToChatServerIRC(&item)
		case "slack":
			g.ChatClientMap[item.Name] = g.ConnectToChatServerSlack(&item)
		default:
			log.Fatal("cannot guess chat server type")
		}
	}
}
