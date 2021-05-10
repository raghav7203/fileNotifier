package fileNotifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
)

func parseMap(aMap1 map[string]interface{}, aMap2 map[string]interface{}) {
	for key, val := range aMap2 {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			// fmt.Println(key)
			parseMap(aMap1[key].(map[string]interface{}), aMap2[key].(map[string]interface{}))
		case []interface{}:
			// fmt.Println(key)
			parseArray(aMap1[key].([]interface{}), aMap2[key].([]interface{}))
		default:
			// fmt.Println("Index", i, ":", concreteVal)
			// fmt.Println(key, ":", aMap1[key], ":", concreteVal)
			if aMap1[key] != concreteVal {
				fmt.Println(key, " changed ", aMap1[key], " to ", concreteVal)
			}
		}
	}
}

func parseArray(anArray1 []interface{}, anArray2 []interface{}) {
	for i, val := range anArray2 {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			// fmt.Println("Index:", i)
			parseMap(anArray1[i].(map[string]interface{}), anArray2[i].(map[string]interface{}))
		case []interface{}:
			// fmt.Println("Index:", i)
			parseArray(anArray1[i].([]interface{}), anArray2[i].([]interface{}))
		default:
			// fmt.Println("Index", i, ":", concreteVal)
			if anArray1[i] != concreteVal {
				fmt.Println("Index ", i, " changed ", anArray1[i], " to ", concreteVal)
			}
		}
	}
}

func reconfigMap(old_config_map map[string]interface{}, new_config_map map[string]interface{}) {
	for key, val := range new_config_map {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			// fmt.Println(key)
			reconfigMap(old_config_map[key].(map[string]interface{}), new_config_map[key].(map[string]interface{}))
		case []interface{}:
			// fmt.Println(key)
			reconfigArray(old_config_map[key].([]interface{}), new_config_map[key].([]interface{}))
		default:
			// fmt.Println("Index", i, ":", concreteVal)
			// fmt.Println(key, ":", aMap1[key], ":", concreteVal)
			old_config_map[key] = concreteVal
		}
	}
}

func reconfigArray(old_config_array []interface{}, new_config_array []interface{}) {
	for key, val := range new_config_array {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			// fmt.Println(key)
			reconfigMap(old_config_array[key].(map[string]interface{}), new_config_array[key].(map[string]interface{}))
		case []interface{}:
			// fmt.Println(key)
			reconfigArray(old_config_array[key].([]interface{}), new_config_array[key].([]interface{}))
		default:
			// fmt.Println("Index", i, ":", concreteVal)
			// fmt.Println(key, ":", aMap1[key], ":", concreteVal)
			old_config_array[key] = concreteVal
		}
	}
}

func Add(filePath string) {

	old_config := make(map[string]interface{})
	new_config := make(map[string]interface{})

	body, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(body, &old_config)
	json.Unmarshal(body, &new_config)

	// parseMap(old_config, new_config)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

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
					parseMap(old_config, new_config)
					fmt.Println("-----------------------")
					reconfigMap(old_config, new_config)
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
	select {}
}
