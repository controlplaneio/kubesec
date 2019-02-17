package rules

import (
	"bytes"
	"github.com/thedevsaddam/gojsonq"
)

func RequestsCPU(json []byte) int {
	spec := getSpecSelector(json)
	found := 0

	paths := gojsonq.New().Reader(bytes.NewReader(json)).
		From(spec + ".containers").
		Only(".resources.requests.cpu")

	// TODO(ajm) this doesn't currently validate the actual value exists
	if paths != nil {
		found++
	}

	//if paths != nil && paths.Count() > 0 {
	//  found++
	//  fmt.Println("FOUND path", paths.Count())
	//  paths := paths.Get()
	//  fmt.Println("FOUND path", paths)
	//  fmt.Println("Paths: ", paths)
	//  fmt.Printf("%#v\n", paths)
	//  fmt.Printf("%v\n", paths)
	//} else {
	//  fmt.Printf("%v\n", paths.Get())
	//}

	//if paths != nil && !strings.HasSuffix(fmt.Sprintf("%v", paths), "[]") {
	// found++
	//}
	//

	return found
}
