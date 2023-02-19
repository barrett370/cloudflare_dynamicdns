package cron

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/barrett370/cloudflare_dynamicdns/workflow"
)

type Service struct {
	logger *log.Logger
	op     workflow.Workflower
	ticker *time.Ticker
	done   chan struct{}
}

func New(name string, wf workflow.Workflower, interval time.Duration) *Service {
	return &Service{
		logger: log.New(os.Stdout, fmt.Sprintf("[CRON: %s] ", name), log.Ldate|log.Ltime|log.Lshortfile|log.LUTC),
		op:     wf,
		ticker: time.NewTicker(interval),
		done:   make(chan struct{}),
	}
}

func (s *Service) Start() {
	s.logger.Println("starting")
	go s.loop()
}

func (s *Service) Stop() {
	s.logger.Println("stopping")
	s.done <- struct{}{}
	<-s.done
}

func (s *Service) loop() {
	for {
		select {
		case <-s.ticker.C:
			err := s.op.Run(s.logger)
			if err != nil {
				s.logger.Printf("error while running work func. err: %v\n", err)
			}
		case <-s.done:
			s.logger.Println("received stop signal, stopping")
			close(s.done)
			return
		}
	}
}
