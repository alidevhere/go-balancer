package main

func main() {

	loadBalancer := NewGoBalancer(Config[Job]{
		MaxJobsToExecuteAtSameTime: 1,
	})
	loadBalancer.AddJob(Job{Name: "test1"})
	loadBalancer.AddJob(Job{Name: "test2"})

	loadBalancer.WaitForRunningJobs()
}
