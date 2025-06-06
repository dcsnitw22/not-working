/*
Namf_Communication

AMF Communication Service © 2022, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved. 

API version: 1.0.8
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi_commn_client

import (
	"encoding/json"
	"os"
)

// checks if the NonUeN2InfoNotifyRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NonUeN2InfoNotifyRequest{}

// NonUeN2InfoNotifyRequest struct for NonUeN2InfoNotifyRequest
type NonUeN2InfoNotifyRequest struct {
	JsonData *N2InformationNotification `json:"jsonData,omitempty"`
	BinaryDataN2Information **os.File `json:"binaryDataN2Information,omitempty"`
}

// NewNonUeN2InfoNotifyRequest instantiates a new NonUeN2InfoNotifyRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNonUeN2InfoNotifyRequest() *NonUeN2InfoNotifyRequest {
	this := NonUeN2InfoNotifyRequest{}
	return &this
}

// NewNonUeN2InfoNotifyRequestWithDefaults instantiates a new NonUeN2InfoNotifyRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNonUeN2InfoNotifyRequestWithDefaults() *NonUeN2InfoNotifyRequest {
	this := NonUeN2InfoNotifyRequest{}
	return &this
}

// GetJsonData returns the JsonData field value if set, zero value otherwise.
func (o *NonUeN2InfoNotifyRequest) GetJsonData() N2InformationNotification {
	if o == nil || IsNil(o.JsonData) {
		var ret N2InformationNotification
		return ret
	}
	return *o.JsonData
}

// GetJsonDataOk returns a tuple with the JsonData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NonUeN2InfoNotifyRequest) GetJsonDataOk() (*N2InformationNotification, bool) {
	if o == nil || IsNil(o.JsonData) {
		return nil, false
	}
	return o.JsonData, true
}

// HasJsonData returns a boolean if a field has been set.
func (o *NonUeN2InfoNotifyRequest) HasJsonData() bool {
	if o != nil && !IsNil(o.JsonData) {
		return true
	}

	return false
}

// SetJsonData gets a reference to the given N2InformationNotification and assigns it to the JsonData field.
func (o *NonUeN2InfoNotifyRequest) SetJsonData(v N2InformationNotification) {
	o.JsonData = &v
}

// GetBinaryDataN2Information returns the BinaryDataN2Information field value if set, zero value otherwise.
func (o *NonUeN2InfoNotifyRequest) GetBinaryDataN2Information() *os.File {
	if o == nil || IsNil(o.BinaryDataN2Information) {
		var ret *os.File
		return ret
	}
	return *o.BinaryDataN2Information
}

// GetBinaryDataN2InformationOk returns a tuple with the BinaryDataN2Information field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NonUeN2InfoNotifyRequest) GetBinaryDataN2InformationOk() (**os.File, bool) {
	if o == nil || IsNil(o.BinaryDataN2Information) {
		return nil, false
	}
	return o.BinaryDataN2Information, true
}

// HasBinaryDataN2Information returns a boolean if a field has been set.
func (o *NonUeN2InfoNotifyRequest) HasBinaryDataN2Information() bool {
	if o != nil && !IsNil(o.BinaryDataN2Information) {
		return true
	}

	return false
}

// SetBinaryDataN2Information gets a reference to the given *os.File and assigns it to the BinaryDataN2Information field.
func (o *NonUeN2InfoNotifyRequest) SetBinaryDataN2Information(v *os.File) {
	o.BinaryDataN2Information = &v
}

func (o NonUeN2InfoNotifyRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NonUeN2InfoNotifyRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.JsonData) {
		toSerialize["jsonData"] = o.JsonData
	}
	if !IsNil(o.BinaryDataN2Information) {
		toSerialize["binaryDataN2Information"] = o.BinaryDataN2Information
	}
	return toSerialize, nil
}

type NullableNonUeN2InfoNotifyRequest struct {
	value *NonUeN2InfoNotifyRequest
	isSet bool
}

func (v NullableNonUeN2InfoNotifyRequest) Get() *NonUeN2InfoNotifyRequest {
	return v.value
}

func (v *NullableNonUeN2InfoNotifyRequest) Set(val *NonUeN2InfoNotifyRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableNonUeN2InfoNotifyRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableNonUeN2InfoNotifyRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNonUeN2InfoNotifyRequest(val *NonUeN2InfoNotifyRequest) *NullableNonUeN2InfoNotifyRequest {
	return &NullableNonUeN2InfoNotifyRequest{value: val, isSet: true}
}

func (v NullableNonUeN2InfoNotifyRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNonUeN2InfoNotifyRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


