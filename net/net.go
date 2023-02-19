package net

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type forwardedAddresses []string

type IFConfig struct {
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
	s := string(bs)
	addrs := strings.Split(s, ",")
	var cleanAddrs []string
	for _, addr := range addrs {
		cleanAddrs = append(cleanAddrs, strings.TrimSpace(addr))
	}
	*c = cleanAddrs
	return
}

func CurrentIFConfig() (IFConfig, error) {
	resp, err := http.Get(addrIFConfig)
	if err != nil {
		return IFConfig{}, err
	}
	defer resp.Body.Close()
	var ifConfig IFConfig
	err = json.NewDecoder(resp.Body).Decode(&ifConfig)
	return ifConfig, err
}

var noTraceCurrentIPRe = regexp.MustCompile(`ip=([\d\.]+)`)

func noTraceCurrentIP() (string, error) {
	resp, err := http.Get("https://1.1.1.1/cdn-cgi/trace")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// body := string(bodyBytes)
	ipLine := string(noTraceCurrentIPRe.Find(bodyBytes))
	ipLineParts := strings.Split(ipLine, "=")
	if len(ipLineParts) < 2 {
		return "", errors.New("malformed current ip response")
	}
	return strings.TrimSpace(ipLineParts[1]), nil
}

func CurrentIP(trace bool) (string, error) {
	if trace {
		ifConfig, err := CurrentIFConfig()
		if err != nil {
			return "", err
		}
		return ifConfig.IPAddress, nil
	}
	return noTraceCurrentIP()
}
