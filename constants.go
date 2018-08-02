package goftd

const (
	grantTypePassword string = "password"
	grantTypeCustom   string = "custom_token"

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
