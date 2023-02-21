// 2>/dev/null || set -o errexit; set -o nounset
// 2>/dev/null || if [ -f "go.mod" ]; then echo "Unable to exec because a go module is present in the working directory!"; exit 1; fi
// 2>/dev/null || mod_path="hub.lol/foo/bar"; mod_gover="1.19"; mod_pkgs=('github.com/spf13/pflag v1.0.5' 'gopkg.in/yaml.v3 v3.0.1')
//usr/bin/env -S printf "module %s\n\ngo %s\n\nrequire (\n%s\n)" "${mod_path}" "${mod_gover}" "$(IFS=$'\n'; echo "${mod_pkgs[*]}")" > go.mod
//usr/bin/env -S go mod tidy; set +o errexit; go run "$0" "$@"; exit_status="$?"; rm -f go.mod go.sum; exit "${exit_status}"
package main

import (
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// ...
