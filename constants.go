package goftd

const (
	grantTypePassword string = "password"
	grantTypeCustom   string = "custom_token"

	apiPOST   string = "POST"
	apiPUT    string = "PUT"
	apiDELETE string = "DELETE"
	apiGET    string = "GET"

	apiBasePath                 string = "api/fdm/v1/"
	apiTokenEndpoint            string = "fdm/token"
	apiNetworksEndpoint         string = "object/networks"
	apiNetworkGroupsEndpoint    string = "object/networkgroups"
	apiTCPPortObjectsEndpoint   string = "object/tcpports"
	apiUDPPortObjectsEndpoint   string = "object/udpports"
	apiPortObjectGroupsEndpoint string = "object/portgroups"

	// TypeUDPPortObject object type udp port
	TypeUDPPortObject string = "udpportobject"
	// TypeTCPPortObject object type tcp port
	TypeTCPPortObject string = "tcpportobject"

	//DuplicateActionError Error on duplicate
	DuplicateActionError int = 0

	//DuplicateActionDoNothing Don't do anything
	DuplicateActionDoNothing int = 1

	//DuplicateActionReplace Replace
	DuplicateActionReplace int = 2

	//LogActionNone LOG_NONE
	LogActionNone string = "LOG_NONE"

	//LogActionFlowStart LOG_FLOW_START
	LogActionFlowStart string = "LOG_FLOW_START"

	//RuleActionPermit PERMIT
	RuleActionPermit string = "PERMIT"
)
