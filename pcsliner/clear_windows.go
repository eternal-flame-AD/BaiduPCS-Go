package pcsliner

import (
	_ "unsafe" // for go:linkname

	"github.com/peterh/liner"
)

//go:linkname eraseScreen github.com/eternal-flame-AD/BaiduPCS-Go/vendor/github.com/peterh/liner.(*State).eraseScreen
func eraseScreen(s *liner.State)

// ClearScreen 清空屏幕
func (pl *PCSLiner) ClearScreen() {
	eraseScreen(pl.State)
}
