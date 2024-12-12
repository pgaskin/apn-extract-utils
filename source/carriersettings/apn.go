package carriersettings

import (
	"fmt"
	"strconv"

	"github.com/pgaskin/apn-extract-utils/aosp/apn"
	"github.com/pgaskin/apn-extract-utils/aosp/carrier_list"
	"github.com/pgaskin/apn-extract-utils/aosp/carrier_settings"
)

// ConvertAPN converts src to an AOSP ApnSetting. It is as lenient as possible
// while still rejecting values it cannot represent. It does not include the
// carrier match attributes (mcc/mnc/carrier_id/mvno_type/mvno_match_data).
func ConvertAPN(src *carrier_settings.ApnItem) (apn.Setting, error) {
	s := apn.Empty()
	s.CarrierEnabled = true
	s.InfrastructureBitmask = 0 // clear it so it isn't set in the xml

	s.EntryName = src.GetName()
	s.APNName = src.GetValue()

	s.APNTypeBitmask = 0
	for _, t := range src.GetType() {
		switch t {
		case carrier_settings.ApnItem_ALL:
			s.APNTypeBitmask |= apn.TYPE_ALL
		case carrier_settings.ApnItem_DEFAULT:
			s.APNTypeBitmask |= apn.TYPE_DEFAULT
		case carrier_settings.ApnItem_MMS:
			s.APNTypeBitmask |= apn.TYPE_MMS
		case carrier_settings.ApnItem_SUPL:
			s.APNTypeBitmask |= apn.TYPE_SUPL
		case carrier_settings.ApnItem_DUN:
			s.APNTypeBitmask |= apn.TYPE_DUN
		case carrier_settings.ApnItem_HIPRI:
			s.APNTypeBitmask |= apn.TYPE_HIPRI
		case carrier_settings.ApnItem_FOTA:
			s.APNTypeBitmask |= apn.TYPE_FOTA
		case carrier_settings.ApnItem_IMS:
			s.APNTypeBitmask |= apn.TYPE_IMS
		case carrier_settings.ApnItem_CBS:
			s.APNTypeBitmask |= apn.TYPE_CBS
		case carrier_settings.ApnItem_IA:
			s.APNTypeBitmask |= apn.TYPE_IA
		case carrier_settings.ApnItem_EMERGENCY:
			s.APNTypeBitmask |= apn.TYPE_EMERGENCY
		case carrier_settings.ApnItem_XCAP:
			s.APNTypeBitmask |= apn.TYPE_XCAP
		case carrier_settings.ApnItem_UT:
			s.APNTypeBitmask |= apn.TYPE_XCAP
		case carrier_settings.ApnItem_RCS:
			s.APNTypeBitmask |= apn.TYPE_RCS
		default:
			return s, fmt.Errorf("unhandled apn type %#v", t)
		}
	}

	if v := src.GetBearerBitmask(); v != "0" {
		if err := s.BearerBitmask.UnmarshalText([]byte(v)); err != nil {
			return s, fmt.Errorf("parse bearer bitmask: %w", err)
		}
		s.NetworkTypeBitmask = apn.ConvertBearerBitmaskToNetworkTypeBitmask(s.BearerBitmask)
	}

	s.Server = src.GetServer()
	s.ProxyAddress = src.GetProxy()
	if v := src.GetPort(); v != "" {
		v, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return s, fmt.Errorf("parse proxy port: %w", err)
		}
		s.ProxyPort = int(v)
	}
	s.User = src.GetUser()
	s.Password = src.GetPassword()
	s.AuthType = apn.AuthType(src.GetAuthtype())

	s.MMSC = src.GetMmsc()
	s.MMSProxyAddress = src.GetMmscProxy()
	if v := src.GetMmscProxyPort(); v != "" {
		v, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return s, fmt.Errorf("parse mmsc proxy port: %w", err)
		}
		s.MMSProxyPort = int(v)
	}

	switch v := src.GetProtocol(); v {
	case carrier_settings.ApnItem_IP:
		s.Protocol = apn.PROTOCOL_IP
	case carrier_settings.ApnItem_IPV6:
		s.Protocol = apn.PROTOCOL_IPV6
	case carrier_settings.ApnItem_IPV4V6:
		s.Protocol = apn.PROTOCOL_IPV4V6
	case carrier_settings.ApnItem_PPP:
		s.Protocol = apn.PROTOCOL_PPP
	case 4: // TODO: update pb
		s.Protocol = apn.PROTOCOL_NON_IP
	default:
		return s, fmt.Errorf("unhandled protocol %#v", v)
	}

	switch v := src.GetRoamingProtocol(); v {
	case carrier_settings.ApnItem_IP:
		s.RoamingProtocol = apn.PROTOCOL_IP
	case carrier_settings.ApnItem_IPV6:
		s.RoamingProtocol = apn.PROTOCOL_IPV6
	case carrier_settings.ApnItem_IPV4V6:
		s.RoamingProtocol = apn.PROTOCOL_IPV4V6
	case carrier_settings.ApnItem_PPP:
		s.RoamingProtocol = apn.PROTOCOL_PPP
	case 4: // TODO: update pb
		s.Protocol = apn.PROTOCOL_NON_IP
	default:
		return s, fmt.Errorf("unhandled roaming protocol %#v", v)
	}

	if v := src.GetMtu(); v != 0 {
		// is this right?
		if s.Protocol != apn.PROTOCOL_IPV6 || s.RoamingProtocol != apn.PROTOCOL_IPV6 {
			s.MTUv4 = int(v)
		}
		if s.Protocol != apn.PROTOCOL_IP || s.RoamingProtocol != apn.PROTOCOL_IP {
			s.MTUv6 = int(v)
		}
	}

	if src.ProfileId != nil {
		s.ProfileID = int(*src.ProfileId)
	}
	s.MaxConns = int(src.GetMaxConns())
	s.WaitTime = int(src.GetWaitTime())
	s.MaxConnsTime = int(src.GetMaxConnsTime())
	s.Persistent = src.GetModemCognitive()
	s.APNSetID = int(src.GetApnSetId())

	switch v := src.GetSkip_464Xlat(); v {
	case carrier_settings.ApnItem_SKIP_464XLAT_DEFAULT:
		s.Skip464XLAT = apn.SKIP_464XLAT_DEFAULT
	case carrier_settings.ApnItem_SKIP_464XLAT_DISABLE:
		s.Skip464XLAT = apn.SKIP_464XLAT_DISABLE
	case carrier_settings.ApnItem_SKIP_464XLAT_ENABLE:
		s.Skip464XLAT = apn.SKIP_464XLAT_ENABLE
	default:
		return s, fmt.Errorf("unhandled skip 464xlat value %#v", v)
	}

	s.UserEditable = src.GetUserEditable()
	s.UserVisible = src.GetUserVisible()

	return s, nil
}

// WithAPNCarrier returns a copy of s with the carrier match attributes for an
// APN setting. It will clear them first, including the carrier id.
func WithAPNCarrier(s apn.Setting, attr *carrier_list.CarrierId) (apn.Setting, error) {
	s.CarrierID = 0
	if v := attr.GetMccMnc(); v != "" {
		s.OperatorNumeric = v
	} else {
		s.OperatorNumeric = "000000" // in aosp, this matches any carrier
	}
	if v := attr.GetMvnoData(); v != nil {
		switch v := v.(type) {
		case *carrier_list.CarrierId_Spn:
			s.MVNOType = apn.MVNO_TYPE_SPN
			s.MVNOMatchData = v.Spn
		case *carrier_list.CarrierId_Imsi:
			s.MVNOType = apn.MVNO_TYPE_IMSI
			s.MVNOMatchData = v.Imsi
		case *carrier_list.CarrierId_Gid1:
			s.MVNOType = apn.MVNO_TYPE_GID
			s.MVNOMatchData = v.Gid1
		default:
			return s, fmt.Errorf("unhandled mvno match type %#v", v)
		}
	} else {
		s.MVNOType = apn.MVNO_TYPE_UNKNOWN
		s.MVNOMatchData = ""
	}
	return s, nil
}
