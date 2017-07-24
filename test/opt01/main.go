package main

import (
	"github.com/cpmech/gosl/io"
	jsoniter "github.com/json-iterator/go"

	"../../inv"
)

func init() {
	io.Verbose = true
}

func main() {
	var opts = make([]*inv.Option, 0)
	for i := 0; i < 3; i++ {
		opts = append(opts, inv.NewOption())
	}
	io.Pf("original:\n%v\n", opts)
	//
	ofn := "./option.json"
	b, err := jsoniter.Marshal(&opts)
	if err != nil {
		panic(err)
	}
	io.WriteBytesToFile(ofn, b)
	//
	bo, err := io.ReadFile(ofn)
	if err != nil {
		panic(err)
	}
	err = jsoniter.Unmarshal(bo, &opts)
	if err != nil {
		panic(err)
	}
	io.Pfgreen2("Unmarshal:\n")
	for i, o := range opts {
		io.Pfgreen2("Option %d:\n%v\n", i+1, o)
	}

}
