package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/acidlemon/rds-throwlog/mysqlslow"
	"github.com/acidlemon/rds-throwlog/restrds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/internal/protocol/rest"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/fluent/fluent-logger-golang/fluent"
)

func main() {
	var dbid = flag.String("database", "", "database identifier")
	var path = flag.String("path", "", "log file path")
	var fluentHost = flag.String("fluent-host", "", "fluentd hostname")
	var fluentPort = flag.Int("fluent-port", 24224, "fluentd forward port (default 24224)")
	var fluentTag = flag.String("fluent-tag", "mysql.slowquery", "fluentd tag")
	var raw = flag.Bool("raw", false, "output raw data (stdout only)")
	flag.Parse()

	if *dbid == "" || *path == "" {
		fmt.Println("Usage: rds-throwlog --database=[database identifier] --path=[log file path]")
		return
	}

	stream, err := Fetch(dbid, path)
	if err != nil {
		log.Println(err)
		return
	}
	defer stream.Close()
	log.Println("download completed")

	// output raw
	if *raw {
		io.Copy(os.Stdout, stream)
		return
	}

	// prepare fluent-logger
	var logger *fluent.Fluent
	if *fluentHost != "" {
		logger, err = fluent.New(fluent.Config{
			FluentPort: *fluentPort,
			FluentHost: *fluentHost,
		})
		if err != nil {
			log.Println("fluent.New returned error:", err)
			return
		}
	}

	records := mysqlslow.Parse(stream)
	for _, r := range records {
		if logger != nil {
			t, msg := r.ToFluentLog()
			log.Println(t)
			logger.PostWithTime(*fluentTag, t, msg)
		} else {
			data, err := json.Marshal(r)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))
		}
	}
}

func Fetch(dbid, path *string) (io.ReadCloser, error) {
	svc := rds.New(&aws.Config{
		Region: aws.String("ap-northeast-1"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	// RDSのHandlerはQuery APIになっているのでをREST APIに変更
	svc.Handlers.Build.Clear()
	svc.Handlers.Build.PushBack(rest.Build)
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.Unmarshal.PushBack(rest.Unmarshal)

	out, err := restrds.DownloadCompleteDBLogFile(svc, &restrds.DownloadCompleteDBLogFileInput{
		DBInstanceIdentifier: dbid,
		LogFileName:          path,
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
