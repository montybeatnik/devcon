package juniper

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
)

type InterfaceTerse struct {
	InterfaceInformation struct {
		PhysicalInterface []struct {
			Text             string `xml:",chardata"`
			Name             string `xml:"name"`
			AdminStatus      string `xml:"admin-status"`
			OperStatus       string `xml:"oper-status"`
			LogicalInterface []struct {
				Text              string `xml:",chardata"`
				Name              string `xml:"name"`
				AdminStatus       string `xml:"admin-status"`
				OperStatus        string `xml:"oper-status"`
				FilterInformation string `xml:"filter-information"`
				AddressFamily     []struct {
					Text              string `xml:",chardata"`
					AddressFamilyName string `xml:"address-family-name"`
					InterfaceAddress  []struct {
						Text     string `xml:",chardata"`
						IfaLocal struct {
							Text string `xml:",chardata"`
							Emit string `xml:"emit,attr"`
						} `xml:"ifa-local"`
						IfaDestination struct {
							Text string `xml:",chardata"`
							Emit string `xml:"emit,attr"`
						} `xml:"ifa-destination"`
					} `xml:"interface-address"`
				} `xml:"address-family"`
			} `xml:"logical-interface"`
		} `xml:"physical-interface"`
	} `xml:"interface-information"`
}

type JuniperClient struct {
	SSHClient *devcon.SSHClient
}

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
