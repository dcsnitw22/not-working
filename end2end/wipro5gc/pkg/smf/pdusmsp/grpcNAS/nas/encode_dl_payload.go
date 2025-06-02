package nas

import (
	"errors"
)

func EncodeDLPayload(payloadtype string, payload interface{}) ([]byte, error) {
	if payloadtype == "N1 SM information" {
		switch payload.(type) {
		case PduSessionEstablishmentAccept:
			pl := payload.(PduSessionEstablishmentAccept)
			result, err := ReRouteEncode(pl.ExtendedProtocolDiscriminator, pl.MessageType, pl)
			if err != nil {
				return nil, err
			}
			return result, nil

		case PduSessionEstablishmentReject:
			pl := payload.(PduSessionEstablishmentReject)
			result, err := ReRouteEncode(pl.ExtendedProtocolDiscriminator, pl.MessageType, pl)
			if err != nil {
				return nil, err
			}
			return result, nil

		case PduSessionModificationReject:
			pl := payload.(PduSessionModificationReject)
			result, err := ReRouteEncode(pl.ExtendedProtocolDiscriminator, pl.MessageType, pl)
			if err != nil {
				return nil, err
			}
			return result, nil

		case PduSessionReleaseReject:
			pl := payload.(PduSessionReleaseReject)
			result, err := ReRouteEncode(pl.ExtendedProtocolDiscriminator, pl.MessageType, pl)
			if err != nil {
				return nil, err
			}
			return result, nil

		}

	}

	return nil, errors.New("currently not handling this container type")

}
