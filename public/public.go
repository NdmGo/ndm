package public

import (
	"embed"
)

//go:embed conf/*
var Conf embed.FS

//go:embed static
var Static embed.FS

//go:embed template/*
var Template embed.FS
