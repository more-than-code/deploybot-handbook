package util

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestTaringFiles(t *testing.T) {
	buf, err := TarFiles("/Users/binartist/Projects/deploybot/geoy-webapp/")

	if err != nil {
		t.Fatal(err)
	}

	// Open and iterate through the files in the archive.
	tr := tar.NewReader(buf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			t.Fatal(err)
		}
		fmt.Println()
	}
}
