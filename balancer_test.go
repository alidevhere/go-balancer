package goBalancer

import (
	"testing"
	"time"
)

type Job struct {
	Name string
}

func (j Job) Run() {
	for i := 0; i < 5; i++ {
		println(j.Name)
		time.Sleep(1 * time.Second)
	}
}

func TestNewGoBalancer(t *testing.T) {

	loadBalancer := NewGoBalancer(Config[Job]{
		MaxJobsToExecuteAtSameTime: 1,
	})
	loadBalancer.AddJob(Job{Name: "test1"})
	loadBalancer.AddJob(Job{Name: "test2"})

	loadBalancer.WaitForRunningJobs()

}
