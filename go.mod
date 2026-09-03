module github.com/nothub/website

go 1.26.4

require (
	github.com/alecthomas/chroma/v2 v2.27.0
	github.com/elnormous/contenttype v1.0.4
	github.com/mangoumbrella/goldmark-figure v1.4.0
	github.com/spf13/pflag v1.0.10
	github.com/yuin/goldmark v1.8.6
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc
	github.com/yuin/goldmark-meta v1.1.0
	go.abhg.dev/goldmark/anchor v0.2.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	codeberg.org/fhuebner/ocipack v0.2.0 // indirect
	github.com/dlclark/regexp2/v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

// ocipack moved to github, but its go.mod still declares the codeberg module
// path, so the module identity has to stay codeberg until upstream retags.
replace codeberg.org/fhuebner/ocipack => github.com/nothub/ocipack v0.2.0

tool codeberg.org/fhuebner/ocipack/cmd/ocipack
