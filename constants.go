package goftd

const (
	grantTypePassword        string = "password"
	grantTypeCustom          string = "custom_token"
	duplicateActionDoNothing int    = 1
	duplicateActionReplace   int    = 2

	//LogActionNone LOG_NONE
	LogActionNone string = "LOG_NONE"

	//LogActionFlowStart LOG_FLOW_START
	LogActionFlowStart string = "LOG_FLOW_START"

	//RuleActionPermit PERMIT
	RuleActionPermit string = "PERMIT"
)
