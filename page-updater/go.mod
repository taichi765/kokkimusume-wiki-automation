module github.com/taichi765/kokkimusume-wiki-automation/page-updater

go 1.26.5

require github.com/joho/godotenv v1.5.1

require github.com/taichi765/kokkimusume-wiki-automation/common v0.0.0

replace github.com/taichi765/kokkimusume-wiki-automation/common => ../common/

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
