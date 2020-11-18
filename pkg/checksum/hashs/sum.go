package hashs

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

// calculate sum of file, with filename, as checksum
func SumFilename(h hash.Hash, file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return SumLineFromReader(h, f, file)
}

// calculate sum of file
func SumFile(h hash.Hash, file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return SumReader(h, f)
}

// calculate sum of reader
func SumLineFromReader(h hash.Hash, r io.Reader, filename string) (string, error) {
	sum, err := SumReader(h, r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s  %s\n", sum, filename), nil
}

func SumSfvLine(filename string, sum string) string {
	return fmt.Sprintf("%s %s\n", filename, sum)
}

func SumLine(filename string, sum string) string {
	return fmt.Sprintf("%s *%s\n", sum, filename)
}

func SumFromSumSfvLine(line string) string {
	split := strings.Split(line, " ")
	return strings.TrimSpace(split[len(split)-1])
}

func SumFromSumLine(line string) string {
	return strings.Split(line, " ")[0]
}

// calculate sum of reader
func SumReader(h hash.Hash, r io.Reader) (string, error) {
	h.Reset()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
