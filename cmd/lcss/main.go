package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var version = "dev"
var commit = ""
var date = ""

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	versionJSON := flag.Bool("version-json", false, "Print version info as JSON and exit")

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stdout, "lcss - Lattice CSS compiler")
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Usage:")
		_, _ = fmt.Fprintln(os.Stdout, "  lcss [options]")
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag || *versionJSON {
		info := struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
			Date    string `json:"date"`
		}{
			Version: version,
			Commit:  commit,
			Date:    date,
		}

		if *versionJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(info)
			os.Exit(0)
		}

		fmt.Printf("lcss %s\ncommit: %s\nbuilt: %s\n", info.Version, info.Commit, info.Date)
		os.Exit(0)
	}

	flag.Usage()
}
