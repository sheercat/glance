package glance

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

type ChatChannel struct {
	Name string `json:"name"`
	Ch   string `json:"ch"`
}

type ServerFilter struct {
	Regexp       string        `json:"regexp"`
	ChatChannels []ChatChannel `json:"chat_channels"`
	CompiledRegexp *regexp.Regexp
}

type GlanceServer struct {
	Hosts               []string       `json:"hosts"`
	TailFilter          string         `json:"tail_filter"`
	LogPattern          string         `json:"log_pattern"`
	DefaultChatChannels []ChatChannel  `json:"default_chat_channels"`
	ServerFilters       []ServerFilter `json:"server_filters"`
}

type ChatServer struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Password   string `json:"password"`
	WebHookURL string `json:"webhook_url"`
}

type Config struct {
	ChatServers []ChatServer `json:"chat_servers"`
}

type GlanceServerList struct {
	Servers []GlanceServer `json:"server_list"`
	Config  Config         `json:"config"`
}

func LoadSetting(settingFile string) (GlanceServerList, error) {
	var d GlanceServerList
	jstring, err := ioutil.ReadFile(settingFile)
	if err != nil {
		log.Println(err)
		return d, err
	}
	err = json.Unmarshal(jstring, &d)
	if err != nil {
		log.Println(err)
		return d, err
	}
	log.Printf("%+v\n", d)
	return d, nil
}
