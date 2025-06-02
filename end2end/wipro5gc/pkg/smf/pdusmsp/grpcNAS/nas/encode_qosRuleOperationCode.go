package nas

import "errors"

// Logic to Encode QoS Rule Operation Code
func EncodeQoSRuleOperationCode(qoc string) (byte, error) {
	switch qoc {
	case "Create new QoS rule":
		return CreateNewQoSRule, nil
	case "Delete existing QoS rule":
		return DeleteExistingQoSRule, nil
	case "Modify existing QoS rule and add packet filters":
		return ModifyExistingQoSRuleAndAddPacketFilters, nil
	case "Modify existing QoS rule and replace all packet filters":
		return ModifyExistingQoSRuleAndReplaceAllPacketFilters, nil
	case "Modify existing QoS rule and delete packet filters":
		return ModifyExistingQoSRuleAndDeletePacketFilters, nil
	case "Modify existing QoS rule without modifying packet filters":
		return ModifyExistingQoSRuleWithoutModifyingPacketFilters, nil
	}

	return 0, errors.New("invalid Input")

}
