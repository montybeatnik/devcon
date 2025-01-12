package juniper

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
)

// JuniperClient represents a Juniper Device's SSH client.
// It holds methods to interact with the remote device, such as
// operational mode and configuration commands.
type JuniperClient struct {
	SSHClient *devcon.SSHClient
}

// NewJuniperClient is a factory function that sets up a
// JuniperClient.
func NewJuniperClient(user, target string, opts ...devcon.Option) *JuniperClient {
	khfp := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		target,
		devcon.WithPassword(os.Getenv("SSH_PASSWORD")),
		devcon.WithTimeout(time.Second*1),
		devcon.WithHostKeyCallback(khfp),
	)
	return &JuniperClient{SSHClient: client}
}

// InterfacesTerse establishes a remote connection to the device,
// runs the 'show interfaces terse' command, asking for the output in XML, and
// unmarshals the output into an InterfaceTerse struct.
func (jc *JuniperClient) InterfacesTerse() (InterfaceTerse, error) {
	cmd := "show interfaces terse | display xml"
	out, err := jc.SSHClient.Run(cmd)
	if err != nil {
		return InterfaceTerse{}, fmt.Errorf("command failed: %w", err)
	}

	var intTerse InterfaceTerse
	if err := xml.Unmarshal([]byte(out), &intTerse); err != nil {
		return InterfaceTerse{}, fmt.Errorf("unmarshal failed: %w", err)
	}
	return intTerse, nil
}

// prepareDiff takes in a config and adds the necessary wrapper syntax
// to apply the conig, print a diff, and roll it back.
func prepareDiff(cfg []string) []string {
	prepCfg := []string{
		"configure private",
		"show | compare",
		"rollback",
		"exit",
		"exit",
	}
	prepCfg = append(prepCfg[:1], append(cfg, prepCfg[1:]...)...)
	return prepCfg
}

// prepareCfg takes in the actual config and inserts
// it into the command necessary to drop into config
// mode, print a diff, and exit config mode and then
// the device altogether.
// TODO:
// take in a commit message.
func prepareCfg(cfg []string) []string {
	prepCfg := []string{
		"configure private",
		"show | compare",
		"commit and-quit",
		"exit",
	}
	prepCfg = append(prepCfg[:1], append(cfg, prepCfg[1:]...)...)
	return prepCfg
}

// Diff logs into the Juniper device, applies the config, prints
// a diff, rolls it back, and then logs out of the device.
func (jc *JuniperClient) Diff(cfg []string) (string, error) {
	cfg = prepareDiff(cfg)
	return jc.SSHClient.RunAll(cfg...)
}

// ApplyConfig logs into the Juniper device, applies the config, prints
// a diff, commits the config, and then logs out of the device.
func (jc *JuniperClient) ApplyConfig(cfg []string) (string, error) {
	cfg = prepareCfg(cfg)
	return jc.SSHClient.RunAll(cfg...)
}

// BGPSummary establishes a remote connection to the device,
// runs the bgp summary command, asking for the output in XML, and
// unmarshals the output into a BGPSummary struct.
func (jc *JuniperClient) BGPSummary() (BGPSummary, error) {
	cmd := "show bgp summary | display xml"
	output, err := jc.SSHClient.Run(cmd)
	if err != nil {
		return BGPSummary{}, fmt.Errorf("failed to run cmd: %w", err)
	}
	var bgpSummary BGPSummary
	if err := xml.Unmarshal([]byte(output), &bgpSummary); err != nil {
		return BGPSummary{}, fmt.Errorf("failed to unmarshal output: %w", err)
	}
	return bgpSummary, nil
}
