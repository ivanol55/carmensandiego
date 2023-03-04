// Sets the package name for the main script
package setupScan

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/secrets/databaseManagement"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dgraph-io/badger/v3"
)

func SetupScan() []*badger.DB {
	var patternsDatabase *badger.DB
	var filesDatabase *badger.DB
	var resultsDatabase *badger.DB
	patternsDatabase = databaseManagement.InitializeDatabase()
	filesDatabase = databaseManagement.InitializeDatabase()
	resultsDatabase = databaseManagement.InitializeDatabase()
	populatePatternsDatabase(patternsDatabase)
	populateFilesDatabase(filesDatabase)
	var InitializedDatabases []*badger.DB
	InitializedDatabases = append(InitializedDatabases, patternsDatabase)
	InitializedDatabases = append(InitializedDatabases, filesDatabase)
	InitializedDatabases = append(InitializedDatabases, resultsDatabase)
	return InitializedDatabases
}

func populatePatternsDatabase(patternsDatabase *badger.DB) {
	type PatternData struct {
		Name       string `json:"name"`
		Regex      string `json:"regex"`
		Confidence string `json:"confidence"`
	}
	type Pattern struct {
		Pattern PatternData `json:"pattern"`
	}
	type Patterns struct {
		Patterns []Pattern `json:"patterns"`
	}
	var patternJsonContents []byte
	patternJsonContents, _ = ioutil.ReadFile("src/patterns/regex-patterns.json")
	var patternsStruct Patterns
	// Unmarshal and write the configContents json file contents into the configStruct struct variable
	json.Unmarshal(patternJsonContents, &patternsStruct)
	// Read the seleced profile from the object and return it to the caller
	var patternsArray Patterns
	patternsArray = patternsStruct
	var pattern Pattern
	for _, pattern = range patternsArray.Patterns {
		fmt.Println(pattern.Pattern.Regex)
	}
}

func populateFilesDatabase(filesDatabase *badger.DB) {

}
