// Sets the package name for the main script
package scannerSetup

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/helpers/configManagement"
	"carmensandiego/src/golang/functions/helpers/errorManagement"
	"carmensandiego/src/golang/functions/helpers/threadManagement"
	"carmensandiego/src/golang/functions/secrets/databaseManagement"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

// Prepare everything necessary to run a secret scan
func SetupScan(profileName string) []*badger.DB {
	// We create three separate databases necessary to store the regex patterns, all the target files and their contents, and a database to store our results on, so we can be as asynchronous as possible.
	var patternsDatabase *badger.DB
	var filesDatabase *badger.DB
	var resultsDatabase *badger.DB
	fmt.Println("Initializing databases...")
	patternsDatabase = databaseManagement.InitializeDatabase()
	filesDatabase = databaseManagement.InitializeDatabase()
	resultsDatabase = databaseManagement.InitializeDatabase()
	fmt.Println("Populating databases...")
	// Populate our databases with the necessary data
	populatePatternsDatabase(patternsDatabase)
	populateFilesDatabase(filesDatabase, profileName)
	fmt.Println("Clearing exemptions from the files database...")
	// Clear exemptions we set on our profile from
	clearFileExemptions(filesDatabase, profileName)
	var InitializedDatabases []*badger.DB
	// Return the databases to our original caller
	InitializedDatabases = append(InitializedDatabases, patternsDatabase)
	InitializedDatabases = append(InitializedDatabases, filesDatabase)
	InitializedDatabases = append(InitializedDatabases, resultsDatabase)
	return InitializedDatabases
}

func populatePatternsDatabase(patternsDatabase *badger.DB) {
	// Prepare and load the patterns for regex matching that live inside of the src folder
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
	// Write the patterns with their name and value into the patterns database
	var patternsArray Patterns
	patternsArray = patternsStruct
	var pattern Pattern
	var err error
	for _, pattern = range patternsArray.Patterns {
		err = databaseManagement.WriteEntry(patternsDatabase, pattern.Pattern.Name, pattern.Pattern.Regex)
		errorManagement.CheckError(err)
	}
}

func populateFilesDatabase(filesDatabase *badger.DB, profileName string) {
	var profile configManagement.Profile
	var threads int
	profile = configManagement.GetProfile(profileName)
	threads = profile.Threads
	// Read the list of all the files we would like to scan
	var err error = nil
	var baseDir string
	baseDir = "target/"
	var fileList []string
	err = getFileList(baseDir, &fileList)
	errorManagement.CheckError(err)
	// Prepare as many queues as we have routine threads for the profile
	var queues [][]string
	queues = threadManagement.GenerateQueues(profileName)
	// Separate all files into the queues
	queues = threadManagement.PopulateQueues(profileName, fileList, queues)
	// Read file contents into the database as threads to speed up the process
	var filePopulationWaitGroup sync.WaitGroup
	filePopulationWaitGroup.Add(threads)
	var thread int
	for thread = 1; thread <= threads; thread = thread + 1 {
		go readFilesToDatabase(&filePopulationWaitGroup, queues[thread-1], filesDatabase)
	}
	filePopulationWaitGroup.Wait()
}

func getFileList(dir string, fileList *[]string) error {
	// Get the list of files we want to scan from the system and store them in fileList
	var entries []fs.DirEntry
	var entry fs.DirEntry
	var err error = nil
	var fullPath string
	// Recursively list the files and folders of a source directory. If it finds a folder, iterate over it until it has no folder structure left.
	entries, err = os.ReadDir(dir)
	errorManagement.CheckError(err)
	for _, entry = range entries {
		fullPath = filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			err = getFileList(fullPath, fileList)
			errorManagement.CheckError(err)
		} else {
			*fileList = append(*fileList, fullPath)
		}
	}
	return nil
}

func readFilesToDatabase(filePopulationWaitGroup *sync.WaitGroup, fileQueue []string, filesDatabase *badger.DB) {
	// For each goroutine that runs, read all of the files that are available in the specified lists, so the key for the database is the filename, and the value is the file's contents.
	defer filePopulationWaitGroup.Done()
	var file string
	var fileContents []byte
	var fileContentsString string
	var err error = nil
	for _, file = range fileQueue {
		fileContents, err = ioutil.ReadFile(file)
		fileContentsString = string(fileContents)
		file = strings.TrimPrefix(file, "target/")
		databaseManagement.WriteEntry(filesDatabase, file, fileContentsString)
		errorManagement.CheckError(err)
	}
}

func clearFileExemptions(filesDatabase *badger.DB, profileName string) {
	// Get the exemptions for the profiles, i.e. the files we can ignore as we have validated their secret contents, like unit tests
	var profile configManagement.Profile
	profile = configManagement.GetProfile(profileName)
	var exemptions []string
	exemptions = profile.Exemptions
	var exemption string
	// Drop each exemption from the files database
	for _, exemption = range exemptions {
		databaseManagement.DeleteEntry(filesDatabase, exemption)
	}
}
