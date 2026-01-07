package Scheduler

import (
	"fmt"
	"project/Technical_Service/Pattern"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	cron           *cron.Cron
	patternChecker *Pattern.PatternChecker
}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		cron:           cron.New(),
		patternChecker: Pattern.NewPatternChecker(),
	}
}

func (s *SchedulerService) Start() {
	fmt.Println("⏳ Starting Scheduler Service...")

	// Run every minute
	_, err := s.cron.AddFunc("* * * * *", func() {
		fmt.Println("⏰ Scheduler Tick: Checking patterns...")
		s.patternChecker.CheckPatterns()
	})

	if err != nil {
		fmt.Println("❌ Error adding cron job:", err)
		return
	}

	s.cron.Start()
	fmt.Println("✅ Scheduler Started (Running every minute)")
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
}
