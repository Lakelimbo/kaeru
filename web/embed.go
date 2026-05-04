package web

import "embed"

//go:generate pnpm i
//go:generate pnpm build
//go:embed all:build
var Frontend embed.FS
