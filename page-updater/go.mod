module github.com/taichi765/kokkimusume-wiki-automation/page-updater

go 1.26.5

require github.com/joho/godotenv v1.5.1

require github.com/taichi765/kokkimusume-wiki-automation/types v0.0.0

require github.com/taichi765/kokkimusume-wiki-automation/wikiwiki v0.0.0

replace github.com/taichi765/kokkimusume-wiki-automation/types => ../types/

replace github.com/taichi765/kokkimusume-wiki-automation/wikiwiki => ../wikiwiki/

require github.com/stretchr/testify v1.12.1

require go.yaml.in/yaml/v3 v3.0.5 // indirect
