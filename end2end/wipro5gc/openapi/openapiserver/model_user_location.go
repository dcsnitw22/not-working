/*
 * Nsmf_PDUSession
 *
 * SMF PDU Session Service. © 2021, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved.
 *
 * API version: 1.0.6
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapiserver

import "encoding/json"

type UserLocation struct {
	EutraLocation EutraLocation `json:"eutraLocation,omitempty"`

	NrLocation NrLocation `json:"nrLocation,omitempty"`

	N3gaLocation N3gaLocation `json:"n3gaLocation,omitempty"`
}

// UnmarshalJSON sets *m to a copy of data while respecting defaults if specified.
func (m *UserLocation) UnmarshalJSON(data []byte) error {

	type Alias UserLocation // To avoid infinite recursion
	return json.Unmarshal(data, (*Alias)(m))
}

// AssertUserLocationRequired checks if the required fields are not zero-ed
func AssertUserLocationRequired(obj UserLocation) error {
	if err := AssertEutraLocationRequired(obj.EutraLocation); err != nil {
		return err
	}
	if err := AssertNrLocationRequired(obj.NrLocation); err != nil {
		return err
	}
	if err := AssertN3gaLocationRequired(obj.N3gaLocation); err != nil {
		return err
	}
	return nil
}

// AssertUserLocationConstraints checks if the values respects the defined constraints
func AssertUserLocationConstraints(obj UserLocation) error {
	return nil
}
