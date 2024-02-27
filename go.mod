module github.com/jippi/dottie

go 1.21.1

toolchain go1.21.6

replace github.com/go-playground/validator/v10 => github.com/jippi/go-validator/v10 v10.0.0-20240202193343-be965b89f3aa

replace github.com/reugn/pkgslog => github.com/jippi/pkgslog v0.0.0-20240224183226-3fdc9d9d89a3

require (
	github.com/caarlos0/go-version v0.1.1
	github.com/charmbracelet/huh v0.3.0
	github.com/charmbracelet/lipgloss v0.9.1
	github.com/davecgh/go-spew v1.1.1
	github.com/go-playground/validator/v10 v10.17.0
	github.com/golang-cz/devslog v0.0.8
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gosimple/slug v1.14.0
	github.com/hashicorp/go-getter v1.7.3
	github.com/lmittmann/tint v1.0.4
	github.com/muesli/termenv v0.15.2
	github.com/neilotoole/slogt v1.1.0
	github.com/reugn/pkgslog v0.0.0-20240224183226-3fdc9d9d89a3
	github.com/samber/slog-multi v1.0.2
	github.com/sebdah/goldie/v2 v2.5.3
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	github.com/teacat/noire v1.1.0
	github.com/veqryn/slog-context v0.5.0
	github.com/veqryn/slog-dedup v0.3.0
	go.uber.org/multierr v1.11.0
	mvdan.cc/sh/v3 v3.8.0
)

require (
	cloud.google.com/go v0.110.4 // indirect
	cloud.google.com/go/compute v1.21.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aws/aws-sdk-go v1.44.122 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/catppuccin/go v0.2.0 // indirect
	github.com/charmbracelet/bubbles v0.17.2-0.20240108170749-ec883029c8e6 // indirect
	github.com/charmbracelet/bubbletea v0.25.0 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.10.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/grpc v1.58.3 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/b/v2 v2.1.0 // indirect
)
