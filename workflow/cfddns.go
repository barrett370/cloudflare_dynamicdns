package workflow

import (
	"errors"
	"fmt"
	"log"

	"github.com/barrett370/cloudflare_dynamicdns/cloudflare"
	"github.com/barrett370/cloudflare_dynamicdns/net"
)

type cfServicer interface {
	GetDomainRecord(zoneID, recordName string) (cloudflare.Record, error)
	UpdateDNSRecordIP(record cloudflare.Record, newIP string) error
	ListZones() ([]cloudflare.Zone, error)
}

type CloudflareSyncWorkflow struct {
	cfService        cfServicer
	zoneID, zoneName string
}

func NewCloudflareSyncWorkflow(cfAPIToken, cfZoneName string) (*CloudflareSyncWorkflow, error) {
	cfService := cloudflare.New(cfAPIToken)
	zones, err := cfService.ListZones()
	if err != nil {
		return nil, fmt.Errorf("error getting zones, %v", err)
	}
	if len(zones) == 0 {
		return nil, errors.New("no zones found")
	}
	var cfZoneID string
	for _, zone := range zones {
		if zone.Name == cfZoneName {
			cfZoneID = zone.ID
			break
		}
	}
	if cfZoneID == "" {
		return nil, fmt.Errorf("could not find zone matching zone name: %s", cfZoneName)
	}
	return &CloudflareSyncWorkflow{
		cfService: cfService,
		zoneID:    cfZoneID,
		zoneName:  cfZoneName,
	}, nil
}

func (w *CloudflareSyncWorkflow) Run(logger *log.Logger) (err error) {
	var (
		record    cloudflare.Record
		currentIP string
	)
	record, err = w.cfService.GetDomainRecord(w.zoneID, w.zoneName)
	if err != nil {
		return err
	}
	currentIP, err = net.CurrentIP(false)
	if err != nil {
		return err
	}
	if record.Content != currentIP {
		logger.Printf("DNS record target does not match current IP, updating... \n record_target: %s, current_ip: %s\n", record.Content, currentIP)
		return w.cfService.UpdateDNSRecordIP(record, currentIP)
	} else {
		logger.Println("IPs match, nothing to update")
	}
	return
}
