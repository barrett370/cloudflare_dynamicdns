package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/barrett370/cloudflare_dynamicdns/cloudflare"
	"github.com/barrett370/cloudflare_dynamicdns/cron"
	"github.com/barrett370/cloudflare_dynamicdns/net"
	"github.com/barrett370/cloudflare_dynamicdns/workflow"
)

type cfServicer interface {
	GetDomainRecord(zoneID, recordName string) (cloudflare.Record, error)
	UpdateDNSRecordIP(record cloudflare.Record, newIP string) error
	ListZones() ([]cloudflare.Zone, error)
}

var (
	CFAPIToken        string
	CFZoneName        string
	CFIntervalSeconds int

	cfZoneID string
)

var (
	cfService cfServicer
)

func updateDNSRecord(logger *log.Logger) (err error) {
	var (
		record    cloudflare.Record
		currentIP net.IPConfig
	)
	record, err = cfService.GetDomainRecord(cfZoneID, CFZoneName)
	if err != nil {
		return err
	}
	currentIP, err = net.CurrentIP()
	if err != nil {
		return err
	}
	if record.Content != currentIP.IPAddress {
		logger.Printf("DNS record target does not match current IP, updating... \n record_target: %s, current_ip: %s\n", record.Content, currentIP.IPAddress)
		return cfService.UpdateDNSRecordIP(record, currentIP.IPAddress)
	} else {
		logger.Println("IPs match, nothing to update")
	}
	return
}

func parseEnvironment() error {
	CFAPIToken = os.Getenv("CF_API_TOKEN")
	CFZoneName = os.Getenv("CF_ZONE_NAME")
	cfIntervalSecondsStr := os.Getenv("CF_INTERVAL_SECONDS")
	var err error
	CFIntervalSeconds, err = strconv.Atoi(cfIntervalSecondsStr)
	return err
}

func main() {
	err := parseEnvironment()
	if err != nil {
		log.Fatalf("error loading environment, %v", err)
	}
	syncWorkflow, err := workflow.NewCloudflareSyncWorkflow(CFAPIToken, CFZoneName)
	if err != nil {
		log.Fatal(err)
	}
	updateCron := cron.New(fmt.Sprintf("Cloudflare DynamicDNS Service [%s:%s]", CFZoneName, cfZoneID), syncWorkflow, time.Second*time.Duration(CFIntervalSeconds))
	updateCron.Start()
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt, syscall.SIGTERM)
	<-interruptC
	log.Println("received os interrupt or kill, stopping update processes..")
	updateCron.Stop()
}
