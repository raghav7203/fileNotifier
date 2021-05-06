package fileNotifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
)

type Config struct {
	Port     string
	Endpoint string
}

func Add(filePath string) {

	var old_config Config
	var new_config Config

	body, _ := ioutil.ReadFile(filePath)

	json.Unmarshal(body, &old_config)
	json.Unmarshal(body, &new_config)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					// log.Println("modified file:", event.Name)
					body1, _ := ioutil.ReadFile(filePath)
					json.Unmarshal(body1, &new_config)
					if old_config.Port != new_config.Port {
						fmt.Println(old_config.Port, " changed ", new_config.Port)
					}
					if old_config.Endpoint != new_config.Endpoint {
						fmt.Println(old_config.Endpoint, " changed ", new_config.Endpoint)
					}
					old_config = new_config
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
