package public

import "embed"

//go:embed conf/*
var Conf embed.FS

// //go:embed all:dist
// var Public embed.FS
