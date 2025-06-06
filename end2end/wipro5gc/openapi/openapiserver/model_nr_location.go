/*
 * Nsmf_PDUSession
 *
 * SMF PDU Session Service. © 2021, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved.
 *
 * API version: 1.0.6
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapiserver

import (
	"encoding/json"
	"errors"
	"time"
)

type NrLocation struct {
	Tai Tai `json:"tai"`

	Ncgi Ncgi `json:"ncgi"`

	AgeOfLocationInformation int32 `json:"ageOfLocationInformation,omitempty"`

	UeLocationTimestamp time.Time `json:"ueLocationTimestamp,omitempty"`

	GeographicalInformation string `json:"geographicalInformation,omitempty"`

	GeodeticInformation string `json:"geodeticInformation,omitempty"`

	GlobalGnbId *GlobalRanNodeId `json:"globalGnbId,omitempty"`
}

// UnmarshalJSON sets *m to a copy of data while respecting defaults if specified.
func (m *NrLocation) UnmarshalJSON(data []byte) error {

	type Alias NrLocation // To avoid infinite recursion
	return json.Unmarshal(data, (*Alias)(m))
}

// AssertNrLocationRequired checks if the required fields are not zero-ed
func AssertNrLocationRequired(obj NrLocation) error {
	elements := map[string]interface{}{
		"tai":  obj.Tai,
		"ncgi": obj.Ncgi,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertTaiRequired(obj.Tai); err != nil {
		return err
	}
	if err := AssertNcgiRequired(obj.Ncgi); err != nil {
		return err
	}
	if obj.GlobalGnbId != nil {
		if err := AssertGlobalRanNodeIdRequired(*obj.GlobalGnbId); err != nil {
			return err
		}
	}
	return nil
}

// AssertNrLocationConstraints checks if the values respects the defined constraints
func AssertNrLocationConstraints(obj NrLocation) error {
	if obj.AgeOfLocationInformation < 0 {
		return &ParsingError{Err: errors.New(errMsgMinValueConstraint)}
	}
	if obj.AgeOfLocationInformation > 32767 {
		return &ParsingError{Err: errors.New(errMsgMaxValueConstraint)}
	}
	return nil
}
