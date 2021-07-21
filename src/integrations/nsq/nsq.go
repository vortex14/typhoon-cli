package nsq

import (
	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"
	"strings"
	"typhoon-cli/src/typhoon/config"
)

type Service struct {
	Config *config.Config
}

func (s *Service) TestConnect() bool  {
	status := false
	projectConfig := s.Config
	NsqlookupdIP := strings.ReplaceAll(projectConfig.NsqlookupdIP, "http://","")
	client := nsq.Client{Address: NsqlookupdIP}
	err := client.Ping()
	if err != nil {
		color.Red("%s", err)
	} else {
		status = true
	}

	return status
}

