package seed

import "embed"

//go:embed *.sql
var Seed embed.FS
