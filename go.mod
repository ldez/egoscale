module github.com/exoscale/egoscale

require (
	github.com/BluntSporks/abbreviation v0.0.0-20150522120346-096cdb48bafa
	github.com/deepmap/oapi-codegen v1.9.1
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/google/uuid v1.4.0
	github.com/hashicorp/go-retryablehttp v0.7.1
	github.com/pb33f/libopenapi v0.11.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dprotaso/go-yit v0.0.0-20220510233725-9ba8df137936 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/viper v1.18.2
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/vmware-labs/yaml-jsonpath v0.3.2 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.20

retract (
	v1.19.1 // Retracts the previous version
	v1.19.0 // Published accidentally.
)
