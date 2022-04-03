module typhoon-cli

go 1.14

require (
	github.com/brianvoe/gofakeit/v6 v6.5.0
	github.com/coreybutler/go-timer v1.0.2 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/fatih/color v1.13.0
	github.com/gobuffalo/packr/v2 v2.5.1
	github.com/goccy/go-yaml v1.8.10 // indirect
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nsqio/go-nsq v1.0.8 // indirect
	github.com/osamingo/checkdigit v1.0.0 // indirect
	github.com/urfave/cli v1.22.5 // indirect
	github.com/urfave/cli/v2 v2.3.0
	github.com/vbauerster/mpb/v7 v7.0.3 // indirect
	github.com/vortex14/gofetcher v0.0.0-00010101000000-000000000000 // indirect
	github.com/vortex14/gotyphoon v0.0.0-20210907025807-7918b40fe308
	go.mongodb.org/mongo-driver v1.7.2
	go.opentelemetry.io/otel v0.20.0 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.25 // indirect
)

replace github.com/vortex14/gotyphoon => ../gotyphoon

replace github.com/vortex14/gofetcher => ../fetcher
