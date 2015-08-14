package glance

import (
	"fmt"
	_ "github.com/k0kubun/pp"
	"github.com/thoj/go-ircevent"
	"log"
)

type ChatServerClientIRC struct {
	ChatServerClient
	Cli *irc.Connection
}

func (c ChatServerClientIRC) Say(ch, msg string) {
	fmt.Printf("%s : %s\n", ch, msg)
	c.Cli.Join(ch)
	c.Cli.Privmsg(ch, msg)
}

func (g *Glancer) ConnectToChatServerIRC(srv *ChatServer) ChatServerClient {
	host := srv.Host
	port := srv.Port
	pass := srv.Password

	cli := irc.IRC("gl", "glance")
	cli.Password = pass
	if err := cli.Connect(host + ":" + port); err != nil {
		log.Fatal(err)
	}

	return ChatServerClientIRC{Cli: cli}
}
