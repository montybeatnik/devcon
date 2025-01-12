package juniper

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

type BGPSummary struct {
	BgpInformation struct {
		Text          string `xml:",chardata"`
		Xmlns         string `xml:"xmlns,attr"`
		GroupCount    string `xml:"group-count"`
		PeerCount     string `xml:"peer-count"`
		DownPeerCount string `xml:"down-peer-count"`
		BgpRib        struct {
			Text                          string `xml:",chardata"`
			Style                         string `xml:"style,attr"`
			Name                          string `xml:"name"`
			TotalPrefixCount              string `xml:"total-prefix-count"`
			ReceivedPrefixCount           string `xml:"received-prefix-count"`
			AcceptedPrefixCount           string `xml:"accepted-prefix-count"`
			ActivePrefixCount             string `xml:"active-prefix-count"`
			SuppressedPrefixCount         string `xml:"suppressed-prefix-count"`
			HistoryPrefixCount            string `xml:"history-prefix-count"`
			DampedPrefixCount             string `xml:"damped-prefix-count"`
			TotalExternalPrefixCount      string `xml:"total-external-prefix-count"`
			ActiveExternalPrefixCount     string `xml:"active-external-prefix-count"`
			AcceptedExternalPrefixCount   string `xml:"accepted-external-prefix-count"`
			SuppressedExternalPrefixCount string `xml:"suppressed-external-prefix-count"`
			TotalInternalPrefixCount      string `xml:"total-internal-prefix-count"`
			ActiveInternalPrefixCount     string `xml:"active-internal-prefix-count"`
			AcceptedInternalPrefixCount   string `xml:"accepted-internal-prefix-count"`
			SuppressedInternalPrefixCount string `xml:"suppressed-internal-prefix-count"`
			PendingPrefixCount            string `xml:"pending-prefix-count"`
			BgpRibState                   string `xml:"bgp-rib-state"`
		} `xml:"bgp-rib"`
		BgpPeer []struct {
			Text            string `xml:",chardata"`
			Style           string `xml:"style,attr"`
			Heading         string `xml:"heading,attr"`
			PeerAddress     string `xml:"peer-address"`
			PeerAs          string `xml:"peer-as"`
			InputMessages   string `xml:"input-messages"`
			OutputMessages  string `xml:"output-messages"`
			RouteQueueCount string `xml:"route-queue-count"`
			FlapCount       string `xml:"flap-count"`
			ElapsedTime     struct {
				Text    string `xml:",chardata"`
				Seconds string `xml:"seconds,attr"`
			} `xml:"elapsed-time"`
			PeerState struct {
				Text   string `xml:",chardata"`
				Format string `xml:"format,attr"`
			} `xml:"peer-state"`
			BgpRib struct {
				Text                  string `xml:",chardata"`
				Name                  string `xml:"name"`
				ActivePrefixCount     string `xml:"active-prefix-count"`
				ReceivedPrefixCount   string `xml:"received-prefix-count"`
				AcceptedPrefixCount   string `xml:"accepted-prefix-count"`
				SuppressedPrefixCount string `xml:"suppressed-prefix-count"`
			} `xml:"bgp-rib"`
		} `xml:"bgp-peer"`
	} `xml:"bgp-information"`
	Cli struct {
		Text   string `xml:",chardata"`
		Banner string `xml:"banner"`
	} `xml:"cli"`
}
