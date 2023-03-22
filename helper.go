// Sets the package name for the main script
package main

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/helpers/generalHelp"
	"carmensandiego/src/golang/functions/helpers/greeting"
	"carmensandiego/src/golang/functions/helpers/requiredArgCount"
	"carmensandiego/src/golang/functions/secrets/setupScan"
	"os"
)

// Function that runs when the program is started, executes the main application logic
func main() {
	// Show the application greeting with ascii art title
	greeting.ShowGreeting()
	// Check if enough args were provided. If not, the program shows an error and exits
	requiredArgCount.CheckForArgs(os.Args, 2, "general")
	// Check which tool to run depending on the first script argument
	switch os.Args[1] {
	// If "help" is provided as the first script argument, show the application general help page and exit the program
	case "help":
		generalHelp.ShowHelp()
	// If "infra" is provided as the first script argument, run the infra application logic check
	case "scan":
		// Check if the required amount of arguments were sent. If not, show an error to the user and exit the program
		requiredArgCount.CheckForArgs(os.Args, 3, "scan")
		setupScan.SetupScan(os.Args[2])
	// If the first argument doesn't match any supported arguments, show the general application help and exit the program
	default:
		generalHelp.ShowHelp()
	}
}
