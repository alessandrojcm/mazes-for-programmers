module mfp

go 1.19

require (
	github.com/gen2brain/raylib-go/raylib v0.0.0-20221023125634-dfbb5ab75453
	github.com/stretchr/testify v1.8.1
)

// Replace with my own fork while this is merged: https://github.com/gen2brain/raylib-go/pull/210
replace github.com/gen2brain/raylib-go/raylib v0.0.0-20221023125634-dfbb5ab75453 => github.com/alessandrojcm/raylib-go/raylib v0.0.0-20221031141958-971b4522172c

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
