package juniper

import "testing"

func TestPrepareDiff(t *testing.T) {
	testCases := []struct {
		desc string
		cfg  []string
	}{
		{
			desc: "interface config",
			cfg: []string{
				"set interfaces ge-0/0/0 flexible-vlan-tagging",
				"set interfaces ge-0/0/0 unit 12 vlan-id 12",
				"set interfaces ge-0/0/0 unit 12 family inet address 172.16.12.2/24",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			diffCfg := prepareDiff(tC.cfg)
			for _, ln := range diffCfg {
				t.Log(ln)
			}
		})
	}
}
