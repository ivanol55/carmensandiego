// Sets the package name for the main script
package scannerRunner

// Imports necessary packages for the main logic loop to run the necessary helpers and tools based on script arguments
import (
	"carmensandiego/src/golang/functions/helpers/configManagement"
	"carmensandiego/src/golang/functions/helpers/databaseManagement"
	"carmensandiego/src/golang/functions/helpers/errorManagement"
	"carmensandiego/src/golang/functions/helpers/threadManagement"
	"fmt"
	"regexp"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

// Prepare everything necessary to run a secret scan
func RunScan(profileName string, initializedDatabases []*badger.DB) {
	// Prepare the databases we've set up so we can use them
	var patternsDatabase *badger.DB
	var filesDatabase *badger.DB
	var resultsDatabase *badger.DB
	patternsDatabase = initializedDatabases[0]
	filesDatabase = initializedDatabases[1]
	resultsDatabase = initializedDatabases[2]
	// Get regex rules based on the profile filtering
	var regexDictionary []databaseManagement.KeyValueEntry
	fmt.Println("Fetching the desired regex rules...")
	regexDictionary = getChosenRegexRules(profileName, patternsDatabase)
	fmt.Println("Scanning the files for secrets...")
	scanDatabaseFiles(profileName, regexDictionary, filesDatabase, resultsDatabase)
	showResults(resultsDatabase)
}

func getChosenRegexRules(profileName string, patternsDatabase *badger.DB) []databaseManagement.KeyValueEntry {
	var profileRulesets []string
	var profile configManagement.Profile
	profile = configManagement.GetProfile(profileName)
	profileRulesets = profile.Ruleset
	var regexDBEntries []databaseManagement.KeyValueEntry
	if profileRulesets[0] == "all" && len(profileRulesets) == 1 {
		regexDBEntries = databaseManagement.FetchAllKeysAndValues(patternsDatabase)
	} else {
		var regexRule string
		for _, regexRule = range profileRulesets {
			var fetchedRegexDBEntries []databaseManagement.KeyValueEntry
			var fetchedRegexDBEntry databaseManagement.KeyValueEntry
			fetchedRegexDBEntries = databaseManagement.FetchPrefixedKeysAndValues(patternsDatabase, regexRule)
			for _, fetchedRegexDBEntry = range fetchedRegexDBEntries {
				regexDBEntries = append(regexDBEntries, fetchedRegexDBEntry)
			}
		}
	}
	return regexDBEntries
}

func scanDatabaseFiles(profileName string, regexDictionary []databaseManagement.KeyValueEntry, filesDatabase *badger.DB, resultsDatabase *badger.DB) {
	// Prepare as many queues as we have routine threads for the profile
	var queues [][]string
	var profile configManagement.Profile
	var threads int
	profile = configManagement.GetProfile(profileName)
	threads = profile.Threads
	queues = threadManagement.GenerateQueues(profileName)
	// Separate all targets into the queues
	var fileList []string
	fileList = databaseManagement.FetchAllKeys(filesDatabase)
	queues = threadManagement.PopulateQueues(profileName, fileList, queues)
	// Run the scan as a parallel task per queue
	var secretScanningWaitGroup sync.WaitGroup
	secretScanningWaitGroup.Add(threads)
	var thread int
	for thread = 1; thread <= threads; thread = thread + 1 {
		go scanFilesForSecrets(&secretScanningWaitGroup, queues[thread-1], filesDatabase, regexDictionary, resultsDatabase)
	}
	secretScanningWaitGroup.Wait()
}

// Helper multithreaded function for scanDatabaseFiles to run a scan for each queue
func scanFilesForSecrets(secretScanningWaitGroup *sync.WaitGroup, queue []string, filesDatabase *badger.DB, regexDictionary []databaseManagement.KeyValueEntry, resultsDatabase *badger.DB) {
	defer secretScanningWaitGroup.Done()
	var filenameToScan string
	var fileContents []byte
	var regexEntry databaseManagement.KeyValueEntry
	var regexpObject *regexp.Regexp
	var matchedStrings [][]byte
	var matchedStringsEntry []byte
	var resultingScanField string
	var err error
	for _, filenameToScan = range queue {
		fileContents = databaseManagement.ReadEntry(filesDatabase, filenameToScan)
		for _, regexEntry = range regexDictionary {
			regexpObject, err = regexp.Compile(regexEntry.Value)
			errorManagement.CheckError(err)
			matchedStrings = regexpObject.FindAll(fileContents, -1)
			if matchedStrings != nil {
				resultingScanField = "Matched the file \"" + filenameToScan + "\" with rule \"" + regexEntry.Key + "\" for the following strings:\n"
				for _, matchedStringsEntry = range matchedStrings {
					resultingScanField = resultingScanField + "\t" + string(matchedStringsEntry) + "\n"
				}
				databaseManagement.WriteEntry(resultsDatabase, filenameToScan+regexEntry.Key, resultingScanField)
			}
		}
	}
}

func showResults(resultsDatabase *badger.DB) {
	var scanResults []databaseManagement.KeyValueEntry
	var scanResult databaseManagement.KeyValueEntry
	scanResults = databaseManagement.FetchAllKeysAndValues(resultsDatabase)
	if len(scanResults) == 0 {
		fmt.Println("No secrets found")
	} else {
		for _, scanResult = range scanResults {
			fmt.Println(scanResult.Value)
		}
	}
}
