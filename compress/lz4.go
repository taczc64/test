package compress

import (
	"fmt"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
)

func LZ4Compress() {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, dir+file.Name())
	}
	err = archiver.TarLz4.Make("cailiao.tar.lz4", filePaths)
	if err != nil {
		fmt.Println(err)
	}
}

func LZ4Decompress() {
	os.Mkdir("lz4dir/", 0777)
	err := archiver.TarLz4.Open("cailiao.tar.lz4", "lz4dir")
	if err != nil {
		fmt.Println(err)
	}
}
