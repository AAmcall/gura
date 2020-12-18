package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/go-git/go-git"
)


type response struct {
	Error       string `json:"error"`
	Version     int    `json:"version"`
	Type        string `json:"type"`
	ResultCount int    `json:"resultcount"`
	Results     []Pkg  `json:"results"`
}

type Pkg struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// Search queries the aur
func Search(values string) {
	var p response
	args := &url.URL{Path: values}
	ur := args.String()
	res, err := http.Get("https://aur.archlinux.org/rpc/?v=5&type=search&arg=" + ur)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}

	str := fmt.Sprintf("%v", p.Results)
	replacer := strings.NewReplacer("{", "", "}", "\n", "[", " ", "]", "")
	replacer2 := strings.NewReplacer("[", "", "]", "")
	test := replacer2.Replace(str)
	output := replacer.Replace(str)
	if test != "" {
		fmt.Println(output)
	} else {
		fmt.Println(values + " not found.")
	}
}

// Git uses git clone to download the pkgbuild from the aur
func Git(pkg string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	
	_, err = git.PlainClone(path + "/"  + pkg, false, &git.CloneOptions{
		URL:      "https://aur.archlinux.org/" + pkg + ".git",
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	var Serpkg = flag.String("Sr", "", "search the aur for a `package`.")
	var Gitpkg = flag.String("Git", "", "use git clone to download a `package` from the aur.")

	flag.Parse()

	switch {
	case len(*Serpkg) > 0 && len(*Gitpkg) > 0:
		fmt.Println("Multiple flags at once is not supported")
	case len(*Gitpkg) > 0:
		Git(*Gitpkg)
		fmt.Println(len(*Serpkg))
	case len(*Serpkg) > 0:
		// Gura has to check if Sr is given multiple search args as appending a whitespace to the end of some package names break the return value.
		if len(flag.Arg(0)) > 0 {
			Search(*Serpkg + " " + (strings.Join(flag.Args(), " ")))
		} else {
			Search(*Serpkg)
		}
	default:
		fmt.Println("no flags given for help use -h")
	}
}
