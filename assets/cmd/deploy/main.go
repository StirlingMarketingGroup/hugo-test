package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cenkalti/backoff"
	"gopkg.in/yaml.v3"
)

func main() {
	cmd := exec.Command("aws", "s3", "ls")
	cmd.Stderr = nil
	cmd.Stdout = nil
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("failed to ping s3: %w", err))
	}

	start := time.Now()

	cancel, in, out := StartWorkerPool(100)
	defer cancel()

	d, err := os.Open("sites")
	if err != nil {
		panic(fmt.Errorf("failed to open the sites folder: %w", err))
	}

	sites, err := d.Readdirnames(0)
	if err != nil {
		panic(fmt.Errorf(`failed to read dir names in "sites": %w`, err))
	}

	for _, site := range sites {
		in <- &WorkRequest{
			Op: HugoBuild,

			Site:   site,
			Dir:    "sites/" + site,
			Bucket: "v3." + site,
		}
	}

	redirectsFile, err := os.Open("redirects.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to open redirects.yaml: %w", err))
	}

	var redirects map[string]string
	err = yaml.NewDecoder(redirectsFile).Decode(&redirects)
	if err != nil {
		panic(fmt.Errorf("failed to decode redirects.yaml: %w", err))
	}

	currRedirectWorkers := make(map[string]int)
	currCloudFrontWorkers := 0

Events:
	for resp := range out {
		if resp.Err != nil {
			panic(err)
		}

		switch wo := resp.Wo; wo.Op {
		case HugoBuild:
			wo.Op = S3Sync
			in <- wo
		case S3Sync:
			wo.Op = GetConfParent
			in <- wo
		case GetConfParent:
			for k, v := range redirects {
				in <- &WorkRequest{
					Op:               MakeRedirect,
					Site:             resp.Wo.Site,
					Dir:              resp.Wo.Dir,
					Bucket:           resp.Wo.Bucket,
					Key:              k,
					RedirectLocation: strings.Replace(v, "$PARENT", resp.Config.Params.Parent, 1),
					Config:           resp.Config,
				}
				currRedirectWorkers[wo.Site]++
			}
		case MakeRedirect:
			currRedirectWorkers[wo.Site]--
			if currRedirectWorkers[wo.Site] == 0 {
				wo.Op = CloudFrontInvalidate
				in <- wo
				currCloudFrontWorkers++
			}
		case CloudFrontInvalidate:
			currCloudFrontWorkers--
			if currCloudFrontWorkers == 0 {
				break Events
			}
		default:
			panic(fmt.Errorf("unhandled op %q", resp.Wo.Op))
		}
	}

	fmt.Println("finished deploying in", time.Since(start))
}

type op = uint

const (
	HugoBuild op = iota + 1
	S3Sync
	GetConfParent
	MakeRedirect
	CloudFrontInvalidate
)

type siteConfig struct {
	Params struct {
		Parent                 string
		CloudFrontDistribution string
	}
}

type WorkRequest struct {
	Op               op
	Site             string
	Dir              string
	Bucket           string
	Key              string
	RedirectLocation string
	Config           *siteConfig
}

type WorkResponse struct {
	Wo     *WorkRequest
	Config *siteConfig
	Err    error
}

func Worker(ctx context.Context, id int, in chan *WorkRequest, out chan *WorkResponse) {
	for { //run indefinitely
		select {
		case <-ctx.Done():
			return
		case wo := <-in:
			out <- workerExec(wo)
		}
	}
}

func workerExec(wo *WorkRequest) *WorkResponse {
	wr := &WorkResponse{Wo: wo}

	err := backoff.Retry(func() error {
		switch wo.Op {
		case HugoBuild:
			if err := hugoBuild(wo); err != nil {
				return backoff.Permanent(err)
			}
			return nil
		case S3Sync:
			return s3Sync(wo)
		case GetConfParent:
			config, err := getConf(wo)
			if err != nil {
				return backoff.Permanent(err)
			}
			wr.Config = config
			return nil
		case MakeRedirect:
			return makeRedirect(wo)
		case CloudFrontInvalidate:
			return cloudFrontInvalidate(wo)
		default:
			return fmt.Errorf("unhandled op %q", wo.Op)
		}
	}, backoff.NewExponentialBackOff())

	wr.Err = err
	return wr
}

func StartWorkerPool(numberWorkers int) (context.CancelFunc, chan *WorkRequest, chan *WorkResponse) {
	in := make(chan *WorkRequest, numberWorkers)
	out := make(chan *WorkResponse, numberWorkers)
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < numberWorkers; i++ {
		go Worker(ctx, i+1, in, out)
	}
	return cancel, in, out
}

func hugoBuild(wo *WorkRequest) error {
	cmd := exec.Command("hugo", "--minify")
	cmd.Dir = wo.Dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run `hugo --minify`: %w", err)
	}

	return nil
}

func s3Sync(wo *WorkRequest) error {
	cmd := exec.Command("aws", "s3", "sync", "public", "s3://"+wo.Bucket, "--delete", "--acl", "public-read")
	cmd.Dir = wo.Dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to sync site to s3: %w", err)
	}

	return nil
}

func getConf(wo *WorkRequest) (*siteConfig, error) {
	config := new(siteConfig)
	_, err := toml.DecodeFile(filepath.Join(wo.Dir, "config.toml"), config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s/config.toml: %w", wo.Site, err)
	}

	return config, nil
}

func makeRedirect(wo *WorkRequest) error {
	k := wo.Key
	if strings.HasSuffix(k, "/") {
		k += "index.html"
	}
	k = k[1:]

	log.Println(k, "->", wo.RedirectLocation)

	cmd := exec.Command("aws", "s3api", "put-object", "--website-redirect-location", wo.RedirectLocation, "--bucket", wo.Bucket, "--key", k, "--acl", "public-read", "--content-type", "text/html")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to make redirect %q -> %q: %w", k, wo.RedirectLocation, err)
	}

	return nil
}

func cloudFrontInvalidate(wo *WorkRequest) error {
	cmd := exec.Command("aws", "cloudfront", "create-invalidation", "--distribution-id", wo.Config.Params.CloudFrontDistribution, "--paths", "/*")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to invalidate cloudfront for %q: %w", wo.Site, err)
	}

	return nil
}
