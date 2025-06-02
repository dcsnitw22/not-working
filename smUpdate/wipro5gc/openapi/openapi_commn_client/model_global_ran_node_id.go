package openapi_commn_client

type GlobalRanNodeId struct {
	PlmnId  *PlmnId `json:"plmnId" yaml:"plmnId" bson:"plmnId"`
	N3IwfId string  `json:"n3IwfId,omitempty" yaml:"n3IwfId" bson:"n3IwfId"`
	GNbId   *GNbId  `json:"gNbId,omitempty" yaml:"gNbId" bson:"gNbId"`
	NgeNbId string  `json:"ngeNbId,omitempty" yaml:"ngeNbId" bson:"ngeNbId"`
}
