package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Package represents the fields that program requires
// from the response of NPM API
type Package struct {
	Name     string             `json:"name"`
	Tags     Tags               `json:"dist-tags"`
	Versions map[string]Version `json:"versions"`
	Error    string             `json:"error"`
}

type Tags struct {
	Latest string `json:"latest"`
}

type Version struct {
	Dependencies map[string]string `json:"dependencies"`
	Dist         Distribution      `json:"dist"`
}

type Distribution struct {
	Size int `json:"unpackedSize"`
}

// PackageList holds "list" map which hold name of the pacakge and
// its size.
type PackageList struct {
	list map[string]int
	// mutex to protect data from being accessed by multiple goroutines
	// concurrently.
	mu sync.Mutex
}

var wg sync.WaitGroup // wg waits for all goroutines to finish execution.

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a package name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("-- Your input was empty")
			continue
		}

		pkgList := &PackageList{
			list: make(map[string]int),
		}

		wg.Add(1)
		go makeApiCall(input, pkgList)

		wg.Wait()
		calculateTotalPkgSize(pkgList.list)
	}
}

// makeApiCall checks if the package already exists in the pkgList,
// and if it exits, the function returns immediately.
//
// Else, an HTTP request is made to the npmjs api and data is parsed.
func makeApiCall(packageName string, pkgList *PackageList) {
	defer wg.Done()

	pkgList.mu.Lock()
	_, ok := pkgList.list[packageName]
	pkgList.mu.Unlock()
	if ok {
		return
	}

	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	parseData(body, pkgList)
}

// parseData takes in "data" - data to be parsed, and
// pkgList - an instance of PackageList
//
// After unmarshalling, mutex is locked, and pkgList
// is updated with package and size.
//
// makeApiCall is recursively called for all other
// external dependencies.
func parseData(data []byte, pkgList *PackageList) {
	var p Package
	err := json.Unmarshal(data, &p)
	if err != nil {
		log.Fatalln(err)
	}

	if p.Error == "Not found" {
		fmt.Println("Package Not Found. Please try with a different name")
		return
	}

	pkgName := p.Name
	latestVersion := p.Tags.Latest
	latestVersionSize := p.Versions[latestVersion].Dist.Size
	externalDependencies := p.Versions[latestVersion].Dependencies

	pkgList.mu.Lock()
	pkgList.list[pkgName] = latestVersionSize
	pkgList.mu.Unlock()

	for externalPkg := range externalDependencies {
		wg.Add(1)
		go makeApiCall(externalPkg, pkgList) // recursively make api calls for all external packages.
	}
}

func calculateTotalPkgSize(list map[string]int) {
	total := 0

	fmt.Println()
	fmt.Println("List of all the dependant packages and their size")
	for k, v := range list {
		fmt.Printf("%s : %d\n", k, v)
		total += v
	}

	fmt.Println()
	fmt.Printf("Estimated Total Size: %.2f MB\n", float64(total)/1024.0/1024.0)
	fmt.Println()
}
