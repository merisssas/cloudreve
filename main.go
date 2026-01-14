//go:debug rsa1024min=0
package main

import (
	_ "embed"
	"flag"
	"os"
	"strings"

	"github.com/merisssas/Cloudreve/v4/cmd"
	"github.com/merisssas/Cloudreve/v4/pkg/util"
)

var (
	confPath   string
	scriptName string
)

// envBool membaca env var boolean yang “manusiawi”.
// Accepted true: 1, true, yes, y, on (case-insensitive)
func envBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

func init() {
	// Default: pakai executable directory (sesuai util.UseWorkingDir = false)
	// Bisa di-override dengan env (berguna untuk Docker/CI):
	//   CLOUDREVE_USE_WORKING_DIR=1
	defaultUseWD := envBool("CLOUDREVE_USE_WORKING_DIR", false)
	util.UseWorkingDir = defaultUseWD

	// Default config Cloudreve v4 adalah data/conf.ini. 2
	defaultConf := os.Getenv("CLOUDREVE_CONFIG")
	if strings.TrimSpace(defaultConf) == "" {
		defaultConf = util.DataPath("conf.ini")
	}

	flag.BoolVar(
		&util.UseWorkingDir,
		"use-working-dir",
		defaultUseWD,
		"Use working directory instead of executable directory (env: CLOUDREVE_USE_WORKING_DIR)",
	)

	flag.StringVar(
		&confPath,
		"c",
		defaultConf,
		"Path to the config file (default: data/conf.ini; env: CLOUDREVE_CONFIG)",
	)

	flag.StringVar(
		&scriptName,
		"database-script",
		"",
		"Name of database util script (reserved; prefer cmd subcommands if available).",
	)

	// IMPORTANT:
	// Jangan flag.Parse() di sini kalau cmd.Execute() memakai Cobra/pflag dan menggabungkan Go flags.
	// Kalau kamu parse di sini, bisa bentrok dengan flags milik Cobra.
}

func main() {
	// Entry point utama Cloudreve CLI (server/migrate/eject/etc).
	cmd.Execute()
}
