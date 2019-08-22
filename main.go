package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

type config struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Root      string `json:"root"`
	LocalPath string `json:"localPath"`
}

func loadConfiguration(file string) config {
	var config config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func getFiles(directory string, client *gowebdav.Client) []os.FileInfo {
	files, err := client.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
	}
	return files
}

func walkWebdavDir(directory string, currentLocalPath string, client *gowebdav.Client) {
	os.MkdirAll(currentLocalPath, os.ModePerm)

	for _, file := range getFiles(directory, client) {
		webDavSubDirectory := filepath.Join(directory, file.Name())
		localPath := filepath.Join(currentLocalPath, file.Name())
		if file.IsDir() {
			walkWebdavDir(webDavSubDirectory, localPath, client)
		} else {
			bytes, _ := client.Read(webDavSubDirectory)
			ioutil.WriteFile(localPath, bytes, 0644)
			fmt.Println(localPath)
		}
	}
}

func main() {
	config := loadConfiguration("./config.json")

	client := gowebdav.NewClient(config.Root, config.Username, config.Password)

	walkWebdavDir("product_images/", config.LocalPath, client)
}
