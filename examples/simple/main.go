package main

import (
	"time"

	b "github.com/alidevhere/go_balancer"
)

type Job struct {
	Name string
}

func (j Job) Run() {
	for i := 0; i < 5; i++ {
		println(i, j.Name)
		time.Sleep(1 * time.Second)
	}
}
func main() {

	loadBalancer := b.NewGoBalancer(b.Config[Job]{
		MaxJobsToExecuteAtSameTime: 1,
	})
	loadBalancer.AddJob(Job{Name: "test1"})
	loadBalancer.AddJob(Job{Name: "test2"})

	loadBalancer.WaitForRunningJobs()
}
