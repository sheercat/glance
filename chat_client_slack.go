package glance

import (
	"encoding/json"
	"github.com/k0kubun/pp"
	"log"
	"net/http"
	"net/url"
)

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type attachment struct {
	Fallback string   `json:"fallback"`
	Pretext  string   `json:"pretext"`
	Color    string   `json:"color"`
	Fields   []*field `json:"fields"`
}

type payload struct {
	Attachments []*attachment `json:"attachments"`
}

type textPayload struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

type ChatServerClientSlack struct {
	ChatServerClient
	WebHookURL string
}

func (c ChatServerClientSlack) Say(ch string, msg string) {
	pp.Println(ch, msg)
	title := "alert"
	// 引用でおくるとき
	_, err := json.Marshal(&payload{Attachments: []*attachment{
		&attachment{
			Fallback: title,
			Pretext:  title,
			Fields: []*field{
				&field{
					Title: "",
					Value: msg,
					Short: false,
				},
			},
		},
	},
	})
	// フラットなテキスト + channel override
	p, err := json.Marshal(&textPayload{Text: msg, Channel: ch})
	res, err := http.PostForm(c.WebHookURL, url.Values{"payload": []string{string(p)}})
	if err != nil {
		log.Println("error on slack")
	}
	if res.StatusCode != 200 {
		pp.Println(res.Status)
	}
}

func (g *Glancer) ConnectToChatServerSlack(c *ChatServer) ChatServerClient {

	return ChatServerClientSlack{WebHookURL: c.WebHookURL}
}
