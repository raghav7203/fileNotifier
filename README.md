# File Notifier for Go
#### It notifies what parameters have been changed in config file.

## To install:
```console
go get github.com/raghav7203/fileNotifier
```
<h4>Cross platform: Windows, Linux.</h4>

## Usage:
- For JSON config
```go
package main

import (
	"github.com/raghav7203/fileNotifier"
)

func main() {
	fileNotifier.AddJson("./config.json", change)
}
func change(m map[string]interface{}, s string) {
	fmt.Println(m) 
	fmt.Println(s)
}
``` 
- For YAML config
```go
package main

import (
	"github.com/raghav7203/fileNotifier"
)

func main() {
	fileNotifier.AddYaml("./config.yaml", change)
}
func change(m map[interface{}]interface{}, s string) {
	fmt.Println(m)
	fmt.Println(s)
}
``` 
