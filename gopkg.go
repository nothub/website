package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"path"
	"strings"
)

type GoPkgNfo struct {
	Path string `yaml:"path"`
	Repo string `yaml:"repo"`
}

var pkgs = make(map[string]GoPkgNfo)

func initGoPkgs(router *gin.Engine) (err error) {
	log.Println("aliasing go packages")

	bytes, err := fs.ReadFile("data/gopkgs.yaml")
	if err != nil {
		return err
	}
	var pkgInfos []GoPkgNfo
	err = yaml.Unmarshal(bytes, &pkgInfos)
	if err != nil {
		return err
	}
	for _, pkg := range pkgInfos {
		pkgs[path.Base(pkg.Path)] = pkg
	}

	log.Printf("go packages aliased: %q\n", maps.Keys(pkgs))

	router.Use(gopkg)

	return nil
}

var gopkg = func(ctx *gin.Context) {
	if ctx.Query("go-get") != "1" {
		// not a go package request
		ctx.Next()
		return
	}

	reqPath := strings.TrimPrefix(ctx.Request.URL.Path, "/")
	reqPath = strings.TrimSuffix(reqPath, "/")
	reqPath = strings.TrimSpace(reqPath)
	if len(reqPath) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	pkg, ok := pkgs[reqPath]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "gopkg.gohtml", map[string]any{
		"path":   pkg.Path,
		"repo":   pkg.Repo,
		"gopher": gopher(),
	})
}
