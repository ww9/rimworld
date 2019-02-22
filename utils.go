package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "", "\ngenweap: generate stats csv from Vanilla-Friendly Weapon Pack by SY4\ndefault: check mods for problems")
	flag.Parse()

	checkFoldersWithoutAboutXML()
	checkAboutXMLDeprecatedVersion()
}

func checkFoldersWithoutAboutXML() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			if _, err := os.Stat("./" + f.Name() + "/About/About.xml"); os.IsNotExist(err) {
				fmt.Println(f.Name() + " missing /About/About.xml")
			}
		}
	}
}

func checkAboutXMLDeprecatedVersion() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			aboutFile := "./" + f.Name() + "/About/About.xml"
			if _, err := os.Stat(aboutFile); !os.IsNotExist(err) {
				fileBytes, err := ioutil.ReadFile(aboutFile)
				if err != nil {
					fmt.Printf("\nCouldn't read %s: %s", aboutFile, err)
					continue
				}
				exists, err := tagExists(string(fileBytes), "targetVersion")
				if err != nil {
					fmt.Printf("\nCouldn't check targetVersion tag in %s: %s", aboutFile, err)
					continue
				}
				if exists {
					fmt.Printf("\nReplace <targetVersion> with <supportedVersions><li>1.0</li></supportedVersions> in: %s", aboutFile)
				}
			}
		}
	}
}

func tagExists(src, tag string) (bool, error) {
	decoder := xml.NewDecoder(strings.NewReader(src))

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == tag {
				return true, nil
			}
		}
	}
}

// func isDir(filePath string) bool {
// 	fi, err := os.Stat(filePath)
// 	if err != nil {
// 		return false
// 	}
// 	return fi.IsDir()
// }
