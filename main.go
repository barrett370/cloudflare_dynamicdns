package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/barrett370/cloudflare_dynamicdns/config"
	"github.com/barrett370/cloudflare_dynamicdns/workflow"
	"github.com/barrett370/crongo"
)

var (
	CFAPIToken        string
	CFZoneName        string
	CFIntervalSeconds int
	CFConfigPath      string
)

func parseEnvironment() error {
	CFConfigPath = os.Getenv("CF_CONFIG_PATH")
	if CFConfigPath != "" {
		return config.Load(CFConfigPath)
	}
	CFAPIToken = os.Getenv("CF_API_TOKEN")
	CFZoneName = os.Getenv("CF_ZONE_NAME")
	cfIntervalSecondsStr := os.Getenv("CF_INTERVAL_SECONDS")
	var err error
	CFIntervalSeconds, err = strconv.Atoi(cfIntervalSecondsStr)
	config.Config.Cloudflare = []config.CloudflareConfig{{ZoneName: CFZoneName, IntervalSeconds: int64(CFIntervalSeconds), Authentication: config.CloudflareAuthentication{APIToken: CFAPIToken}}}
	return err
}

func startUpdateCrons() (crons []*crongo.Scheduler) {
	for _, zone := range config.Config.Cloudflare {
		syncWorkflow, err := workflow.NewCloudflareSyncWorkflow(zone.Authentication.APIToken, zone.ZoneName)
		if err != nil {
			log.Fatal(err)
		}
		updateCron := crongo.New(fmt.Sprintf("Cloudflare DynamicDNS Service [%s]", zone.ZoneName), syncWorkflow, time.Second*time.Duration(zone.IntervalSeconds))
		updateCron.Start()
		crons = append(crons, updateCron)
	}
	return
}

func stopUpdateCrons(crons []*crongo.Scheduler) {
	var wg sync.WaitGroup
	wg.Add(len(crons))
	for _, c := range crons {
		go func(c *crongo.Scheduler) {
			defer wg.Done()
			c.Stop()
		}(c)
	}
	wg.Wait()
}

func main() {
	err := parseEnvironment()
	if err != nil {
		log.Fatalf("error loading environment, %v", err)
	}
	crons := startUpdateCrons()
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt, syscall.SIGTERM)
	<-interruptC
	log.Println("received os interrupt or kill, stopping update processes..")
	stopUpdateCrons(crons)
}
