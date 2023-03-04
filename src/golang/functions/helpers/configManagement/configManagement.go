// Sets the package name to import from the helper runner
package configManagement

// Imports necessary packages for the function to process data from disk, in this case reading documentation, and encoding json data into golang structs
import (
	"encoding/json"
	"io/ioutil"
)

// Declares the base struct of the configuration file, consisting of preferences and profiles
type Config struct {
	Profiles map[string]Profile `json:"profiles"`
}

// Declares the struct for profiles, consisting of A name, a description and an array of strings that are considered targets
type Profile struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Confidence  []string `json:"confidence"`
	Ruleset     []string `json:"ruleset"`
}

// Declare the function that retrieves the profile from the configuration file and returns it back as a struct
func GetProfile(profileName string) Profile {
	// Read the configuration json file from the current folder
	var configContents []byte
	configContents, _ = ioutil.ReadFile("config.json")
	// Prepare a base variable with the Config structure to store all data we need
	var configStruct Config
	// Unmarshal and write the configContents json file contents into the configStruct struct variable
	json.Unmarshal(configContents, &configStruct)
	// Read the seleced profile from the object and return it to the caller
	var profileObject Profile
	profileObject = configStruct.Profiles[profileName]
	return profileObject
}
