package context

import (
	"time"

	"github.com/sirupsen/logrus"
	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

type RelAction int

const (
	RanUeNgapIdUnspecified int64 = 0xffffffff
)

const (
	UeContextN2NormalRelease RelAction = iota
	UeContextReleaseHandover
	UeContextReleaseUeContext
)

type RanUe struct {
	/* UE identity*/
	RanUeNgapId int64
	AmfUeNgapId int64

	/* HandOver Info*/
	/*HandOverType        ngapType.HandoverType
	SuccessPduSessionId []int32
	SourceUe            *RanUe
	TargetUe           *RanUe */

	/* UserLocation*/
	Tai      openapi_commn_client.Tai
	Location openapi_commn_client.UserLocation
	/* context about udm */
	SupportVoPSn3gpp  bool
	SupportVoPS       bool
	SupportedFeatures string
	LastActTime       *time.Time

	/* Related Context*/
	AmfUe        *AmfUe
	Ran          *AmfRan
	HoldingAmfUe *AmfUe // The AmfUe that is already exist (CM-Idle, Re-Registration)

	/* Routing ID */
	RoutingID string
	/* Trace Recording Session Reference */
	Trsr string
	/* Ue Context Release Action */
	ReleaseAction RelAction
	/* context used for AMF Re-allocation procedure */
	OldAmfName            string
	InitialUEMessage      []byte
	RRCEstablishmentCause string // Received from initial ue message; pattern: ^[0-9a-fA-F]+$
	UeContextRequest      bool   // Receive UEContextRequest IE from RAN

	/* send initial context setup request or not*/
	InitialContextSetup bool

	/* logger */
	Log *logrus.Entry
}
