package context

import (
	"sync"

	"w5gc.io/wipro5gcore/openapi/openapi_commn_client"
)

type IDGenerator struct {
	lock     sync.Mutex
	minValue int64
	maxValue int64
	valRange int64
	offset   int64
	usedMap  map[int64]bool
}

func NewGenerator(minValue, maxValue int64) *IDGenerator {
	idGenerator := &IDGenerator{}
	idGenerator.init(minValue, maxValue)
	return idGenerator
}

func (x *IDGenerator) init(minVal int64, maxVal int64) {
	x.minValue = minVal
	x.maxValue = maxVal
	x.valRange = maxVal - minVal + 1
	x.offset = 0
	x.usedMap = make(map[int64]bool)
}

type AMFContext struct {
	EventSubscriptionIDGenerator *IDGenerator
	EventSubscriptions           sync.Map
	UePool                       sync.Map // map[supi]*AmfUe
	RanUePool                    sync.Map // map[AmfUeNgapID]*RanUe
	AmfRanPool                   sync.Map // map[net.Conn]*AmfRan
	// LadnPool                     map[string]factory.Ladn // dnn as key
	SupportTaiLists []openapi_commn_client.Tai
	ServedGuamiList []openapi_commn_client.Guami
	// PlmnSupportList        []factory.PlmnSupportItem
	RelativeCapacity int64
	NfId             string
	Name             string
	// NfService              map[models.ServiceName]models.NfService // nfservice that amf support
	// UriScheme              models.UriScheme
	BindingIPv4            string
	SBIPort                int
	RegisterIPv4           string
	HttpIPv6Address        string
	TNLWeightFactor        int64
	SupportDnnLists        []string
	AMFStatusSubscriptions sync.Map // map[subscriptionID]models.SubscriptionData
	NrfUri                 string
	NrfCertPem             string
	// SecurityAlgorithm            SecurityAlgorithm
	// NetworkName                  factory.NetworkName
	NgapIpList             []string // NGAP Server IP
	NgapPort               int
	T3502Value             int    // unit is second
	T3512Value             int    // unit is second
	Non3gppDeregTimerValue int    // unit is second
	TimeZone               string // "[+-]HH:MM[+][1-2]", Refer to TS 29.571 - 5.2.2 Simple Data Types
	// read-only fields
	/*T3513Cfg factory.TimerValue
	T3522Cfg factory.TimerValue
	T3550Cfg factory.TimerValue
	T3560Cfg factory.TimerValue
	T3565Cfg factory.TimerValue
	T3570Cfg factory.TimerValue
	T3555Cfg factory.TimerValue*/
	Locality string

	OAuth2Required bool
}
