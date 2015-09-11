package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/acidlemon/rds-throwlog/mysqlslow"
	"github.com/acidlemon/rds-throwlog/restrds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/internal/protocol/rest"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	var dbid = flag.String("database", "", "database identifier")
	var path = flag.String("path", "", "log file path")
	flag.Parse()

	if *dbid == "" || *path == "" {
		fmt.Println("Usage: rds-throwlog --database=[database identifier] --path=[log file path]")
		return
	}

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
		log.Println(err)
	}
	defer out.Body.Close()
	//result, err := ioutil.ReadAll(out.Body)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	log.Println("download completed")

	records := mysqlslow.Parse(out.Body)
	for _, r := range records {
		data, err := json.Marshal(r)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	}

	/*
		dir := filepath.Dir(*path)
		if dir != "." {
			os.MkdirAll(dir, 0755)
		}

		err = ioutil.WriteFile(*path, result, 0644)
		if err != nil {
			log.Println(err)
		}
	*/
}
