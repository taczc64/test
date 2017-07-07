package compress

//main to test rate of comress, speed of compress and decompress

import (
	"os"
	"testing"
)

func BenchmarkZipCompress(b *testing.B) {
	os.Remove("cailiao.zip")
	os.RemoveAll("./zipdir/")
	for i := 0; i < b.N; i++ {
		ZipCompress()
	}
	// CompressRate("cailiao.zip")
}

func BenchmarkLZ4Compress(b *testing.B) {
	os.Remove("cailiao.tar.lz4")
	for i := 0; i < b.N; i++ {
		LZ4Compress()
	}
	// CompressRate("cailiao.tar.lz4")
}

func BenchmarkZipDecompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ZipDecompress()
	}
}

func BenchmarkLZ4Decompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LZ4Decompress()
	}
}
