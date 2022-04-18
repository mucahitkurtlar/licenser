package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var LicensePath = "/opt/licenser/licenses/"

func main() {
	// if no parameters available print usage and exit the program
	if len(os.Args) <= 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: go run main.go <source file> <license name>")
		fmt.Println("Available licenses: mit, gpl, apache, lgpl, mpl")

		os.Exit(1)
	}

	var commentBegin string
	var commentEnd string
	var licenseStr string
	var sourceFileName = os.Args[1]
	var licenseName = os.Args[2]

	commentBegin, commentEnd = getCommentSymbols(sourceFileName)
	licenseStr, err := getLicense(licenseName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = appendLicense(sourceFileName, licenseStr, commentBegin, commentEnd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func getCommentSymbols(fileName string) (commentBegin string, commentEnd string) {

	// print files extension
	fileExt := filepath.Ext(fileName)

	// if case is Go, C, C++, Rust, Java, Kotlin, D, C#, Objective-C, PHP etc. set comment begin and end
	switch fileExt {
	case ".py":
		commentBegin = "\"\"\""
		commentEnd = "\"\"\""
	case ".rb":
		commentBegin = "=begin"
		commentEnd = "=end"
	case ".pl":
		commentBegin = "=pod"
		commentEnd = "=cut"
	case ".sh":
		commentBegin = ": '"
		commentEnd = "'"
	default:
		commentBegin = "/*"
		commentEnd = "*/"
	}

	return
}

func getLicense(licenseName string) (licenseStr string, err error) {
	switch licenseName {
	case "mit":
		// read license file and read the content into a License variable
		licenseStr, err = readLicense("../licenses/mit.txt")

	case "gpl":
		licenseStr, err = readLicense("../licenses/gpl.txt")

	case "apache":
		licenseStr, err = readLicense("../licenses/apache.txt")

	case "lgpl":
		licenseStr, err = readLicense("../licenses/lgpl.txt")
	case "mpl":
		licenseStr, err = readLicense("../licenses/mpl.txt")

	default:
		fmt.Println("Unknown License")
		os.Exit(1)
	}

	return
}

func readLicense(licenseFileName string) (licenseStr string, err error) {
	// open license file and read the content into a License variable
	licenseFile, err := os.ReadFile(LicensePath + licenseFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	licenseStr = string(licenseFile)

	return
}

func appendLicense(sourceFileName string, licenseStr string, commentBegin string, commentEnd string) error {
	tempFile, err := os.Create("license_swap_temp.txt")
	if err != nil {
		return err
	}

	// open the file to be appended to for read
	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		os.Remove("license_swap_temp.txt")
		return err
	}

	// append at the start
	_, err = tempFile.WriteString(commentBegin + "\n" + licenseStr + commentEnd + "\n\n")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(sourceFile)

	// read the file to be appended to and output all of it
	for scanner.Scan() {
		_, err = tempFile.WriteString(scanner.Text())
		if err != nil {
			return err
		}

		_, err = tempFile.WriteString("\n")
		if err != nil {
			return err
		}

	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	// ensure all lines are written
	err = tempFile.Sync()
	if err != nil {
		return err
	}

	// over write the old file with the new one
	err = os.Rename("license_swap_temp.txt", sourceFileName)
	if err != nil {
		return err
	}

	tempFile.Close()
	sourceFile.Close()

	return err
}
