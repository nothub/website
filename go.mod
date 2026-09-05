module github.com/nothub/website

go 1.27.1

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
	github.com/dlclark/regexp2/v2 v2.2.1 // indirect
	github.com/nothub/ocipack v0.3.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

tool github.com/nothub/ocipack/cmd/ocipack
