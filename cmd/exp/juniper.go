package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
	"github.com/montybeatnik/devcon/vendors/juniper"
)

func main() {
	khfp := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	jnprClient := juniper.NewJuniperClient(
		os.Getenv("SSH_USER"),
		"10.0.0.86",
		devcon.WithPassword(os.Getenv("SSH_PASSWORD")),
		devcon.WithTimeout(time.Second*1),
		devcon.WithHostKeyCallback(khfp),
	)
	bgpSummary, err := jnprClient.BGPSummary()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("DownPeers: %v\n", bgpSummary.BgpInformation.DownPeerCount)
	fmt.Printf("PeerCount: %v\n", bgpSummary.BgpInformation.PeerCount)
	for _, peer := range bgpSummary.BgpInformation.BgpPeer {
		fmt.Printf("Addr: %v, ASN: %v, State: %v\n", peer.PeerAddress, peer.PeerAs, peer.PeerState)
	}
}
