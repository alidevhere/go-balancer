package main

import (
	"encoding/csv"
	"fmt"
	"os"

	b "github.com/alidevhere/go_balancer"
)

type FileJob struct {
	FileName string
}

func (j FileJob) Run() {
	f, err := os.Open(j.FileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}

		fmt.Printf("Record: %+v\n", record)
	}
	fmt.Println("Finished processing file: ", j.FileName)
	fmt.Println()
}
func main() {

	loadBalancer := b.NewGoBalancer(b.Config[FileJob]{
		MaxJobsToExecuteAtSameTime: 1,
	})
	loadBalancer.AddJob(FileJob{FileName: "test1.csv"})
	loadBalancer.AddJob(FileJob{FileName: "test2.csv"})

	loadBalancer.WaitForRunningJobs()
}
