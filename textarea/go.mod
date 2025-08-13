module github.com/fyrna/x/textarea

go 1.24.5

require (
	github.com/fyrna/x/term v0.0.0
	github.com/fyrna/x/color v0.0.0
)

require (
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
)

replace (
	github.com/fyrna/x/term => ../term
	github.com/fyrna/x/color => ../color
)
