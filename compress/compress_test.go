package compress

//main to test rate of comress, speed of compress and decompress

import (
	"testing"
)

func BenchmarkZipCompress(b *testing.B) {
	ZipCompress()
}

func BenchmarkSnappyCompress(b *testing.B) {

}

func BenchmarkLZ4Compress(b *testing.B) {

}

func BenchmarkZipDecompress(b *testing.B) {

}

func BenchmarkSnappyDecompress(b *testing.B) {

}

func BenchmarkLZ4Decompress(b *testing.B) {

}
