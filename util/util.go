package util

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vote_bot/log"
	"github.com/vote_bot/models"
	"gopkg.in/yaml.v2"
)

var DirsMap map[string]string

func GetDirs() {
	rootDir := GetRootPath()
	DirsMap = make(map[string]string)
	DirsMap["modulesDir"] = rootDir + "/modules"
}

// URL call
// ex)
// CallUrl(url) : 2sec (default)
// CallUrl(url,  10) : 10sec time out
func CallUrl(url string, sec ...int) (string, error) {
	timeout := 2 * time.Second

	if len(sec) > 0 {
		timeout = time.Duration(sec[0]) * time.Second
	}
	client := resty.New()
	client.SetHeader("Accept", "application/json")
	client.SetTimeout(timeout) // timeout check
	client.SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	})
	req, err := client.R().EnableTrace().Get(url)
	if err != nil {
		return "", err
	}

	return string(req.Body()), nil
}

// Import Yaml Files
// ex)
// yamlName : local.env , prod.env
func LoadYaml(yamlName string) {
	yamlFile, _ := os.ReadFile(yamlName)
	err := yaml.Unmarshal(yamlFile, &models.Config)
	//err := godotenv.Load(yamlName)

	if err != nil {
		log.Logger.Error.Fatal("Error loading "+yamlName, err)
	}

	// log.Logger.Trace.Println(envName)
}

// Init
// log, configuration file, db setting
// ex) Init() Enter numbers according to the location of the current folder. Default values are not required
// Run as Init(1) if you enter one folder
func Init() {
	rootPath := GetRootPath()
	os.Setenv("ROOT_PATH", rootPath) // Required by logger

	// logger Init
	log.LogInit()

	// yaml
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		LoadYaml(rootPath + "/config/local.yaml")
		log.Logger.Trace.Println("local.yaml")
	} else {
		LoadYaml(rootPath + "/config/prod.yaml")
		log.Logger.Trace.Println("prod.yaml")
	}

	// log.Logger.Trace.Println("runtime.GOOS", runtime.GOOS)

	// get Dirs
	GetDirs()
	// DB
	err := models.ConnectDatabase()
	if err != nil {
		log.Logger.Error.Fatal("Error loading Database : ", err)
	}

}

// Required to get root folder absolute path
func GetRootPath() string {
	// Get the absolute path of the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		// Check if a go.mod file exists in the current directory.
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Move up one directory.
		parentDir := filepath.Dir(dir)

		// If we've reached the root directory ("/"), exit the loop.
		if parentDir == dir {
			break
		}

		// Continue searching in the parent directory.
		dir = parentDir
	}

	return ""
}
