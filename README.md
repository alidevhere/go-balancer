# Go Balancer
Go Balancer is package to load balance any type task. For example you have an API that receives files, process them and do something with it. If you many users upload file to your system it may over burden your system.

To deal with this scenario you can use Go Balancer. It provides out of box functionality to load balance any task.
 
 You just need to implement interface 

```

type runnable interface {
	// Run() is the function that will be executed by the load balancer.
	Run()
}

```

Just define a function Run() on your struct and write the logic you want to execute for each task.

## Use

```
go get github.com/alidevhere/go_balancer
```

## Examples:
Examples are in examples directory of this repo.
