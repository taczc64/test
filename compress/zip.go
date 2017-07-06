package compress

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const dir = "./cailiao/"

var buffersize = 4 * 1024 * 1024

func ZipCompress() {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("read dir error :", err)
		return
	}
	zipf, err := os.Create("cailiao.zip")
	if err != nil {
		fmt.Println(err)
		return
	}
	zfw := zip.NewWriter(zipf)
	defer zfw.Close()

	buf := make([]byte, buffersize)
	for _, file := range fs {
		subfw, _ := zfw.Create(file.Name())
		f, err := os.Open(dir + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		freader := bufio.NewReader(f)
		for {
			n, err := freader.Read(buf)
			if err != nil && err == io.EOF {
				break
			}
			_, err = subfw.Write(buf[:n])
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func ZipDecompress() {

}
