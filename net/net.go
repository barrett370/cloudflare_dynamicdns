package net

import (
	"encoding/json"
	"net/http"
	"strings"
)

type forwardedAddresses []string

type IPConfig struct {
	IPAddress  string             `json:"ip_addr,omitempty"`
	RemoteHost string             `json:"remote_host,omitempty"`
	UserAgent  string             `json:"user_agent,omitempty"`
	Port       int64              `json:"port"`
	Method     string             `json:"method,omitempty"`
	MimeType   string             `json:"mime,omitempty"`
	Via        string             `json:"via,omitempty"`
	Forwarded  forwardedAddresses `json:"forwarded,omitempty"`
}

const addrIFConfig = "https://ifconfig.me/all.json"

var _ json.Unmarshaler = &forwardedAddresses{}

func (c *forwardedAddresses) UnmarshalJSON(bs []byte) (err error) {
	var s string
	err = json.Unmarshal(bs, &s)
	if err != nil {
		return
	}
	addrs := strings.Split(s, ",")
	var cleanAddrs []string
	for _, addr := range addrs {
		cleanAddrs = append(cleanAddrs, strings.TrimSpace(addr))
	}
	*c = cleanAddrs
	return
}

func CurrentIP() (IPConfig, error) {
	resp, err := http.Get(addrIFConfig)
	if err != nil {
		return IPConfig{}, err
	}
	defer resp.Body.Close()
	var ifConfig IPConfig
	err = json.NewDecoder(resp.Body).Decode(&ifConfig)
	return ifConfig, err
}
