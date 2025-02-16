module github.com/jippi/dottie

go 1.24.0

replace github.com/go-playground/validator/v10 => github.com/jippi/go-validator/v10 v10.0.0-20240202193343-be965b89f3aa

replace github.com/reugn/pkgslog => github.com/jippi/pkgslog v0.0.0-20240224183226-3fdc9d9d89a3

replace github.com/reeflective/console => github.com/jippi/go-console v0.0.0-20240302001452-e5453feb8929

// replace github.com/reeflective/console => ./3rd-party/go-console

require (
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/caarlos0/go-version v0.2.0
	github.com/carapace-sh/carapace v1.7.1
	github.com/charmbracelet/huh v0.6.0
	github.com/charmbracelet/lipgloss v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-playground/validator/v10 v10.24.0
	github.com/golang-cz/devslog v0.0.11
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gosimple/slug v1.15.0
	github.com/hashicorp/go-getter/v2 v2.2.3
	github.com/lmittmann/tint v1.0.7
	github.com/muesli/termenv v0.15.3-0.20241212154518-8c990cd6cf4b
	github.com/neilotoole/slogt v1.1.0
	github.com/reeflective/console v0.1.22
	github.com/reugn/pkgslog v0.2.0
	github.com/samber/slog-multi v1.4.0
	github.com/sebdah/goldie/v2 v2.5.5
	github.com/spf13/cobra v1.9.0
	github.com/spf13/pflag v1.0.6
	github.com/stretchr/testify v1.10.0
	github.com/teacat/noire v1.1.0
	github.com/veqryn/slog-context v0.7.0
	github.com/veqryn/slog-dedup v0.5.0
	go.uber.org/multierr v1.11.0
	mvdan.cc/sh/v3 v3.10.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/carapace-sh/carapace-shlex v1.0.1 // indirect
	github.com/catppuccin/go v0.3.0 // indirect
	github.com/charmbracelet/bubbles v0.20.0 // indirect
	github.com/charmbracelet/bubbletea v1.3.3 // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/charmbracelet/x/exp/strings v0.0.0-20250213125511-a0c32e22e4fc // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/reeflective/readline v1.0.13 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rsteube/carapace v0.46.3-0.20231214181515-27e49f3c3b69 // indirect
	github.com/rsteube/carapace-shlex v0.1.1 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/exp v0.0.0-20250215185904-eff6e970281f // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/b/v2 v2.1.2 // indirect
)
