package main

type lang string

const (
	dockerfile lang = "dockerfile"
	golang     lang = "go"
	java       lang = "java"
	lua        lang = "lua"
	perl       lang = "perl"
	python     lang = "python"
	shell      lang = "shell"
)

var langColors = make(map[lang]string)

func init() {
	langColors[dockerfile] = "#384D54"
	langColors[golang] = "#00ADD8"
	langColors[java] = "#B07219"
	langColors[lua] = "#000080"
	langColors[perl] = "#0298C3"
	langColors[python] = "#3572A5"
	langColors[shell] = "#89E051"
}
