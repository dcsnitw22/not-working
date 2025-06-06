/*
Namf_Communication

AMF Communication Service © 2022, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved.

API version: 1.0.8
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi_commn_client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	// "os"
	"strings"
	"fmt"
)

// N1N2MessageCollectionDocumentAPIService N1N2MessageCollectionDocumentAPI service
type N1N2MessageCollectionDocumentAPIService service

type ApiN1N2MessageTransferRequest struct {
	ctx context.Context
	ApiService *N1N2MessageCollectionDocumentAPIService
	ueContextId string
	n1N2MessageTransferReqData *N1N2MessageTransferReqData
	binaryDataN1MessageContent []byte
	binaryDataN2InfoContent []byte
}

func (r ApiN1N2MessageTransferRequest) N1N2MessageTransferReqData(n1N2MessageTransferReqData N1N2MessageTransferReqData) ApiN1N2MessageTransferRequest {
	r.n1N2MessageTransferReqData = &n1N2MessageTransferReqData
	return r
}

func (r ApiN1N2MessageTransferRequest) BinaryDataN1MessageContent(binaryDataN1MessageContent []byte) ApiN1N2MessageTransferRequest {
	r.binaryDataN1MessageContent = binaryDataN1MessageContent
	return r
}

func (r ApiN1N2MessageTransferRequest) BinaryDataN2InfoContent(binaryDataN2InfoContent []byte) ApiN1N2MessageTransferRequest {
	r.binaryDataN2InfoContent = binaryDataN2InfoContent
	return r
}

func (r ApiN1N2MessageTransferRequest) Execute() (*N1N2MessageTransferRspData, *http.Response, error) {
	return r.ApiService.N1N2MessageTransferExecute(r)
}

/*
N1N2MessageTransfer Namf_Communication N1N2 Message Transfer (UE Specific) service Operation

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ueContextId UE Context Identifier
 @return ApiN1N2MessageTransferRequest
*/
func (a *N1N2MessageCollectionDocumentAPIService) N1N2MessageTransfer(ctx context.Context, ueContextId string) ApiN1N2MessageTransferRequest {
	return ApiN1N2MessageTransferRequest{
		ApiService: a,
		ctx: ctx,
		ueContextId: ueContextId,
	}
}

// Execute executes the request
//  @return N1N2MessageTransferRspData
func (a *N1N2MessageCollectionDocumentAPIService) N1N2MessageTransferExecute(r ApiN1N2MessageTransferRequest) (*N1N2MessageTransferRspData, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *N1N2MessageTransferRspData
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "N1N2MessageCollectionDocumentAPIService.N1N2MessageTransfer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/ue-contexts/{ueContextId}/n1-n2-messages"
	localVarPath = strings.Replace(localVarPath, "{"+"ueContextId"+"}", url.PathEscape(parameterValueToString(r.ueContextId, "ueContextId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	// if r.n1N2MessageTransferReqData == nil {
	// 	return localVarReturnValue, nil, reportError("n1N2MessageTransferReqData is required and must be specified")
	// }

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json", "multipart/related"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json", "application/problem+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	//localVarPostBody = r.n1N2MessageTransferReqData

	if r.n1N2MessageTransferReqData!=nil{
		paramJson, err := parameterToJson(*r.n1N2MessageTransferReqData)
		if err != nil{
			return localVarReturnValue, nil, err
		}
		localVarFormParams.Add("jsonData", paramJson)
	}

	var binaryDataN1MessageContentLocalVarFormFileName string
	var binaryDataN1MessageContentLocalVarFileName     string
	var binaryDataN1MessageContentLocalVarFileBytes    []byte

	binaryDataN1MessageContentLocalVarFormFileName = "binaryDataN1MessageContent"

	if (r.binaryDataN1MessageContent!=nil){
		binaryDataN1MessageContentLocalVarFileBytes = r.binaryDataN1MessageContent
		binaryDataN1MessageContentLocalVarFileName = "binaryDataN1MessageContent"
		formFiles = append(formFiles, formFile{fileBytes: binaryDataN1MessageContentLocalVarFileBytes, fileName: binaryDataN1MessageContentLocalVarFileName, formFileName: binaryDataN1MessageContentLocalVarFormFileName})
		fmt.Println("Form files :", formFiles)
	}
	
	/*
	binaryDataN1MessageContentLocalVarFile := r.binaryDataN1MessageContent
	fmt.Println("binaryDataN1MessageContentLocalVarFile: ", binaryDataN1MessageContentLocalVarFile)
	if binaryDataN1MessageContentLocalVarFile != nil {
		fbs, _ := io.ReadAll(binaryDataN1MessageContentLocalVarFile)

		binaryDataN1MessageContentLocalVarFileBytes = fbs
		binaryDataN1MessageContentLocalVarFileName = binaryDataN1MessageContentLocalVarFile.Name()
		binaryDataN1MessageContentLocalVarFile.Close()
		formFiles = append(formFiles, formFile{fileBytes: binaryDataN1MessageContentLocalVarFileBytes, fileName: binaryDataN1MessageContentLocalVarFileName, formFileName: binaryDataN1MessageContentLocalVarFormFileName})
		fmt.Println("Form files :", formFiles)
	}
	*/

	var binaryDataN2InfoContentLocalVarFormFileName string
	var binaryDataN2InfoContentLocalVarFileName     string
	var binaryDataN2InfoContentLocalVarFileBytes    []byte

	binaryDataN2InfoContentLocalVarFormFileName = "binaryDataN2InfoContent"

	if r.binaryDataN2InfoContent!=nil{
		binaryDataN2InfoContentLocalVarFileBytes = r.binaryDataN2InfoContent
		binaryDataN2InfoContentLocalVarFileName = "binaryDataN2InfoContent"
		formFiles = append(formFiles, formFile{fileBytes: binaryDataN2InfoContentLocalVarFileBytes, fileName: binaryDataN2InfoContentLocalVarFileName, formFileName: binaryDataN2InfoContentLocalVarFormFileName})
		fmt.Println("Form files :", formFiles)
	}
	
	/*
	binaryDataN2InfoContentLocalVarFile := r.binaryDataN2InfoContent
	fmt.Println("binaryDataN2InfoContentLocalVarFile: ", binaryDataN2InfoContentLocalVarFile)
	if binaryDataN2InfoContentLocalVarFile != nil {
		fbs, _ := io.ReadAll(binaryDataN2InfoContentLocalVarFile)

		binaryDataN2InfoContentLocalVarFileBytes = fbs
		binaryDataN2InfoContentLocalVarFileName = binaryDataN2InfoContentLocalVarFile.Name()
		binaryDataN2InfoContentLocalVarFile.Close()
		formFiles = append(formFiles, formFile{fileBytes: binaryDataN2InfoContentLocalVarFileBytes, fileName: binaryDataN2InfoContentLocalVarFileName, formFileName: binaryDataN2InfoContentLocalVarFormFileName})
		fmt.Println("Form files :", formFiles)
	}
	*/

	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 307 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 409 {
			var v N1N2MessageTransferError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 411 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 413 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 415 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 503 {
			var v ProblemDetails
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 504 {
			var v N1N2MessageTransferError
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
