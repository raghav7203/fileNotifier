package fileNotifier

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/yaml.v2"
)

func TestFileNotifier(t *testing.T) {
	Convey("matching output string", t, func() {
		old_config := make(map[string]interface{})
		new_config := make(map[string]interface{})
		body1, _ := ioutil.ReadFile("./config1.json")
		body2, _ := ioutil.ReadFile("./config2.json")
		json.Unmarshal(body1, &old_config)
		json.Unmarshal(body2, &new_config)
		parseMapJson(old_config, new_config)
		expectedOutputString := "param2r changed nil  to value2\np2r changed nil  to v2r\np3r changed nil  to v3r\np41r changed nil  to v4ree\nnil to 2\nnil to 3\nnil to 1\nnil to 2\nnil to 7\nnil to 52\nnil to 6\nnil to 8\nnil to 3\np1r changed nil  to v1r\nparam1r changed nil  to value1\nName1 changed Raghav to Raghav Sharma\nIndex 0 changed 74046 to 7404681005\nparam1 changed value1 to value12\np41 changed v4 to v42\nIndex 1 changed 2 to 22\nIndex 4 changed 8 to 76\nnil to 8\np1 changed v1 to v12\nIndex 4 changed 73 to 234\nIndex 5 changed 7 to 73\nnil to 7\nName changed nil  to Raghav Sharma\nnil to 7404681005\nq changed e to er\n"
		arr1 := strings.Split(outputStringJson, "\n")
		sort.Strings(arr1)
		arr2 := strings.Split(expectedOutputString, "\n")
		sort.Strings(arr2)
		So(arr1, ShouldResemble, arr2)
	})

	Convey("matching output map", t, func() {
		old_config := make(map[string]interface{})
		new_config := make(map[string]interface{})
		body1, _ := ioutil.ReadFile("./config1.json")
		body2, _ := ioutil.ReadFile("./config2.json")
		json.Unmarshal(body1, &old_config)
		json.Unmarshal(body2, &new_config)
		parseMapJson(old_config, new_config)
		expectedOutputMap := make(map[string]interface{})
		var p1r interface{} = "v1r"
		var p1 interface{} = "v12"
		var p41 interface{} = "v42"
		var p2r interface{} = "v2r"
		var p3r interface{} = "v3r"
		var p41r interface{} = "v4ree"
		var param1r interface{} = "value1"
		var param1 interface{} = "value12"
		var q interface{} = "er"
		var param2r interface{} = "value2"
		var hobbiesr interface{} = []float64{1, 2, 7, 52, 6, 8, 3}
		var hobbies interface{} = []float64{22, 76, 8}
		var hobbiess interface{} = []float64{234, 73, 7}
		var hobbies2 interface{} = []float64{2, 3}
		var Name interface{} = "Raghav Sharma"
		var PhoneNumber interface{} = []string{"7404681005"}
		var Name1 interface{} = "Raghav Sharma"
		var PhoneNumber1 interface{} = []string{"7404681005"}

		expectedOutputMap["p1r"] = p1r
		expectedOutputMap["p1"] = p1
		expectedOutputMap["p41"] = p41
		expectedOutputMap["p2r"] = p2r
		expectedOutputMap["p3r"] = p3r
		expectedOutputMap["p41r"] = p41r
		expectedOutputMap["param1r"] = param1r
		expectedOutputMap["param1"] = param1
		expectedOutputMap["q"] = q
		expectedOutputMap["param2r"] = param2r
		expectedOutputMap["hobbiesr"] = hobbiesr
		expectedOutputMap["hobbies"] = hobbies
		expectedOutputMap["hobbiess"] = hobbiess
		expectedOutputMap["hobbies2"] = hobbies2
		expectedOutputMap["Name"] = Name
		expectedOutputMap["PhoneNumber"] = PhoneNumber
		expectedOutputMap["Name1"] = Name1
		expectedOutputMap["PhoneNumber1"] = PhoneNumber1

		So(outputDiffJson, ShouldResemble, expectedOutputMap)
	})

	Convey("matching output string yaml", t, func() {
		old_config := make(map[interface{}]interface{})
		new_config := make(map[interface{}]interface{})
		body1, _ := ioutil.ReadFile("./config1.yaml")
		body2, _ := ioutil.ReadFile("./config2.yaml")
		yaml.Unmarshal(body1, &old_config)
		yaml.Unmarshal(body2, &new_config)
		parseMapYaml(old_config, new_config)
		expectedOutputString := "param2r changed nil  to value2\np2r changed nil  to v2r\np3r changed nil  to v3r\np41r changed nil  to v4ree\nnil to 2\nnil to 3\nnil to 1\nnil to 2\nnil to 7\nnil to 52\nnil to 6\nnil to 8\nnil to 3\np1r changed nil  to v1r\nparam1r changed nil  to value1\nName1 changed Raghav to Raghav Sharma\nIndex 0 changed 74046 to 7404681005\nparam1 changed value1 to value12\np41 changed v4 to v42\nIndex 1 changed 2 to 22\nIndex 4 changed 8 to 76\nnil to 8\np1 changed v1 to v12\nIndex 4 changed 73 to 234\nIndex 5 changed 7 to 73\nnil to 7\nName changed nil  to Raghav Sharma\nnil to 7404681005\nq changed e to er\n"
		arr1 := strings.Split(outputStringYaml, "\n")
		sort.Strings(arr1)
		arr2 := strings.Split(expectedOutputString, "\n")
		sort.Strings(arr2)
		So(arr1, ShouldResemble, arr2)
	})

	Convey("clear json", t, func() {
		clearAllJson()
		So(outputStringJson, ShouldBeBlank)
		So(len(outputDiffJson), ShouldEqual, 0)

	})

	Convey("clear yaml", t, func() {
		clearAllYaml()
		So(outputStringYaml, ShouldBeBlank)
		So(len(outputDiffYaml), ShouldEqual, 0)
	})

	Convey("check extension json", t, func() {
		So(checkExtJson("conf.json"), ShouldBeTrue)
		So(checkExtJson("conf.jso"), ShouldBeFalse)
	})

	Convey("check extension yaml", t, func() {
		So(checkExtYaml("conf.yaml"), ShouldBeTrue)
		So(checkExtYaml("conf.yam"), ShouldBeFalse)
	})

	Convey("matching output string in delete(json)", t, func() {
		old_config := make(map[string]interface{})
		new_config := make(map[string]interface{})
		body1, _ := ioutil.ReadFile("./config2.json")
		body2, _ := ioutil.ReadFile("./config1.json")
		json.Unmarshal(body1, &old_config)
		json.Unmarshal(body2, &new_config)
		deleteMapJson(old_config, new_config)
		expectedOutputString := "param1 changed value12 to value1\np41 changed v42 to v4\nIndex 1 changed 22 to 22\nIndex 4 changed 76 to 76\n8 to nil\np1 changed v12 to v1\nparam1r changed value1 to nil\nparam2r changed value2 to nil\np2r changed v2r to nil\np3r changed v3r to nil\np41r changed v4ree to nil\n1 to nil\n2 to nil\n7 to nil\n52 to nil\n6 to nil\n8 to nil\n3 to nil\np1r changed v1r to nil\nIndex 4 changed 234 to 234\nIndex 5 changed 73 to 73\n7 to nil\n2 to nil\n3 to nil\nq changed e to er\nName changed Raghav Sharma to nil\n7404681005 to nil\nName1 changed Raghav Sharma to Raghav\nIndex 0 changed 7404681005 to 74046\n"
		arr1 := strings.Split(outputStringJson, "\n")
		sort.Strings(arr1)
		arr2 := strings.Split(expectedOutputString, "\n")
		sort.Strings(arr2)
		So(arr1, ShouldResemble, arr2)
	})

	Convey("matching output string in delete(yaml)", t, func() {
		old_config := make(map[interface{}]interface{})
		new_config := make(map[interface{}]interface{})
		body1, _ := ioutil.ReadFile("./config2.yaml")
		body2, _ := ioutil.ReadFile("./config1.yaml")
		json.Unmarshal(body1, &old_config)
		json.Unmarshal(body2, &new_config)
		deleteMapYaml(old_config, new_config)
		expectedOutputString := "param1 changed value12 to value1\np41 changed v42 to v4\nIndex 1 changed 22 to 22\nIndex 4 changed 76 to 76\n8 to nil\np1 changed v12 to v1\nparam1r changed value1 to nil\nparam2r changed value2 to nil\np2r changed v2r to nil\np3r changed v3r to nil\np41r changed v4ree to nil\n1 to nil\n2 to nil\n7 to nil\n52 to nil\n6 to nil\n8 to nil\n3 to nil\np1r changed v1r to nil\nIndex 4 changed 234 to 234\nIndex 5 changed 73 to 73\n7 to nil\n2 to nil\n3 to nil\nq changed e to er\nName changed Raghav Sharma to nil\n7404681005 to nil\nName1 changed Raghav Sharma to Raghav\nIndex 0 changed 7404681005 to 74046\n"
		arr1 := strings.Split(outputStringJson, "\n")
		sort.Strings(arr1)
		arr2 := strings.Split(expectedOutputString, "\n")
		sort.Strings(arr2)
		So(arr1, ShouldResemble, arr2)
	})
}
