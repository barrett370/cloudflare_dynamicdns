package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DomainRecords struct {
	Records []struct {
		ID        string `json:"id"`
		ZoneID    string `json:"zone_id"`
		ZoneName  string `json:"zone_name"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Proxiable bool   `json:"proxiable"`
		Proxied   bool   `json:"proxied"`
		TTL       int    `json:"ttl"`
		Locked    bool   `json:"locked"`
		Meta      struct {
			AutoAdded           bool   `json:"auto_added"`
			ManagedByApps       bool   `json:"managed_by_apps"`
			ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
			Source              string `json:"source"`
		} `json:"meta"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
	} `json:"result"`
	Success    bool          `json:"success"`
	Errors     []interface{} `json:"errors"`
	Messages   []interface{} `json:"messages"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"result_info"`
}

type DNSData struct {
	ZoneID    string
	EntryID   string
	currentIP string
}

func getStoredIP() (string, error) {
	dnsData, err := getDNSEntryID()
	if err != nil {
		return "", err
	}
	return dnsData.currentIP, nil

}

func getCurrentIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func updateDNSRecord(newIP string, dnsData DNSData) error {

	type Payload struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
		Proxied bool   `json:"proxied"`
	}

	data := Payload{
		Type:    "A",
		Name:    DOMAIN,
		Content: newIP,
		TTL:     1,
		Proxied: true,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", dnsData.ZoneID, dnsData.EntryID), body)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Email", "barrett370@gmail.com")
	req.Header.Set("X-Auth-Key", TOKEN)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	return nil
}

func getDNSEntryID() (DNSData, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/10a667bda412483884440f0388161c6f/dns_records?&name=%v&page=1&per_page=20&order=type&direction=desc&match=all",DOMAIN), nil)
	if err != nil {
		return DNSData{}, err
	}
	req.Header.Set("X-Auth-Email", AUTH_EMAIL)
	req.Header.Set("X-Auth-Key", TOKEN)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DNSData{}, err
	}
	defer resp.Body.Close()

	records := new(DomainRecords)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DNSData{}, err
	}
	err = json.Unmarshal(body, &records)
	if err != nil {
		return DNSData{}, err
	}
	if len(records.Records) > 0 {
		for _, record := range records.Records {
			if record.ZoneName == DOMAIN {
				log.Println("Found DNS entry")
				return DNSData{record.ZoneID, record.ID, record.Content}, nil
			}
		}
	} else {
		return DNSData{}, errors.New("no records found")
	}

	return DNSData{}, errors.New(fmt.Sprintf("could not find ID for %v",DOMAIN))
}

var (
	TOKEN string
        DOMAIN string
        AUTH_EMAIL string
)

func main() {
	dnsData, err := getDNSEntryID()
	if err != nil {
		log.Fatal(err)
	}
	storedIP := dnsData.currentIP
	currentIP, err := getCurrentIP()
	if err != nil {
		log.Fatal(err)
	}
	if storedIP != currentIP {
		log.Printf("%s, %s IPs do not match, updating DNS record\n", storedIP, currentIP)
		err = updateDNSRecord(currentIP, dnsData)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("%s, %s match, exiting\n", storedIP, currentIP)
	}
}
