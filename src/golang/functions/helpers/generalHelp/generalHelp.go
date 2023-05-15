// Sets the package name to import from the helper runner
package generalHelp

// Imports necessary packages for the function to print text into the terminal
import (
	"fmt"
)

// Declares a function that prints the general system help
func ShowHelp() {
	fmt.Println("This tool offers several scanning options depending on your secret scanning needs. Here's what's available as first-level options (more features in development!):")
	fmt.Println("    - 'secrets' will let you scan and generate reports on your specified folder structure on your profile. run this script with just 'scan' to get more information")
	fmt.Println("")
}
