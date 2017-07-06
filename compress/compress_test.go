package compress

//main to test rate of comress, speed of compress and decompress

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkZipCompress(b *testing.B) {
	os.Remove("cailiao.zip")
	for i := 0; i < b.N; i++ {
		ZipCompress()
	}
	zipCompressRate()
}

// func BenchmarkLZ4Compress(b *testing.B) {
//
// }
//
// func BenchmarkZipDecompress(b *testing.B) {
//
// }
//
// func BenchmarkLZ4Decompress(b *testing.B) {
//
// }

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

func zipCompressRate() {
	dirsize := getDirSize(dir)
	info, err := os.Stat("cailiao.zip")
	if err != nil {
		fmt.Println(err)
	}
	rate := float64(info.Size()) / float64(dirsize)
	fmt.Println("zip compress rate >>>>>", rate*100)
}
