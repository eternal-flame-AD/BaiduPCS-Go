package pcsverbose

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/eternal-flame-AD/BaiduPCS-Go/pcsutil/pcstime"
)

//PrintReader 输出Reader
func PrintReader(r io.Reader) {
	b, _ := ioutil.ReadAll(r)
	fmt.Printf("%s\n", b)
}

func TimePrefix() string {
	return "[" + pcstime.BeijingTimeOption("Refer") + "]"
}
