package targz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func Extract(gzipStream io.Reader) {
	uncompressedStream, err := gzip.NewReader(gzipStream)

	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed\n", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
		case tar.TypeSymlink:
			os.Symlink(header.Linkname, header.Name)
		default:
			log.Fatalf(
				"ExtractTarGz: unknown type: %s in %s",
				string(header.Typeflag),
				header.Name)
		}
	}
}
