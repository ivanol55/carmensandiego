// Sets the package name for the main script
package errorManagement

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"fmt"
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Panicked on error. Here's the error trace:")
		log.Fatal(err)
		os.Exit(1)
	}
}
