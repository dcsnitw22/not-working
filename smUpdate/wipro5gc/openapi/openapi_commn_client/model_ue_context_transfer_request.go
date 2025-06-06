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

// checks if the UEContextTransferRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UEContextTransferRequest{}

// UEContextTransferRequest struct for UEContextTransferRequest
type UEContextTransferRequest struct {
	JsonData *UeContextTransferReqData `json:"jsonData,omitempty"`
	BinaryDataN1Message **os.File `json:"binaryDataN1Message,omitempty"`
}

// NewUEContextTransferRequest instantiates a new UEContextTransferRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUEContextTransferRequest() *UEContextTransferRequest {
	this := UEContextTransferRequest{}
	return &this
}

// NewUEContextTransferRequestWithDefaults instantiates a new UEContextTransferRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUEContextTransferRequestWithDefaults() *UEContextTransferRequest {
	this := UEContextTransferRequest{}
	return &this
}

// GetJsonData returns the JsonData field value if set, zero value otherwise.
func (o *UEContextTransferRequest) GetJsonData() UeContextTransferReqData {
	if o == nil || IsNil(o.JsonData) {
		var ret UeContextTransferReqData
		return ret
	}
	return *o.JsonData
}

// GetJsonDataOk returns a tuple with the JsonData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UEContextTransferRequest) GetJsonDataOk() (*UeContextTransferReqData, bool) {
	if o == nil || IsNil(o.JsonData) {
		return nil, false
	}
	return o.JsonData, true
}

// HasJsonData returns a boolean if a field has been set.
func (o *UEContextTransferRequest) HasJsonData() bool {
	if o != nil && !IsNil(o.JsonData) {
		return true
	}

	return false
}

// SetJsonData gets a reference to the given UeContextTransferReqData and assigns it to the JsonData field.
func (o *UEContextTransferRequest) SetJsonData(v UeContextTransferReqData) {
	o.JsonData = &v
}

// GetBinaryDataN1Message returns the BinaryDataN1Message field value if set, zero value otherwise.
func (o *UEContextTransferRequest) GetBinaryDataN1Message() *os.File {
	if o == nil || IsNil(o.BinaryDataN1Message) {
		var ret *os.File
		return ret
	}
	return *o.BinaryDataN1Message
}

// GetBinaryDataN1MessageOk returns a tuple with the BinaryDataN1Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UEContextTransferRequest) GetBinaryDataN1MessageOk() (**os.File, bool) {
	if o == nil || IsNil(o.BinaryDataN1Message) {
		return nil, false
	}
	return o.BinaryDataN1Message, true
}

// HasBinaryDataN1Message returns a boolean if a field has been set.
func (o *UEContextTransferRequest) HasBinaryDataN1Message() bool {
	if o != nil && !IsNil(o.BinaryDataN1Message) {
		return true
	}

	return false
}

// SetBinaryDataN1Message gets a reference to the given *os.File and assigns it to the BinaryDataN1Message field.
func (o *UEContextTransferRequest) SetBinaryDataN1Message(v *os.File) {
	o.BinaryDataN1Message = &v
}

func (o UEContextTransferRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UEContextTransferRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.JsonData) {
		toSerialize["jsonData"] = o.JsonData
	}
	if !IsNil(o.BinaryDataN1Message) {
		toSerialize["binaryDataN1Message"] = o.BinaryDataN1Message
	}
	return toSerialize, nil
}

type NullableUEContextTransferRequest struct {
	value *UEContextTransferRequest
	isSet bool
}

func (v NullableUEContextTransferRequest) Get() *UEContextTransferRequest {
	return v.value
}

func (v *NullableUEContextTransferRequest) Set(val *UEContextTransferRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableUEContextTransferRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableUEContextTransferRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUEContextTransferRequest(val *UEContextTransferRequest) *NullableUEContextTransferRequest {
	return &NullableUEContextTransferRequest{value: val, isSet: true}
}

func (v NullableUEContextTransferRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUEContextTransferRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


