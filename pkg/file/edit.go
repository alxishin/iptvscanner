package file

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func SetFileContent(path string, data string) {
	file := rewriteFile(path)

	_, err := file.WriteString(data)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	defer file.Close()
}

func GetFileContent(path string) (string, error) {
	b, err := os.ReadFile(path) // just pass the file name
	return string(b), err       // convert content to a 'string'
}

func PrependFile(path string, addline string) {

	outfile, err := os.Create(path + ".tmp")

	if err != nil {
		panic(err)
	}

	defer outfile.Close()

	// open the file to be appended to for read
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	// append at the start
	_, err = outfile.WriteString(addline)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	// read the file to be appended to and output all of it
	for scanner.Scan() {
		_, err = outfile.WriteString(scanner.Text())
		_, err = outfile.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// ensure all lines are written
	outfile.Sync()
	// over write the old file with the new one
	err = os.Rename(path+".tmp", path)
	if err != nil {
		panic(err)
	}
}

func WriteTextBefore(path string, data string, separator string) {
	writeTextBeforeOrAfter(path, data, separator, "before")
}

func WriteTextAfter(path string, data string, separator string) {
	writeTextBeforeOrAfter(path, data, separator, "after")
}

func writeTextBeforeOrAfter(path string, data string, separator string, mode string) {
	outfile, err := os.Create(path + ".tmp")

	if err != nil {
		panic(err)
	}

	defer outfile.Close()

	// open the file to be appended to for read
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	txt := ""
	// read the file to be appended to and output all of it
	for scanner.Scan() {
		txt = scanner.Text()
		arr := strings.Split(txt, separator)

		if len(arr) == 1 {
			_, err = outfile.WriteString(txt)
		} else if len(arr) == 2 {
			if mode == "before" {
				_, err = outfile.WriteString(arr[0] + data + separator + arr[1])

			} else if mode == "after" {
				_, err = outfile.WriteString(arr[0] + separator + data + arr[1])
			} else {
				panic("Wrong mode: " + mode)
			}

		} else {
			panic("WTF")
		}

		_, err = outfile.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// ensure all lines are written
	outfile.Sync()
	// over write the old file with the new one
	err = os.Rename(path+".tmp", path)
	if err != nil {
		panic(err)
	}
}

func AppendFile(path string, data string) {
	file := openFile(path)

	_, err := file.WriteString(data)
	//len, err := file.WriteAt([]byte{'G'}, 0) // Write at 0 beginning
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	defer file.Close()
}

func openFile(path string) *os.File {
	return openFileWithFlag(path, os.O_APPEND|os.O_RDWR)
}

func rewriteFile(path string) *os.File {
	return openFileWithFlag(path, os.O_RDWR|os.O_CREATE)
}

func openFileWithFlag(path string, flag int) *os.File {

	if _, err := os.Stat(path); err == nil {
		file, err := os.OpenFile(path, flag, 0644)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		return file
	} else {
		file, err := os.Create(path)
		// path/to/whatever does *not* exist
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		return file
	}
}
