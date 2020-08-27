package vjson

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestGo(t *testing.T) {

	tcaseList := []struct {
		name    string
		in      string
		matches []string
	}{
		{name: "simple_map", in: `{"testkey":"testval"}`, matches: []string{`{"testkey":"testval"}`}},
		{name: "simple_map_number", in: `{"teststr":"str1","testnum":15}`, matches: []string{`"teststr":"str1"`, `"testnum":15`}},
		{name: "simple_map_bool", in: `{"testbool1":false,"testbool2":true}`, matches: []string{`"testbool1":false`, `"testbool2":true`}},
		{name: "string_esc", in: `{"testkey":"test\"val"}`, matches: []string{`[{]"testkey":"test\\"val"[}]`}},
		{name: "string_esc2", in: `{"testkey":"test\"val\""}`, matches: []string{`[{]"testkey":"test\\"val\\""[}]`}},
		{name: "string_esc3", in: `{"testkey":"te\\st\"val\""}`, matches: []string{`[{]"testkey":"te\\\\st\\"val\\""[}]`}},
		{name: "null1", in: `{"null1":null}`, matches: []string{`"null1":null`}},
		{name: "array1", in: `{"array1":["s1","s2"]}`, matches: []string{`"array1":\["s1","s2"\]`}},
		{name: "array2", in: `{"array2":[]}`, matches: []string{`"array2":\[\]`}},
		{name: "array_of_obj1", in: `{"array1":[{},{"k1":"v1"}]}`, matches: []string{`"array1":...,."k1":"v1".`}},
	}

	for _, tcase := range tcaseList {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {

			var m map[string]interface{}
			err := unmarshal([]byte(tcase.in), &m)
			if err != nil {
				t.Fatal(err)
			}

			b, err := marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Marshal output: %s", b)

			for _, ma := range tcase.matches {
				if !regexp.MustCompile(ma).Match(b) {
					t.Errorf("failed to match in output: %s", ma)
				}
			}

		})
	}
}

func TestTinygo(t *testing.T) {

	dockerImage := "tinygo/tinygo:latest"

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	pwd, err = filepath.EvalSymlinks(pwd)
	if err != nil {
		t.Fatal(err)
	}

	// now execute it inside a container for each input (since it will be built for that target, not necessarily for the host)

	tcaseList := []struct {
		name    string
		in      string
		matches []string
	}{
		{name: "simple_map", in: `{"testkey":"testval"}`, matches: []string{`{"testkey":"testval"}`}},
		{name: "simple_map_number", in: `{"teststr":"str1","testnum":15}`, matches: []string{`"teststr":"str1"`, `"testnum":15`}},
		{name: "simple_map_bool", in: `{"testbool1":false,"testbool2":true}`, matches: []string{`"testbool1":false`, `"testbool2":true`}},
		{name: "string_esc", in: `{"testkey":"test\"val"}`, matches: []string{`[{]"testkey":"test\\"val"[}]`}},
		{name: "string_esc2", in: `{"testkey":"test\"val\""}`, matches: []string{`[{]"testkey":"test\\"val\\""[}]`}},
		{name: "string_esc3", in: `{"testkey":"te\\st\"val\""}`, matches: []string{`[{]"testkey":"te\\\\st\\"val\\""[}]`}},
		{name: "null1", in: `{"null1":null}`, matches: []string{`"null1":null`}},
		{name: "array1", in: `{"array1":["s1","s2"]}`, matches: []string{`"array1":\["s1","s2"\]`}},
		{name: "array2", in: `{"array2":[]}`, matches: []string{`"array2":\[\]`}},
		{name: "array_of_obj1", in: `{"array1":[{},{"k1":"v1"}]}`, matches: []string{`"array1":...,."k1":"v1".`}},
	}

	for _, tcase := range tcaseList {
		tcase := tcase
		t.Run(tcase.name, func(t *testing.T) {

			pgmText := strings.Replace(testPgmTemplate, `__JSON_IN__`, fmt.Sprintf("%q", tcase.in), 1)
			err := ioutil.WriteFile("vjson_test_pgm.go", []byte(pgmText), 0644)
			if err != nil {
				t.Fatal(err)
			}

			// build test program with test case stuff hard coded into it (os.Stdin doesn't seem to work, no args, no env)
			cmd := exec.Command(
				"docker", "run",
				"--rm",
				"-t",
				// "-a","STDERR","-a","STDOUT",
				"-eGOPATH=/root/go",
				"-v"+pwd+":/root/go/src/github.com/vugu/vjson",
				dockerImage,
				"tinygo", "build",
				"-o", "/root/go/src/github.com/vugu/vjson/vjson_test_pgm.out",
				"/root/go/src/github.com/vugu/vjson/vjson_test_pgm.go",
			)
			b, err := cmd.CombinedOutput()
			log.Printf("BUILD OUTPUT: %s", b)
			if err != nil {
				t.Fatal(err)
			}

			// run it
			cmd = exec.Command(
				"docker", "run",
				"--rm",
				"-t",
				// "-a","STDERR","-a","STDOUT",
				"-eGOPATH=/root/go",
				"-v"+pwd+":/root/go/src/github.com/vugu/vjson",
				dockerImage,
				"/root/go/src/github.com/vugu/vjson/vjson_test_pgm.out",
			)
			// cmd.Stdin = bytes.NewReader([]byte(tcase.in))
			b, err = cmd.CombinedOutput()
			log.Printf("RUN OUTPUT for %s: %s", tcase.name, b)
			if err != nil {
				t.Fatal(err)
			}

			for _, ma := range tcase.matches {
				if !regexp.MustCompile(ma).Match(b) {
					t.Errorf("failed to match in output: %s", ma)
				}
			}

			// time.Sleep(time.Second * 2)

		})
	}

}

var testPgmTemplate = `// +build ignore

package main

// This is a test program that is built and run with Tinygo.  See vjson_test.go.

import (
	"fmt"
	//"io/ioutil"
	//"os"

	"github.com/vugu/vjson"
)

func main() {

	var jsonIn = __JSON_IN__

	var m map[string]interface{}
	err := vjson.Unmarshal([]byte(jsonIn), &m)
	if err != nil {
		panic(err)
	}

	b, err := vjson.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", b)

}
`
