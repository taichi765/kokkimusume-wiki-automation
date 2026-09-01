module github.com/taichi765/kokkimusume-wiki-automation/deletion-detector

go 1.26.5

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.23.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.14.1
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.8.0
	github.com/AzureAD/microsoft-authentication-library-for-go v1.8.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
)

require github.com/taichi765/kokkimusume-wiki-automation/wikiwiki v0.0.0

require (
	github.com/alecthomas/repr v0.4.0 // indirect
	github.com/disgoorg/godave v0.1.0 // indirect
	github.com/disgoorg/json/v2 v2.0.0 // indirect
	github.com/disgoorg/omit v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad // indirect
)

replace github.com/taichi765/kokkimusume-wiki-automation/wikiwiki => ../wikiwiki/

require github.com/alecthomas/assert/v2 v2.11.0

require (
	github.com/disgoorg/disgo v0.19.6
	github.com/disgoorg/snowflake/v2 v2.0.3
)
