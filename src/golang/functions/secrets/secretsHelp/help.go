// Sets the package name to import from the helper runner
package secretsHelp

// Imports necessary packages for the function to print text into the console
import (
	"fmt"
)

// Declare a function to show help for the secret commandlet
func ShowHelp() {
	fmt.Println("This is the secret scanner help! Here's the options available for you under 'scan':")
	fmt.Println("    - 'scan [profile] will run a secret matching scan with the options specifcied for that profile on config.json.")
	fmt.Println("")
}
