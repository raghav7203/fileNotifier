package fileNotifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// output maps
var outputDiffJson = make(map[string]interface{})
var outputDiffYaml = make(map[interface{}]interface{})

// output strings
var outputStringJson = ""
var outputStringYaml = ""

// helper arrays
var valArrStr []string
var valArrInt []float64
var valArrStrYaml []string
var valArrIntYaml []int // yaml is not typecasting on its own unlike json thats why int

// to reinitialize map and string
func clearAllJson() {
	outputStringJson = ""
	for k := range outputDiffJson {
		delete(outputDiffJson, k)
	}
}

func checkExtJson(filePath string) bool {
	fileExt := strings.SplitAfter(filePath, ".")
	if fileExt[len(fileExt)-1] != "json" {
		fmt.Println("json config not found")
		return false
	}
	fmt.Println("json config found")
	return true
}

// adding watcher to the json config
func AddJson(filePath string, change func(m map[string]interface{}, s string)) {
	// fsnotify things
	ticker := time.NewTicker(300 * time.Millisecond)
	events := make([]fsnotify.Event, 0)

	// checking for json extension file
	if !checkExtJson(filePath) {
		return
	}

	// maps for old and new state of config file
	old_config := make(map[string]interface{})
	new_config := make(map[string]interface{})

	body, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(body, &old_config)
	// json.Unmarshal(body, &new_config)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Add current event to array for batching.
				events = append(events, event)
			case <-ticker.C:
				// Checks on set interval if there are events.
				if len(events) > 0 {
					// Display messages for write event in batch.
					// for _, event := range events {
					// 	if event.Op&fsnotify.Write == fsnotify.Write {
					// 		// fmt.Printf("\nFile write detected: %#v\n", event)
					// 	}
					// }
					// Run whatever you want here, one time for all batched events.
					// reading newly saved config
					body1, _ := ioutil.ReadFile(filePath)
					json.Unmarshal(body1, &new_config)

					fmt.Println(len(old_config), len(new_config))
					// passing old and new states for comparison
					if len(old_config) > len(new_config) {
						deleteMapJson(old_config, new_config)
					} else {
						parseMapJson(old_config, new_config)
					}

					for i := range old_config {
						delete(old_config, i)
					}
					// assigning new state to old state
					json.Unmarshal(body1, &old_config)
					// clearing the new state(fresh new state will be fetched from reading file again)
					for k := range new_config {
						delete(new_config, k)
					}
					// callback func
					if len(outputDiffJson) != 0 {
						change(outputDiffJson, outputStringJson)
						clearAllJson()
					}
					fmt.Println("------------------")

					// Empty the batch array.
					events = make([]fsnotify.Event, 0)
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

func deleteMapJson(aMap1 map[string]interface{}, aMap2 map[string]interface{}) {
	if len(aMap2) > len(aMap1) {
		parseMapJson(aMap1, aMap2)
	} else {
		for key, val := range aMap1 {
			switch concreteVal := val.(type) {
			case map[string]interface{}:
				// case where a whole new map is added in config
				if _, ok := aMap2[key].(map[string]interface{}); ok {
					parseMapJson(aMap1[key].(map[string]interface{}), aMap2[key].(map[string]interface{}))
				} else {
					parseMapJson(aMap1[key].(map[string]interface{}), nil)
				}
			case []interface{}:
				// case where a whole new array is added in config
				if _, ok := aMap2[key].([]interface{}); ok {
					parseArrayJson(aMap1[key].([]interface{}), aMap2[key].([]interface{}), key)
				} else {
					parseArrayJson(aMap1[key].([]interface{}), nil, key)
				}
			default:
				if aMap2[key] != concreteVal {
					// assigning difference to output map(if nil would be added on it's own in map)
					outputDiffJson[key] = concreteVal
					// assigning diff to output string depending it was modified or added from nil
					if _, ok := aMap2[key].(string); ok {
						oldval := aMap2[key].(string)
						newval := concreteVal.(string)
						outputStringJson += key + " changed " + oldval + " to " + newval + "\n"
					} else {
						outputStringJson += key + " changed " + concreteVal.(string) + " to nil\n"
					}
				}
				// if aMap2[key] != concreteVal {
				// 	fmt.Println(key, " changed from ", concreteVal, " to nil")
				// }
			}
		}
	}
}

func deleteArrayJson(anArray1 []interface{}, anArray2 []interface{}, key string) {
	valArrInt = nil
	valArrStr = nil
	if len(anArray2) > len(anArray1) {
		parseArrayJson(anArray1, anArray2, key)
	} else {
		for i, val := range anArray1 {
			switch concreteVal := val.(type) {
			case map[string]interface{}:
				if len(anArray2) > i {
					parseMapJson(anArray1[i].(map[string]interface{}), anArray2[i].(map[string]interface{}))
				} else {
					parseMapJson(anArray1[i].(map[string]interface{}), nil)
				}
				// if _, ok := anArray1[i].(map[string]interface{}); ok {
				// 	parseMapJson(anArray1[i].(map[string]interface{}), anArray2[i].(map[string]interface{}))
				// } else {
				// 	parseMapJson(nil, anArray2[i].(map[string]interface{}))
				// }
			case []interface{}:
				if len(anArray2) > i {
					parseArrayJson(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				} else {
					parseArrayJson(anArray1[i].([]interface{}), nil, key)
				}
				// if _, ok := anArray1[i].([]interface{}); ok {
				// 	parseArrayJson(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				// } else {
				// 	parseArrayJson(nil, anArray2[i].([]interface{}), key)
				// }
			default:
				// case where iteration is smaller than old config(implies only modification) and else has new values
				if len(anArray2) > i {
					if anArray2[i] != concreteVal {
						// fmt.Println("Index ", i, " changed ", anArray1[i], " to ", concreteVal)
						if _, ok := concreteVal.(string); ok {
							valArrStr = append(valArrStr, concreteVal.(string))
							outputDiffJson[key] = valArrStr
							oldval := fmt.Sprintf("%v", anArray1[i].(string))
							newval := fmt.Sprintf("%v", concreteVal.(string))
							outputStringJson += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						} else {
							valArrInt = append(valArrInt, concreteVal.(float64))
							outputDiffJson[key] = valArrInt
							oldval := fmt.Sprintf("%v", anArray1[i].(float64))
							newval := fmt.Sprintf("%v", concreteVal.(float64))
							outputStringJson += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						}
					}
				} else {
					if _, ok := concreteVal.(string); ok {
						valArrStr = append(valArrStr, concreteVal.(string))
						outputDiffJson[key] = valArrStr
						newval := fmt.Sprintf("%v", concreteVal.(string))
						outputStringJson += newval + " to nil\n"

					} else {
						valArrInt = append(valArrInt, concreteVal.(float64))
						outputDiffJson[key] = valArrInt
						newval := fmt.Sprintf("%v", concreteVal.(float64))
						outputStringJson += newval + " to nil\n"
					}
				}
			}
		}
	}
}

// checking type at each iteration and recursively calling map/array or default
func parseMapJson(aMap1 map[string]interface{}, aMap2 map[string]interface{}) {
	if len(aMap1) > len(aMap2) {
		deleteMapJson(aMap1, aMap2)
	} else {
		for key, val := range aMap2 {
			switch concreteVal := val.(type) {
			case map[string]interface{}:
				// case where a whole new map is added in config
				if _, ok := aMap1[key].(map[string]interface{}); ok {
					parseMapJson(aMap1[key].(map[string]interface{}), aMap2[key].(map[string]interface{}))
				} else {
					parseMapJson(nil, aMap2[key].(map[string]interface{}))
				}
			case []interface{}:
				// case where a whole new array is added in config
				if _, ok := aMap1[key].([]interface{}); ok {
					parseArrayJson(aMap1[key].([]interface{}), aMap2[key].([]interface{}), key)
				} else {
					parseArrayJson(nil, aMap2[key].([]interface{}), key)
				}
			default:
				if aMap1[key] != concreteVal {
					// assigning difference to output map(if nil would be added on it's own in map)
					outputDiffJson[key] = concreteVal
					// assigning diff to output string depending it was modified or added from nil
					if _, ok := aMap1[key].(string); ok {
						oldval := aMap1[key].(string)
						newval := concreteVal.(string)
						outputStringJson += key + " changed " + oldval + " to " + newval + "\n"
					} else {
						outputStringJson += key + " changed nil " + " to " + concreteVal.(string) + "\n"
					}
				}
			}
		}
	}
}

func parseArrayJson(anArray1 []interface{}, anArray2 []interface{}, key string) {
	valArrInt = nil
	valArrStr = nil
	if len(anArray1) > len(anArray2) {
		deleteArrayJson(anArray1, anArray2, key)
	} else {
		for i, val := range anArray2 {
			switch concreteVal := val.(type) {
			case map[string]interface{}:
				if len(anArray1) > i {
					parseMapJson(anArray1[i].(map[string]interface{}), anArray2[i].(map[string]interface{}))
				} else {
					parseMapJson(nil, anArray2[i].(map[string]interface{}))
				}
				// if _, ok := anArray1[i].(map[string]interface{}); ok {
				// 	parseMapJson(anArray1[i].(map[string]interface{}), anArray2[i].(map[string]interface{}))
				// } else {
				// 	parseMapJson(nil, anArray2[i].(map[string]interface{}))
				// }
			case []interface{}:
				if len(anArray1) > i {
					parseArrayJson(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				} else {
					parseArrayJson(nil, anArray2[i].([]interface{}), key)
				}
				// if _, ok := anArray1[i].([]interface{}); ok {
				// 	parseArrayJson(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				// } else {
				// 	parseArrayJson(nil, anArray2[i].([]interface{}), key)
				// }
			default:
				// case where iteration is smaller than old config(implies only modification) and else has new values
				if len(anArray1) > i {
					if anArray1[i] != concreteVal {
						// fmt.Println("Index ", i, " changed ", anArray1[i], " to ", concreteVal)
						if _, ok := concreteVal.(string); ok {
							valArrStr = append(valArrStr, concreteVal.(string))
							outputDiffJson[key] = valArrStr
							oldval := fmt.Sprintf("%v", anArray1[i].(string))
							newval := fmt.Sprintf("%v", concreteVal.(string))
							outputStringJson += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						} else {
							valArrInt = append(valArrInt, concreteVal.(float64))
							outputDiffJson[key] = valArrInt
							oldval := fmt.Sprintf("%v", anArray1[i].(float64))
							newval := fmt.Sprintf("%v", concreteVal.(float64))
							outputStringJson += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						}
					}
				} else {
					if _, ok := concreteVal.(string); ok {
						valArrStr = append(valArrStr, concreteVal.(string))
						outputDiffJson[key] = valArrStr
						newval := fmt.Sprintf("%v", concreteVal.(string))
						outputStringJson += "nil to " + newval + "\n"

					} else {
						valArrInt = append(valArrInt, concreteVal.(float64))
						outputDiffJson[key] = valArrInt
						newval := fmt.Sprintf("%v", concreteVal.(float64))
						outputStringJson += "nil to " + newval + "\n"
					}
				}
			}
		}
	}
}

// to reinitialize map and string
func clearAllYaml() {
	outputStringYaml = ""
	for k := range outputDiffYaml {
		delete(outputDiffYaml, k)
	}
}

func checkExtYaml(filePath string) bool {
	fileExt := strings.SplitAfter(filePath, ".")
	if fileExt[len(fileExt)-1] != "yaml" {
		fmt.Println("yaml config not found")
		return false
	}
	fmt.Println("yaml config found")
	return true
}

func AddYaml(filePath string, change func(m map[interface{}]interface{}, s string)) {
	// fsnotify things
	ticker := time.NewTicker(300 * time.Millisecond)
	events := make([]fsnotify.Event, 0)

	// checking for yaml extension file
	if !checkExtYaml(filePath) {
		return
	}

	old_config := make(map[interface{}]interface{})
	new_config := make(map[interface{}]interface{})

	body, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(body, &old_config)
	yaml.Unmarshal(body, &new_config)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Add current event to array for batching.
				events = append(events, event)
			case <-ticker.C:
				// Checks on set interval if there are events.
				if len(events) > 0 {
					// Display messages for write event in batch.
					// for _, event := range events {
					// 	if event.Op&fsnotify.Write == fsnotify.Write {
					// 		// fmt.Printf("\nFile write detected: %#v\n", event)
					// 	}
					// }
					// Run whatever you want here, one time for all batched events.
					// reading newly saved config
					body1, _ := ioutil.ReadFile(filePath)
					yaml.Unmarshal(body1, &new_config)

					fmt.Println(len(old_config), len(new_config))
					// passing old and new states for comparison
					if len(old_config) > len(new_config) {
						deleteMapYaml(old_config, new_config)
					} else {
						parseMapYaml(old_config, new_config)
					}

					for i := range old_config {
						delete(old_config, i)
					}
					// assigning new state to old state
					yaml.Unmarshal(body1, &old_config)
					// clearing the new state(fresh new state will be fetched from reading file again)
					for k := range new_config {
						delete(new_config, k)
					}
					// callback func
					if len(outputDiffYaml) != 0 {
						change(outputDiffYaml, outputStringYaml)
						clearAllYaml()
					}
					fmt.Println("------------------")

					// Empty the batch array.
					events = make([]fsnotify.Event, 0)
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

func deleteMapYaml(aMap1 map[interface{}]interface{}, aMap2 map[interface{}]interface{}) {
	if len(aMap2) > len(aMap1) {
		parseMapYaml(aMap1, aMap2)
	} else {
		for key, val := range aMap1 {
			switch concreteVal := val.(type) {
			case map[interface{}]interface{}:
				if _, ok := aMap2[key].(map[interface{}]interface{}); ok {
					parseMapYaml(aMap1[key].(map[interface{}]interface{}), aMap2[key].(map[interface{}]interface{}))
				} else {
					parseMapYaml(aMap1[key].(map[interface{}]interface{}), nil)
				}
			case []interface{}:
				if _, ok := aMap2[key].([]interface{}); ok {
					parseArrayYaml(aMap1[key].([]interface{}), aMap2[key].([]interface{}), key)
				} else {
					parseArrayYaml(aMap1[key].([]interface{}), nil, key)
				}
			default:
				if aMap2[key] != concreteVal {
					outputDiffYaml[key] = concreteVal
					if _, ok := aMap2[key].(string); ok {
						oldval := aMap2[key].(string)
						newval := concreteVal.(string)

						outputStringYaml += key.(string) + " changed " + oldval + " to " + newval + "\n"
					} else {
						outputStringYaml += key.(string) + " changed " + concreteVal.(string) + " to nil\n"
					}
				}
			}
		}
	}
}

func deleteArrayYaml(anArray1 []interface{}, anArray2 []interface{}, key interface{}) {
	valArrIntYaml = nil
	valArrStrYaml = nil
	if len(anArray2) > len(anArray1) {
		parseArrayYaml(anArray1, anArray2, key)
	} else {
		for i, val := range anArray1 {
			switch concreteVal := val.(type) {
			case map[interface{}]interface{}:
				if len(anArray2) > i {
					parseMapYaml(anArray1[i].(map[interface{}]interface{}), anArray2[i].(map[interface{}]interface{}))
				} else {
					parseMapYaml(anArray1[i].(map[interface{}]interface{}), nil)
				}
				// if _, ok := anArray1[i].(map[interface{}]interface{}); ok {
				// 	parseMapYaml(anArray1[i].(map[interface{}]interface{}), anArray2[i].(map[interface{}]interface{}))
				// } else {
				// 	parseMapYaml(nil, anArray2[i].(map[interface{}]interface{}))
				// }
			case []interface{}:
				if len(anArray2) > i {
					parseArrayYaml(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				} else {
					parseArrayYaml(anArray1[i].([]interface{}), nil, key)
				}
				// if _, ok := anArray1[i].([]interface{}); ok {
				// 	parseArrayYaml(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				// } else {
				// 	parseArrayYaml(nil, anArray2[i].([]interface{}), key)
				// }
			default:
				if len(anArray2) > i {
					if anArray2[i] != concreteVal {
						if _, ok := concreteVal.(string); ok {
							valArrStrYaml = append(valArrStrYaml, concreteVal.(string))
							outputDiffYaml[key] = valArrStrYaml
							oldval := fmt.Sprintf("%v", anArray2[i].(string))
							newval := fmt.Sprintf("%v", concreteVal.(string))
							outputStringYaml += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						} else {
							valArrIntYaml = append(valArrIntYaml, concreteVal.(int))
							outputDiffYaml[key] = valArrIntYaml
							oldval := fmt.Sprintf("%v", anArray2[i].(int))
							newval := fmt.Sprintf("%v", concreteVal.(int))
							outputStringYaml += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						}
					}
				} else {
					if _, ok := concreteVal.(string); ok {
						valArrStrYaml = append(valArrStrYaml, concreteVal.(string))
						outputDiffYaml[key] = valArrStrYaml
						newval := fmt.Sprintf("%v", concreteVal.(string))
						outputStringYaml += newval + " to nil\n"

					} else {
						valArrIntYaml = append(valArrIntYaml, concreteVal.(int))
						outputDiffYaml[key] = valArrIntYaml
						newval := fmt.Sprintf("%v", concreteVal.(int))
						outputStringYaml += newval + " to nil\n"
					}
				}

			}
		}
	}
}

func parseMapYaml(aMap1 map[interface{}]interface{}, aMap2 map[interface{}]interface{}) {
	if len(aMap1) > len(aMap2) {
		deleteMapYaml(aMap1, aMap2)
	} else {
		for key, val := range aMap2 {
			switch concreteVal := val.(type) {
			case map[interface{}]interface{}:
				if _, ok := aMap1[key].(map[interface{}]interface{}); ok {
					parseMapYaml(aMap1[key].(map[interface{}]interface{}), aMap2[key].(map[interface{}]interface{}))
				} else {
					parseMapYaml(nil, aMap2[key].(map[interface{}]interface{}))
				}
			case []interface{}:
				if _, ok := aMap1[key].([]interface{}); ok {
					parseArrayYaml(aMap1[key].([]interface{}), aMap2[key].([]interface{}), key)
				} else {
					parseArrayYaml(nil, aMap2[key].([]interface{}), key)
				}
			default:
				if aMap1[key] != concreteVal {
					outputDiffYaml[key] = concreteVal
					if _, ok := aMap1[key].(string); ok {
						oldval := aMap1[key].(string)
						newval := concreteVal.(string)

						outputStringYaml += key.(string) + " changed " + oldval + " to " + newval + "\n"
					} else {
						outputStringYaml += key.(string) + " changed nil " + " to " + concreteVal.(string) + "\n"
					}
				}
			}
		}
	}
}

func parseArrayYaml(anArray1 []interface{}, anArray2 []interface{}, key interface{}) {
	valArrIntYaml = nil
	valArrStrYaml = nil
	if len(anArray1) > len(anArray2) {
		deleteArrayYaml(anArray1, anArray2, key)
	} else {
		for i, val := range anArray2 {
			switch concreteVal := val.(type) {
			case map[interface{}]interface{}:
				if len(anArray1) > i {
					parseMapYaml(anArray1[i].(map[interface{}]interface{}), anArray2[i].(map[interface{}]interface{}))
				} else {
					parseMapYaml(nil, anArray2[i].(map[interface{}]interface{}))
				}
				// if _, ok := anArray1[i].(map[interface{}]interface{}); ok {
				// 	parseMapYaml(anArray1[i].(map[interface{}]interface{}), anArray2[i].(map[interface{}]interface{}))
				// } else {
				// 	parseMapYaml(nil, anArray2[i].(map[interface{}]interface{}))
				// }
			case []interface{}:
				if len(anArray1) > i {
					parseArrayYaml(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				} else {
					parseArrayYaml(nil, anArray2[i].([]interface{}), key)
				}
				// if _, ok := anArray1[i].([]interface{}); ok {
				// 	parseArrayYaml(anArray1[i].([]interface{}), anArray2[i].([]interface{}), key)
				// } else {
				// 	parseArrayYaml(nil, anArray2[i].([]interface{}), key)
				// }
			default:
				if len(anArray1) > i {
					if anArray1[i] != concreteVal {
						if _, ok := concreteVal.(string); ok {
							valArrStrYaml = append(valArrStrYaml, concreteVal.(string))
							outputDiffYaml[key] = valArrStrYaml
							oldval := fmt.Sprintf("%v", anArray1[i].(string))
							newval := fmt.Sprintf("%v", concreteVal.(string))
							outputStringYaml += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						} else {
							valArrIntYaml = append(valArrIntYaml, concreteVal.(int))
							outputDiffYaml[key] = valArrIntYaml
							oldval := fmt.Sprintf("%v", anArray1[i].(int))
							newval := fmt.Sprintf("%v", concreteVal.(int))
							outputStringYaml += "Index " + strconv.Itoa(i) + " changed " + oldval + " to " + newval + "\n"
						}
					}
				} else {
					if _, ok := concreteVal.(string); ok {
						valArrStrYaml = append(valArrStrYaml, concreteVal.(string))
						outputDiffYaml[key] = valArrStrYaml
						newval := fmt.Sprintf("%v", concreteVal.(string))
						outputStringYaml += "nil to " + newval + "\n"

					} else {
						valArrIntYaml = append(valArrIntYaml, concreteVal.(int))
						outputDiffYaml[key] = valArrIntYaml
						newval := fmt.Sprintf("%v", concreteVal.(int))
						outputStringYaml += "nil to " + newval + "\n"
					}
				}

			}
		}
	}
}
