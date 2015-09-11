package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/acidlemon/rds-throwlog/mysqlslow"
	"github.com/acidlemon/rds-throwlog/restrds"
)

func main() {
	var dbid = flag.String("database", "", "database identifier")
	var path = flag.String("path", "", "log file path")
	flag.Parse()

	if *dbid == "" || *path == "" {
		fmt.Println("Usage: rds-throwlog --database=[database identifier] --path=[log file path]")
		return
	}

	stream, err := restrds.Fetch(dbid, path)
	if err != nil {
		log.Println(err)
		return
	}
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
