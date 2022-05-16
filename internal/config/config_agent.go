package config

import (
	"flag"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/service/agent"
)

func NewAgentConfig() *agent.Config {
	cfg := agent.Config{
		MaxIdleConns:        agent.DefaultMaxIdleConns,
		MaxIdleConnsPerHost: agent.DefaultMaxIdleConnsPerHost,
		RetryCount:          agent.DefaultRetryCount,
		RetryWaitTime:       agent.DefaultRetryWaitTime,
		RetryMaxWaitTime:    agent.DefaultRetryMaxWaitTime,
		PollMetricsBuffSize: agent.DefaultPollMetricsBuffSize,
		PostWorkersPoolSize: agent.DefaultPostWorkersPoolSize,
	}

	flag.DurationVar(&cfg.ReportInterval, "r", agent.DefaultReportInterval, "REPORT_INTERVAL")
	flag.DurationVar(&cfg.PollInterval, "p", agent.DefaultPollInterval, "POLL_INTERVAL")
	flag.StringVar(&cfg.Key, "k", "", "KEY")

	return &cfg
}
