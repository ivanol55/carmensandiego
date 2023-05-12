// Sets the package name to import from the helper runner
package threadManagement

import (
	"carmensandiego/src/golang/functions/helpers/configManagement"
	"runtime"
)

// Imports necessary packages for the function to prepare threads as necessary, needed in multiple parts of the code to enable multithreading

func GenerateQueues() [][]string {
	var profile configManagement.Profile
	var threads int
	threads = profile.Threads
	if threads < 2 {
		threads = 2
	}
	runtime.GOMAXPROCS(threads)
	// Generate queues for the nomber of goroutine threads we're running, and return the queue set
	var queueSet [][]string
	var queue []string
	var thread int
	for thread = 1; thread <= threads; thread = thread + 1 {
		queue = []string{}
		queueSet = append(queueSet, queue)
	}
	return queueSet
}

func PopulateQueues(fileList []string, queues [][]string) [][]string {
	// For each goroutine thread, populate these queues with an equal amount of files to scan
	var profile configManagement.Profile
	var threads int
	threads = profile.Threads
	var thread int
	thread = 0
	threads = threads - 1
	var file string
	for _, file = range fileList {
		queues[thread] = append(queues[thread], file)
		if thread == threads {
			thread = 0
		} else {
			thread = thread + 1
		}
	}
	return queues
}
