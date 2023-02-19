package cloudflare

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	client httpClient
}

func New(apiToken string) *Service {
	return &Service{
		client: httpClient{
			c:        http.DefaultClient,
			host:     "https://api.cloudflare.com/",
			apiToken: apiToken,
		},
	}
}

func (s *Service) GetDomainRecord(zoneID, recordName string) (Record, error) {
	resp, err := s.client.Get(fmt.Sprintf("/client/v4/zones/%s/dns_records", zoneID))
	if err != nil {
		return Record{}, nil
	}
	defer resp.Body.Close()

	var jsonResp ListDNSRecordsResponse
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return Record{}, nil
	}
	if len(jsonResp.Result) == 0 {
		return Record{}, fmt.Errorf("no records found in zoneID %s with name %s, %+v", zoneID, recordName, jsonResp)
	}
	for _, record := range jsonResp.Result {
		if record.ZoneName == recordName {
			return record, nil
		}
	}

	return Record{}, fmt.Errorf("no records found in zoneID %s with name %s", zoneID, recordName)
}

func (s *Service) UpdateDNSRecordIP(record Record, newIP string) error {

	updateRequest := UpdateDNSRecordRequest{
		Type:    "A",
		Name:    record.Name,
		Content: newIP,
		TTL:     record.TTL,
		Proxied: record.Proxied,
	}

	resp, err := s.client.Put(fmt.Sprintf("/client/v4/zones/%s/dns_records/%s", record.ZoneID, record.ID), updateRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var urr UpdateDNSRecordResponse
	err = json.NewDecoder(resp.Body).Decode(&urr)
	if err != nil {
		return err
	}
	if !urr.Success {
		return fmt.Errorf("error updating DNS record, response: %+v", urr)
	}
	return nil
}

func (s *Service) ListZones() ([]Zone, error) {
	resp, err := s.client.Get("/client/v4/zones")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var lzr ListZonesResponse
	err = json.NewDecoder(resp.Body).Decode(&lzr)

	return lzr.Result, err
}
