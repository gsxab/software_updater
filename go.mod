module software_updater

go 1.18

require (
	github.com/fishjar/gin-i18n v0.0.3
	github.com/gin-gonic/gin v1.9.1
	github.com/gsxab/go-generic_lru v1.0.2
	github.com/gsxab/go-error_util v1.0.0
	github.com/gsxab/go-logs v1.0.0
	github.com/gsxab/go-slice_util v0.1.0
	github.com/gsxab/go-optional v0.2.0
	github.com/gsxab/go-version v1.0.0
	github.com/itchyny/gojq v0.12.13
	github.com/tebeka/selenium v0.9.10-0.20211105214847-e9100b7f5ac1
	golang.org/x/net v0.10.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.5.3
	gorm.io/gen v0.3.19
	gorm.io/gorm v1.25.4
	gorm.io/plugin/dbresolver v1.3.0
)

require (
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/mediabuyerbot/go-crx3 v1.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/datatypes v1.0.7 // indirect
	gorm.io/driver/mysql v1.4.0 // indirect
	gorm.io/hints v1.1.0 // indirect
)

replace github.com/gsxab/go-generic_lru => ../generic_lru

replace github.com/gsxab/go-logs => ../logs

replace github.com/gsxab/go-error_util => ../error_util

replace github.com/gsxab/go-version => ../version

replace github.com/gsxab/go-optional => ../optional

replace github.com/gsxab/go-slice_util => ../slice_util
