package main

// main is a dummy function to satisfy the cgo requirements.
func main() {}

// DEV_NOTES:
//
// 1. Functions intended to be exported to C MUST:
//   1a. have an `//export <func_name>` comment on the line preceding their declaration.
//   1b. be declared in this `main` package.
// 2. C types which are included in one package NEVER match the same C types imported from another package.
//
// For more on cgo, see: https://pkg.go.dev/cmd/cgo
