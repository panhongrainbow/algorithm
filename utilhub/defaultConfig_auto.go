package utilhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// AutoConfig ⛏️ is a type constraint for struct types that store default configuration values used in automated testing. (自动预设配置)
type AutoConfig any

// AutoConfigType ⛏️ is a struct for Bptree automatic unit test configuration.
type AutoConfigType struct {
	Mechanism string   `json:"mechanism" default:"auto"` // 🧪 When the mechanism selection is set to `auto`, all tests will be conducted.
	Record    struct { // 🧪 Record contains configurations related to test record storage.
		TestRecordPath   string `json:"testRecordPath" default:"/temp/test_record"` // 🧪 TestRecordPath specifies the directory path where test records will be saved.
		ManualRecordDate string `json:"manualRecordDate" default:"0000-00-00"`      // 🧪 Manual testing requires specifying the previous test date.
		ManualRecordFile string `json:"manualRecordFile" default:"empty"`           // 🧪 Manual testing requires specifying the previous test record file.
		IsInsideProject  bool   `json:"isInsideProject" default:"true"`             // 🧪 IsInsideProject indicates whether the test records are stored inside the project directory.
	} `json:"record"`
	Parameters struct { // Parameters contains configurations for test execution parameters.
		RandomTotalCount             int64 `json:"randomTotalCount" default:"7500000"`        // 🧪 RandomTotalCount represents the number of elements to be generated for random testing.
		RandomMin                    int64 `json:"randomMin" default:"10"`                    // 🧪 RandomMin represents the minimum value for generating random numbers.
		RandomHitCollisionPercentage int64 `json:"randomHitCollisionPercentage" default:"70"` // 🧪 Random number hit collision percentage.
		// Calculate the maximum random value.
		// randomTotalCount/randomHitCollisionPercentage*100 + randomMin = randomMax
		// 7500000 / 70 * 100 + 10 = 10714295
		RandomMax int64 `json:"randomMax" default:"10714295"` // 🧪 RandomMax represents the maximum value for generating random numbers.
		BpWidth   []int `json:"bpWidth" default:"3,4,5,6,7"`
	} `json:"parameters"`
	PoolStage struct { // This is primarily used to test boundary conditions.
		MinRemovals       int64 `json:"minRemovals" default:"5"`        // 🧪 Lower bound of items to remove in this stage.
		MaxRemovals       int64 `json:"maxRemovals" default:"50"`       // 🧪 Upper bound of items to remove in this stage.
		MinPreserveInPool int64 `json:"minPreserveInPool" default:"10"` // 🧪 Lower bound of items to remain in the pool after this stage.
		MaxPreserveInPool int64 `json:"maxPreserveInPool" default:"20"` // 🧪 Upper bound of items to remain in the pool after this stage.
	} `json:"poolStage"`
	CyclicStress struct { // metal fatigue style endurance test.
		CyclicStressCount int `json:"cyclicStressCount" default:"10"` // 🧪 Number of fatigue test cycles.
	} `json:"cyclicStress"`
	/*
		ManualTest struct { // 使用手动测试，重现之前的错误
			EnableBulkInsertDelete   bool `json:"enableBulkInsertDelete" default:"false"`
			EnableRandomizedBoundary bool `json:"enableRandomizedBoundary" default:"false"`
			EnableNodeEnduranceTest  bool `json:"enableNodeEnduranceTest" default:"false"`
		} `json:"manualTest"`
	*/
}

// ParseAuto ⛏️ loads the default configuration for automated testing from struct tags and applies it to the provided struct.
func ParseAuto(autoConfig *TestProcessConfigType) error {
	// Prepare the variable outside the closure function.
	var cfgFile AutoConfig = new(AutoConfigType)
	var err error
	var projectPath, file string

	// Use Golang's sync.Once to prevent the setting from being overwritten.
	_onesAutoConfig.Do(func() {

		// Get the default configuration directory.
		projectPath, err = GetProjectDir(filepath.Join(ProjectName))
		if err != nil {
			return
		}

		// Get the struct name to use as the filename.
		file, err = GetStructName(&cfgFile)
		if err != nil {
			return
		}

		// Return the result of _parseAuto.
		err = _parseAuto(filepath.Join(projectPath, "config", file+".json"), cfgFile)
		if err != nil {
			return
		}

		// If the record is configured to be inside the project directory,
		// prepend the project path to the test record path
		if cfgFile.(*AutoConfigType).Record.IsInsideProject == true {
			cfgFile.(*AutoConfigType).Record.TestRecordPath = filepath.Join(projectPath, cfgFile.(*AutoConfigType).Record.TestRecordPath)
		}

		// When the mechanism is set to automatic, the test subdirectory will automatically be set to the current date.
		if cfgFile.(*AutoConfigType).Mechanism == "auto" {
			cfgFile.(*AutoConfigType).Record.ManualRecordDate = TestTimeString("2006-01-02", "Asia/Shanghai")
		}

		// HACK => For a B plus tree, you can specify the branching factor here and test each case one by one.
		cfgFile.(*AutoConfigType).Parameters.BpWidth = []int{3, 10, 13, 14, 17} // 指定分叉因子

		// Copy all Parameters from cfgFile to autoConfig.
		autoConfig.Mechanism = cfgFile.(*AutoConfigType).Mechanism
		autoConfig.Record = cfgFile.(*AutoConfigType).Record
		autoConfig.Parameters = cfgFile.(*AutoConfigType).Parameters
		autoConfig.PoolStage = cfgFile.(*AutoConfigType).PoolStage
		autoConfig.CyclicStress = cfgFile.(*AutoConfigType).CyclicStress
	})

	// Return nil to indicate the operation completed successfully.
	return err
}

// _parseAuto ⛏️ loads the default configuration from struct tags and applies it to the provided struct.
// Configuration from the file, if the file exists, and applies and overwrites the struct. (以文件的配置为主,结构体配置为次)
func _parseAuto(filePath string, cfg AutoConfig) error {
	// Check if the config is a pointer to a struct.
	if reflect.ValueOf(cfg).Kind() != reflect.Pointer {
		return errors.New("config must be a pointer to a struct")
	}

	// Read the default configuration from the file.
	var content []byte
	var err error
	content, err = os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Do nothing; the tag will be handled later by applyDefaults. (之后由 applyDefaults 取 tag 决定)
		}
		return err
	}

	// Unmarshal the JSON data into the provided config and overwrite the default values.
	if err = json.Unmarshal(content, cfg); err != nil {
		return err
	}

	// applyDefaults applies the default values from struct tags to the provided config. (主要填入预设值逻辑)
	if err = applyDefaults(cfg); err != nil {
		return err
	}

	// No error occurred, return nil.
	return nil
}

// GetAutoConfig parses and returns the auto-configuration, panicking if parsing fails.
func GetAutoConfig() TestProcessConfigType {
	// Parse the auto-configuration into the global config instance.
	if err := ParseAuto(&_autoConfig); err != nil {
		fmt.Printf("autoParseErr after init: %v\n", err)
		panic(err)
	}

	// Return the initialized configuration.
	return _autoConfig
}
