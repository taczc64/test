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
	os.Mkdir("zipdir/", 0777)
	zipf, err := zip.OpenReader("cailiao.zip")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer zipf.Close()
	for _, file := range zipf.File {
		subzip, err := file.Open()
		if err != nil {
			fmt.Println("open sub zip err :", err)
		}
		f, err := os.Create("zipdir/" + file.Name)
		if err != nil {
			fmt.Println("create sub file err :", err)
			continue
		}
		defer f.Close()
		//TODO dont use copy, just write it block by block
		_, err = io.Copy(f, subzip)
		if err != nil {
			fmt.Println(err)
		}
	}
}
