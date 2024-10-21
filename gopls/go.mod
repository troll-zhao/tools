module golang.custom.org/x/tools/gopls

// go 1.23.1 fixes some bugs in go/types Alias support.
// (golang/go#68894 and golang/go#68905).
go 1.23.1

require (
	github.com/google/go-cmp v0.6.0
	github.com/jba/templatecheck v0.7.0
	golang.custom.org/x/tools v0.0.0-00010101000000-000000000000
	golang.org/x/mod v0.21.0
	golang.org/x/sync v0.8.0
	golang.org/x/sys v0.26.0
	golang.org/x/telemetry v0.0.0-20240927184629-19675431963b
	golang.org/x/text v0.19.0
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
	golang.org/x/vuln v1.0.4
	gopkg.in/yaml.v3 v3.0.1
	honnef.co/go/tools v0.4.7
	mvdan.cc/gofumpt v0.7.0
	mvdan.cc/xurls/v2 v2.5.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/google/safehtml v0.1.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20221212164502-fae10dda9338 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect

)

replace golang.custom.org/x/tools => ../
