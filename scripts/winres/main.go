// Generates the Windows resource object (.syso) that `go build` links into
// the executable: application icon, DPI/common-controls manifest, and version
// info. This mirrors what `wails build` does internally (see wails'
// pkg/commands/build/packager.go) using the same assets in build/windows/,
// since GoReleaser invokes plain `go build` which does not generate one.
//
// Invoked by the GoReleaser pre-build hook for the Windows build:
//
//	go run ./scripts/winres -version <semver>
//
// The output filename carries the _windows_amd64 suffix so Go's implicit
// build constraints keep it out of builds for other targets.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
)

const (
	name        = "BetterDiscord Installer"
	company     = "BetterDiscord"
	copyright   = "© BetterDiscord"
	description = "A simple standalone program to install, repair, and uninstall BetterDiscord."
)

func main() {
	ver := flag.String("version", "0.0.0", "release version (without v prefix)")
	out := flag.String("out", "rsrc_windows_amd64.syso", "output file")
	flag.Parse()

	if err := run(*ver, *out); err != nil {
		fmt.Fprintln(os.Stderr, "winres:", err)
		os.Exit(1)
	}
}

func run(ver, out string) error {
	// Assembly identity and fixed file versions must be numeric (x.y.z), so
	// strip any prerelease suffix like "-alpha.1". GoReleaser passes the
	// version without the leading v, but tolerate one for manual invocations.
	numeric := "0.0.0"
	if m := regexp.MustCompile(`^v?(\d+\.\d+\.\d+)`).FindStringSubmatch(ver); len(m) > 1 {
		numeric = m[1]
	}

	rs := winres.ResourceSet{}

	iconFile, err := os.Open("build/windows/icon.ico")
	if err != nil {
		return err
	}
	defer iconFile.Close()
	ico, err := winres.LoadICO(iconFile)
	if err != nil {
		return fmt.Errorf("couldn't load icon from icon.ico: %w", err)
	}
	if err := rs.SetIcon(winres.RT_ICON, ico); err != nil {
		return err
	}

	replacer := strings.NewReplacer(
		"{{.Name}}", name,
		"{{.Info.ProductVersion}}", numeric,
		"{{.Info.CompanyName}}", company,
		"{{.Info.ProductName}}", name,
		"{{.Info.Copyright}}", copyright,
		"{{.Info.Comments}}", description,
	)

	manifestData, err := os.ReadFile("build/windows/wails.exe.manifest")
	if err != nil {
		return err
	}
	manifest, err := winres.AppManifestFromXML([]byte(replacer.Replace(string(manifestData))))
	if err != nil {
		return err
	}
	rs.SetManifest(manifest)

	versionData, err := os.ReadFile("build/windows/info.json")
	if err != nil {
		return err
	}
	var vi version.Info
	if err := vi.UnmarshalJSON([]byte(replacer.Replace(string(versionData)))); err != nil {
		return err
	}
	rs.SetVersionInfo(vi)

	fout, err := os.Create(out)
	if err != nil {
		return err
	}
	if err := rs.WriteObject(fout, winres.ArchAMD64); err != nil {
		fout.Close()
		return err
	}
	// Close explicitly so flush errors fail the build instead of being
	// swallowed by a deferred Close, which could leave a corrupt syso.
	return fout.Close()
}
