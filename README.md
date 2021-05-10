# File Notifier for Go
#### It notifies what parameters have been changed in config file.

## To install:
```console
go get github.com/raghav7203/fileNotifier
```
<h4>Cross platform: Windows, Linux.</h4>

## Usage:

```go
package main

import (
	"github.com/raghav7203/fileNotifier"
)

func main() {

fileNotifier.Add("./config.json")

}
```

**Config file type must be JSON**

## To-Do:
- Config file can support YAML extension 
