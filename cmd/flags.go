package main

import "github.com/urfave/cli/v3"

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:       "file",
		Category:   "Input:",
		Value:      ".env",
		Sources:    cli.EnvVars("FILE"),
		Usage:      "Load configuration from `FILE`",
		Persistent: true,
		TakesFile:  true,
		OnlyOnce:   true,
	},

	&cli.BoolFlag{
		Name:       "pretty",
		Category:   "Output:",
		Sources:    cli.EnvVars("PRETTY"),
		Usage:      "implies --with-comments --with-blank-lines --with-groups",
		Persistent: true,
		OnlyOnce:   true,
	},
	&cli.BoolFlag{
		Name:       "with-comments",
		Category:   "Output:",
		Sources:    cli.EnvVars("WITH_COMMENTS"),
		Usage:      "Show comments`",
		Persistent: true,
		OnlyOnce:   true,
	},
	&cli.BoolFlag{
		Name:       "with-blank-lines",
		Category:   "Output:",
		Sources:    cli.EnvVars("WITH_BLANK_LINES"),
		Usage:      "Show blank lines between sections`",
		Persistent: true,
		OnlyOnce:   true,
	},
	&cli.BoolFlag{
		Name:       "with-groups",
		Category:   "Output:",
		Sources:    cli.EnvVars("WITH_GROUPS"),
		Usage:      "Show group banners`",
		Persistent: true,
		OnlyOnce:   true,
	},
	&cli.StringFlag{
		Name:       "key-prefix",
		Category:   "Filters:",
		Sources:    cli.EnvVars("FILTER_KEY_PREFIX"),
		Usage:      "Filter by key prefix`",
		Persistent: true,
		OnlyOnce:   true,
	},
	&cli.StringFlag{
		Name:       "group",
		Category:   "Filters:",
		Sources:    cli.EnvVars("FILTER_GROUP"),
		Usage:      "Filter by group name`",
		Persistent: true,
		OnlyOnce:   true,
	},
	&cli.BoolFlag{
		Name:       "include-commented",
		Category:   "Filters:",
		Sources:    cli.EnvVars("INCLUDE_COMMENTED"),
		Usage:      "Include commented KEY/VALUE pairs`",
		Persistent: true,
		OnlyOnce:   true,
	},
}
