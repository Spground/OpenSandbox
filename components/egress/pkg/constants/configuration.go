package constants

const (
	EnvBlockDoH443    = "OPENSANDBOX_EGRESS_BLOCK_DOH_443"
	EnvDoHBlocklist   = "OPENSANDBOX_EGRESS_DOH_BLOCKLIST" // comma-separated IP/CIDR
	EnvEgressMode     = "OPENSANDBOX_EGRESS_MODE"          // dns | dns+nft
	EnvEgressHTTPAddr = "OPENSANDBOX_EGRESS_HTTP_ADDR"
	EnvEgressToken    = "OPENSANDBOX_EGRESS_TOKEN"
	EnvEgressRules    = "OPENSANDBOX_EGRESS_RULES"
)

const (
	PolicyDnsOnly = "dns"
	PolicyDnsNft  = "dns+nft"
)

const (
	DefaultEgressServerAddr = ":18080"
)
