// @formatter:off 2>/dev/null
// 2>/dev/null || set -o errexit; set -o nounset
// 2>/dev/null || if [ -f "go.mod" ]; then echo "Unable to exec because a go module is present in the working directory!"; exit 1; fi
// 2>/dev/null || mod_path="hub.lol/foo/bar"; mod_gover="1.19"; mod_pkgs=('github.com/spf13/pflag v1.0.5' 'gopkg.in/yaml.v3 v3.0.1')
//usr/bin/env -S printf "module %s\n\ngo %s\n\nrequire (\n%s\n)" "${mod_path}" "${mod_gover}" "$(IFS=$'\n'; echo "${mod_pkgs[*]}")" > go.mod
//usr/bin/env -S go mod tidy; set +o errexit; go run "$0" "$@"; exit_status="$?"; rm -f go.mod go.sum; exit "${exit_status}"
// @formatter:on

package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type ReleaseInfo struct {
	Version  string    `yaml:"version"`
	Name     string    `yaml:"name"`
	Category string    `yaml:"category"`
	Date     time.Time `yaml:"date"`
}

type Manifest struct {
	Latest []ReleaseInfo `yaml:"latest"`
}

func main() {
	openttd := fetch("https://cdn.openttd.org/openttd-releases/latest.yaml")
	opengfx := fetch("https://cdn.openttd.org/opengfx-releases/latest.yaml")

	fmt.Printf("%++q\n", openttd)
	fmt.Printf("%++q\n", opengfx)
}

func fetch(url string) ReleaseInfo {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		panic(err)
	}

	var data Manifest
	err = yaml.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	// sort releases (latest first)
	sort.Slice(data.Latest, func(i, j int) bool {
		return data.Latest[j].Date.Before(data.Latest[i].Date)
	})

	for _, release := range data.Latest {
		if release.Name == "stable" {
			return release
		}
	}

	panic("No stable release found!")
}
