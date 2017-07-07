package compress

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func getDirSize(dir string) int64 {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var size int64
	for _, file := range fs {
		size += file.Size()
	}
	return size
}

func CompressRate(filename string) {
	dirsize := getDirSize(dir)
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}
	rate := float64(info.Size()) / float64(dirsize)
	fmt.Println(filename, " compress rate >>>>>", rate*100)
}

func TestRatio(t *testing.T) {
	CompressRate("cailiao.zip")
	CompressRate("cailiao.tar.lz4")
}
