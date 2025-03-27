package ui

import "embed"

//go:embed static/*
var static embed.FS

func GetStatic() embed.FS {
	return static
}
