/*
Nsmf_PDUSession

SMF PDU Session Service. © 2021, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved. 

API version: 1.0.6
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapiclient

import (
	"encoding/json"
)

// checks if the Ncgi type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Ncgi{}

// Ncgi struct for Ncgi
type Ncgi struct {
	PlmnId PlmnId `json:"plmnId"`
	NrCellId string `json:"nrCellId"`
}

// NewNcgi instantiates a new Ncgi object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNcgi(plmnId PlmnId, nrCellId string) *Ncgi {
	this := Ncgi{}
	this.PlmnId = plmnId
	this.NrCellId = nrCellId
	return &this
}

// NewNcgiWithDefaults instantiates a new Ncgi object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNcgiWithDefaults() *Ncgi {
	this := Ncgi{}
	return &this
}

// GetPlmnId returns the PlmnId field value
func (o *Ncgi) GetPlmnId() PlmnId {
	if o == nil {
		var ret PlmnId
		return ret
	}

	return o.PlmnId
}

// GetPlmnIdOk returns a tuple with the PlmnId field value
// and a boolean to check if the value has been set.
func (o *Ncgi) GetPlmnIdOk() (*PlmnId, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PlmnId, true
}

// SetPlmnId sets field value
func (o *Ncgi) SetPlmnId(v PlmnId) {
	o.PlmnId = v
}

// GetNrCellId returns the NrCellId field value
func (o *Ncgi) GetNrCellId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NrCellId
}

// GetNrCellIdOk returns a tuple with the NrCellId field value
// and a boolean to check if the value has been set.
func (o *Ncgi) GetNrCellIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NrCellId, true
}

// SetNrCellId sets field value
func (o *Ncgi) SetNrCellId(v string) {
	o.NrCellId = v
}

func (o Ncgi) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Ncgi) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["plmnId"] = o.PlmnId
	toSerialize["nrCellId"] = o.NrCellId
	return toSerialize, nil
}

type NullableNcgi struct {
	value *Ncgi
	isSet bool
}

func (v NullableNcgi) Get() *Ncgi {
	return v.value
}

func (v *NullableNcgi) Set(val *Ncgi) {
	v.value = val
	v.isSet = true
}

func (v NullableNcgi) IsSet() bool {
	return v.isSet
}

func (v *NullableNcgi) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNcgi(val *Ncgi) *NullableNcgi {
	return &NullableNcgi{value: val, isSet: true}
}

func (v NullableNcgi) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNcgi) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


