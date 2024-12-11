package apn

import (
	"bytes"
	"fmt"
	"iter"
	"math/bits"
	"net/url"
	"slices"
	"strconv"
)

// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/telephony/java/android/telephony/data/ApnSetting.java;drc=4ba139804a0a420c376d8fffbcb7e9f2fa3f65a8
// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/telephony/java/android/telephony/TelephonyManager.java;l=14766;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae

// TODO: use codegen for generating the enum/text-bitmask/numeric-bitmask methods

type NetworkType int

const (
	NETWORK_TYPE_UNKNOWN  NetworkType = 0
	NETWORK_TYPE_GPRS     NetworkType = 1
	NETWORK_TYPE_EDGE     NetworkType = 2
	NETWORK_TYPE_UMTS     NetworkType = 3
	NETWORK_TYPE_CDMA     NetworkType = 4
	NETWORK_TYPE_EVDO_0   NetworkType = 5
	NETWORK_TYPE_EVDO_A   NetworkType = 6
	NETWORK_TYPE_1xRTT    NetworkType = 7
	NETWORK_TYPE_HSDPA    NetworkType = 8
	NETWORK_TYPE_HSUPA    NetworkType = 9
	NETWORK_TYPE_HSPA     NetworkType = 10
	NETWORK_TYPE_IDEN     NetworkType = 11 // deprecated, no longer used starting in Android U
	NETWORK_TYPE_EVDO_B   NetworkType = 12
	NETWORK_TYPE_LTE      NetworkType = 13
	NETWORK_TYPE_EHRPD    NetworkType = 14
	NETWORK_TYPE_HSPAP    NetworkType = 15
	NETWORK_TYPE_GSM      NetworkType = 16
	NETWORK_TYPE_TD_SCDMA NetworkType = 17
	NETWORK_TYPE_IWLAN    NetworkType = 18
	NETWORK_TYPE_LTE_CA   NetworkType = 19 // hidden since Android R (b/170729553 -- tldr: should be 13)
	NETWORK_TYPE_NR       NetworkType = 20 // only used for 5G SA; NSA is LTE
)

type NetworkTypeBitmask int

const (
	NETWORK_TYPE_BITMASK_UNKNOWN  NetworkTypeBitmask = 0
	NETWORK_TYPE_BITMASK_GPRS     NetworkTypeBitmask = (1 << (NETWORK_TYPE_GPRS - 1))
	NETWORK_TYPE_BITMASK_EDGE     NetworkTypeBitmask = (1 << (NETWORK_TYPE_EDGE - 1))
	NETWORK_TYPE_BITMASK_UMTS     NetworkTypeBitmask = (1 << (NETWORK_TYPE_UMTS - 1))
	NETWORK_TYPE_BITMASK_CDMA     NetworkTypeBitmask = (1 << (NETWORK_TYPE_CDMA - 1))
	NETWORK_TYPE_BITMASK_EVDO_0   NetworkTypeBitmask = (1 << (NETWORK_TYPE_EVDO_0 - 1))
	NETWORK_TYPE_BITMASK_EVDO_A   NetworkTypeBitmask = (1 << (NETWORK_TYPE_EVDO_A - 1))
	NETWORK_TYPE_BITMASK_1xRTT    NetworkTypeBitmask = (1 << (NETWORK_TYPE_1xRTT - 1))
	NETWORK_TYPE_BITMASK_HSDPA    NetworkTypeBitmask = (1 << (NETWORK_TYPE_HSDPA - 1))
	NETWORK_TYPE_BITMASK_HSUPA    NetworkTypeBitmask = (1 << (NETWORK_TYPE_HSUPA - 1))
	NETWORK_TYPE_BITMASK_HSPA     NetworkTypeBitmask = (1 << (NETWORK_TYPE_HSPA - 1))
	NETWORK_TYPE_BITMASK_IDEN     NetworkTypeBitmask = (1 << (NETWORK_TYPE_IDEN - 1))
	NETWORK_TYPE_BITMASK_EVDO_B   NetworkTypeBitmask = (1 << (NETWORK_TYPE_EVDO_B - 1))
	NETWORK_TYPE_BITMASK_LTE      NetworkTypeBitmask = (1 << (NETWORK_TYPE_LTE - 1))
	NETWORK_TYPE_BITMASK_EHRPD    NetworkTypeBitmask = (1 << (NETWORK_TYPE_EHRPD - 1))
	NETWORK_TYPE_BITMASK_HSPAP    NetworkTypeBitmask = (1 << (NETWORK_TYPE_HSPAP - 1))
	NETWORK_TYPE_BITMASK_GSM      NetworkTypeBitmask = (1 << (NETWORK_TYPE_GSM - 1))
	NETWORK_TYPE_BITMASK_TD_SCDMA NetworkTypeBitmask = (1 << (NETWORK_TYPE_TD_SCDMA - 1))
	NETWORK_TYPE_BITMASK_IWLAN    NetworkTypeBitmask = (1 << (NETWORK_TYPE_IWLAN - 1))
	NETWORK_TYPE_BITMASK_LTE_CA   NetworkTypeBitmask = (1 << (NETWORK_TYPE_LTE_CA - 1))
	NETWORK_TYPE_BITMASK_NR       NetworkTypeBitmask = (1 << (NETWORK_TYPE_NR - 1))
)

type Type int

const (
	TYPE_NONE       Type = 0               // none (should only be used for initialization)
	TYPE_DEFAULT    Type = 1 << (iota - 1) // default data traffic
	TYPE_MMS                               // MMS (Multimedia Messaging Service) traffic
	TYPE_SUPL                              // SUPL (Secure User Plane Location) assisted GPS
	TYPE_DUN                               // DUN (Dial-up networking) traffic
	TYPE_HIPRI                             // high-priority traffic
	TYPE_FOTA                              // FOTA (Firmware over-the-air) traffic (accessing the carrier's FOTA portal, used for over the air updates)
	TYPE_IMS                               // IMS (IP Multimedia Subsystem) traffic
	TYPE_CBS                               // CBS (Carrier Branded Services) traffic
	TYPE_IA                                // IA (Initial Attach) APN
	TYPE_EMERGENCY                         // emergency PDN (this is not an IA apn, but is used for access to carrier services in an emergency call situation)
	TYPE_MCX                               // MCX (Mission Critical Service) where X can be PTT/Video/Data
	TYPE_XCAP                              // XCAP (XML Configuration Access Protocol) traffic
	TYPE_VSIM                              // Virtual SIM service
	TYPE_BIP                               // Bearer Independent Protocol
	TYPE_ENTERPRISE                        // ENTERPRISE traffic
	TYPE_RCS                               // RCS (Rich Communication Services)
	type_limit
	// all data connections
	TYPE_ALL = TYPE_DEFAULT | TYPE_HIPRI | TYPE_MMS | TYPE_SUPL | TYPE_DUN | TYPE_FOTA | TYPE_IMS | TYPE_CBS
)

const (
	TYPE_ALL_STRING        = "*"
	TYPE_DEFAULT_STRING    = "default"
	TYPE_MMS_STRING        = "mms"
	TYPE_SUPL_STRING       = "supl"
	TYPE_DUN_STRING        = "dun"
	TYPE_HIPRI_STRING      = "hipri"
	TYPE_FOTA_STRING       = "fota"
	TYPE_IMS_STRING        = "ims"
	TYPE_CBS_STRING        = "cbs"
	TYPE_IA_STRING         = "ia"
	TYPE_EMERGENCY_STRING  = "emergency"
	TYPE_MCX_STRING        = "mcx"
	TYPE_XCAP_STRING       = "xcap"
	TYPE_VSIM_STRING       = "vsim"
	TYPE_BIP_STRING        = "bip"
	TYPE_ENTERPRISE_STRING = "enterprise"
	TYPE_RCS_STRING        = "rcs"
)

type AuthType int

const (
	AUTH_TYPE_UNKNOWN     AuthType = iota - 1 // unknown
	AUTH_TYPE_PAP                             // PAP
	AUTH_TYPE_CHAP                            // CHAP
	AUTH_TYPE_PAP_OR_CHAP                     // PAP or CHAP
)

type Skip464XLAT int

const (
	SKIP_464XLAT_DEFAULT Skip464XLAT = iota - 1
	SKIP_464XLAT_DISABLE
	SKIP_464XLAT_ENABLE
)

type Protocol int

const (
	PROTOCOL_UNKNOWN      Protocol = iota - 1 // unknown
	PROTOCOL_IP                               // internet protocol
	PROTOCOL_IPV6                             // internet protocol, version 6
	PROTOCOL_IPV4V6                           // virtual PDP type introduced to handle dual IP stack UE capability
	PROTOCOL_PPP                              // point to point protocol
	PROTOCOL_NON_IP                           // transfer of Non-IP data to external packet data network
	PROTOCOL_UNSTRUCTURED                     // transfer of Unstructured data to the Data Network via N6
)

type MVNOType int

const (
	MVNO_TYPE_UNKNOWN MVNOType = iota - 1 // unset
	MVNO_TYPE_SPN                         // service provider name
	MVNO_TYPE_IMSI                        // IMSI
	MVNO_TYPE_GID                         // group identifier level 1
	MVNO_TYPE_ICCID                       // ICCID
)

type Infrastructure int

const (
	INFRASTRUCTURE_CELLULAR Infrastructure = 1 << iota
	INFRASTRUCTURE_SATELLITE
)

func (x NetworkType) Valid() bool {
	return x.String() != ""
}

func (x NetworkType) String() string {
	switch x {
	case NETWORK_TYPE_GPRS:
		return "GPRS"
	case NETWORK_TYPE_EDGE:
		return "EDGE"
	case NETWORK_TYPE_UMTS:
		return "UMTS"
	case NETWORK_TYPE_CDMA:
		return "CDMA"
	case NETWORK_TYPE_EVDO_0:
		return "EVDO_0"
	case NETWORK_TYPE_EVDO_A:
		return "EVDO_A"
	case NETWORK_TYPE_1xRTT:
		return "1xRTT"
	case NETWORK_TYPE_HSDPA:
		return "HSDPA"
	case NETWORK_TYPE_HSUPA:
		return "HSUPA"
	case NETWORK_TYPE_HSPA:
		return "HSPA"
	case NETWORK_TYPE_IDEN:
		return "IDEN"
	case NETWORK_TYPE_EVDO_B:
		return "EVDO_B"
	case NETWORK_TYPE_LTE:
		return "LTE"
	case NETWORK_TYPE_EHRPD:
		return "EHRPD"
	case NETWORK_TYPE_HSPAP:
		return "HSPAP"
	case NETWORK_TYPE_GSM:
		return "GSM"
	case NETWORK_TYPE_TD_SCDMA:
		return "TD_SCDMA"
	case NETWORK_TYPE_IWLAN:
		return "IWLAN"
	case NETWORK_TYPE_LTE_CA:
		return "LTE_CA"
	case NETWORK_TYPE_NR:
		return "NR"
	default:
		return ""
	}
}

func MakeNetworkTypeBitmask(t ...NetworkType) NetworkTypeBitmask {
	return MakeNetworkTypeBitmaskSeq(slices.Values(t))
}

func MakeNetworkTypeBitmaskSeq(seq iter.Seq[NetworkType]) NetworkTypeBitmask {
	var b NetworkTypeBitmask
	for x := range seq {
		b |= (1 << (x - 1))
	}
	return b
}

func (x NetworkTypeBitmask) Valid() bool {
	return x.String() != ""
}

func (x NetworkTypeBitmask) String() string {
	b, _ := x.MarshalText()
	return string(b)
}

func (x NetworkTypeBitmask) Seq() iter.Seq[NetworkType] {
	return func(yield func(NetworkType) bool) {
		for i := 0; i < bits.UintSize; i++ {
			if x&1 != 0 {
				if !yield(NetworkType(i + 1)) {
					return
				}
			}
			if x >>= 1; x == 0 {
				return
			}
		}
	}
}

func (x *NetworkTypeBitmask) UnmarshalText(b []byte) error {
	var err error
	*x = MakeNetworkTypeBitmaskSeq(func(yield func(NetworkType) bool) {
		if len(b) != 0 {
			for _, t := range bytes.Split(b, []byte{'|'}) {
				v, err1 := strconv.ParseInt(string(t), 10, 0)
				if err1 != nil || NetworkType(v).String() == "" {
					err = fmt.Errorf("invalid network type %q", string(t))
					return
				}
				if !yield(NetworkType(v)) {
					return
				}
			}
		}
	})
	return err
}

func (x NetworkTypeBitmask) MarshalText() ([]byte, error) {
	var b []byte
	for t := range x.Seq() {
		if t.String() == "" {
			return nil, fmt.Errorf("invalid network type %q", t.String())
		}
		if len(b) != 0 {
			b = append(b, '|')
		}
		b = strconv.AppendInt(b, int64(t), 10)
	}
	return b, nil
}

func (x Type) Valid() bool {
	return x&^(type_limit-1) == 0
}

func (x Type) String() string {
	b, _ := x.MarshalText()
	return string(b)
}

func (x Type) Seq() iter.Seq[Type] {
	return func(yield func(Type) bool) {
		for _, t := range []Type{
			TYPE_DEFAULT,
			TYPE_MMS,
			TYPE_SUPL,
			TYPE_DUN,
			TYPE_HIPRI,
			TYPE_FOTA,
			TYPE_IMS,
			TYPE_CBS,
			TYPE_IA,
			TYPE_EMERGENCY,
			TYPE_MCX,
			TYPE_XCAP,
			TYPE_VSIM,
			TYPE_BIP,
			TYPE_ENTERPRISE,
			TYPE_RCS,
		} {
			if x&t != 0 {
				if !yield(t) {
					return
				}
			}
		}
	}
}

func (x *Type) UnmarshalText(b []byte) error {
	*x = TYPE_NONE
	if len(b) != 0 {
		if string(b) == TYPE_ALL_STRING {
			*x = TYPE_ALL
			return nil
		}
		for _, t := range bytes.Split(b, []byte{'|'}) {
			switch t := string(t); t {
			case TYPE_DEFAULT_STRING:
				*x |= TYPE_DEFAULT
			case TYPE_MMS_STRING:
				*x |= TYPE_MMS
			case TYPE_SUPL_STRING:
				*x |= TYPE_SUPL
			case TYPE_DUN_STRING:
				*x |= TYPE_DUN
			case TYPE_HIPRI_STRING:
				*x |= TYPE_HIPRI
			case TYPE_FOTA_STRING:
				*x |= TYPE_FOTA
			case TYPE_IMS_STRING:
				*x |= TYPE_IMS
			case TYPE_CBS_STRING:
				*x |= TYPE_CBS
			case TYPE_IA_STRING:
				*x |= TYPE_IA
			case TYPE_EMERGENCY_STRING:
				*x |= TYPE_EMERGENCY
			case TYPE_MCX_STRING:
				*x |= TYPE_MCX
			case TYPE_XCAP_STRING:
				*x |= TYPE_XCAP
			case TYPE_VSIM_STRING:
				*x |= TYPE_VSIM
			case TYPE_BIP_STRING:
				*x |= TYPE_BIP
			case TYPE_ENTERPRISE_STRING:
				*x |= TYPE_ENTERPRISE
			case TYPE_RCS_STRING:
				*x |= TYPE_RCS
			default:
				return fmt.Errorf("unknown type %q", t)
			}
		}
	}
	return nil
}

func (x Type) MarshalText() ([]byte, error) {
	if !x.Valid() {
		return nil, fmt.Errorf("invalid type bitmask %b", x)
	}
	if x == TYPE_ALL {
		return []byte(TYPE_ALL_STRING), nil
	}
	var b []byte
	for t := range x.Seq() {
		if len(b) != 0 {
			b = append(b, ',')
		}
		switch t {
		case TYPE_DEFAULT:
			b = append(b, TYPE_DEFAULT_STRING...)
		case TYPE_MMS:
			b = append(b, TYPE_MMS_STRING...)
		case TYPE_SUPL:
			b = append(b, TYPE_SUPL_STRING...)
		case TYPE_DUN:
			b = append(b, TYPE_DUN_STRING...)
		case TYPE_HIPRI:
			b = append(b, TYPE_HIPRI_STRING...)
		case TYPE_FOTA:
			b = append(b, TYPE_FOTA_STRING...)
		case TYPE_IMS:
			b = append(b, TYPE_IMS_STRING...)
		case TYPE_CBS:
			b = append(b, TYPE_CBS_STRING...)
		case TYPE_IA:
			b = append(b, TYPE_IA_STRING...)
		case TYPE_EMERGENCY:
			b = append(b, TYPE_EMERGENCY_STRING...)
		case TYPE_MCX:
			b = append(b, TYPE_MCX_STRING...)
		case TYPE_XCAP:
			b = append(b, TYPE_XCAP_STRING...)
		case TYPE_VSIM:
			b = append(b, TYPE_VSIM_STRING...)
		case TYPE_BIP:
			b = append(b, TYPE_BIP_STRING...)
		case TYPE_ENTERPRISE:
			b = append(b, TYPE_ENTERPRISE_STRING...)
		case TYPE_RCS:
			b = append(b, TYPE_RCS_STRING...)
		default:
			panic("wtf")
		}
	}
	return b, nil
}

func (x Protocol) Valid() bool {
	return x.String() != ""
}

func (x Protocol) String() string {
	b, _ := x.MarshalText()
	return string(b)
}

func (x Protocol) MarshalText() ([]byte, error) {
	switch x {
	case PROTOCOL_IP:
		return []byte("IP"), nil
	case PROTOCOL_IPV6:
		return []byte("IPV6"), nil
	case PROTOCOL_IPV4V6:
		return []byte("IPV4V6"), nil
	case PROTOCOL_PPP:
		return []byte("PPP"), nil
	case PROTOCOL_NON_IP:
		return []byte("NON-IP"), nil
	case PROTOCOL_UNSTRUCTURED:
		return []byte("UNSTRUCTURED"), nil
	default:
		return nil, fmt.Errorf("unknown protocol %#v", x)
	}
}

func (p *Protocol) UnmarshalText(t []byte) (Protocol, error) {
	switch t := string(t); t {
	case "IP":
		return PROTOCOL_IP, nil
	case "IPV6":
		return PROTOCOL_IPV6, nil
	case "IPV4V6":
		return PROTOCOL_IPV4V6, nil
	case "PPP":
		return PROTOCOL_PPP, nil
	case "NON-IP":
		return PROTOCOL_NON_IP, nil
	case "UNSTRUCTURED":
		return PROTOCOL_UNSTRUCTURED, nil
	default:
		return PROTOCOL_UNKNOWN, fmt.Errorf("unknown protocol %q", t)
	}
}

func (x MVNOType) Valid() bool {
	return x.String() != ""
}

func (x MVNOType) String() string {
	b, _ := x.MarshalText()
	return string(b)
}

func (x MVNOType) MarshalText() ([]byte, error) {
	switch x {
	case MVNO_TYPE_SPN:
		return []byte("spn"), nil
	case MVNO_TYPE_IMSI:
		return []byte("imsi"), nil
	case MVNO_TYPE_GID:
		return []byte("gid"), nil
	case MVNO_TYPE_ICCID:
		return []byte("iccid"), nil
	default:
		return nil, fmt.Errorf("unknown mvno type %#v", x)
	}
}

func (x *MVNOType) UnmarshalText(t []byte) (MVNOType, error) {
	switch t := string(t); t {
	case "spn":
		return MVNO_TYPE_SPN, nil
	case "imsi":
		return MVNO_TYPE_IMSI, nil
	case "gid":
		return MVNO_TYPE_GID, nil
	case "iccid":
		return MVNO_TYPE_ICCID, nil
	default:
		return MVNO_TYPE_UNKNOWN, fmt.Errorf("unknown mvno type %q", t)
	}
}

func (x Infrastructure) Valid() bool {
	return x.String() != ""
}

func (x Infrastructure) String() string {
	b, _ := x.MarshalText()
	return string(b)
}

func (x Infrastructure) MarshalText() ([]byte, error) {
	switch x {
	case INFRASTRUCTURE_CELLULAR:
		return []byte("cellular"), nil
	case INFRASTRUCTURE_SATELLITE:
		return []byte("satellite"), nil
	case INFRASTRUCTURE_CELLULAR | INFRASTRUCTURE_SATELLITE:
		return []byte("cellular|satellite"), nil
	default:
		return nil, fmt.Errorf("unknown infrastructure type %#v", x)
	}
}

func (p *Infrastructure) UnmarshalText(t []byte) (Infrastructure, error) {
	switch t := string(t); t {
	case "cellular":
		return INFRASTRUCTURE_CELLULAR, nil
	case "satellite":
		return INFRASTRUCTURE_SATELLITE, nil
	case "cellular|satellite", "satellite|cellular":
		return INFRASTRUCTURE_CELLULAR | INFRASTRUCTURE_SATELLITE, nil
	default:
		return INFRASTRUCTURE_CELLULAR | INFRASTRUCTURE_SATELLITE, fmt.Errorf("unknown infrastructure type %q", t)
	}
}

type Setting struct {
	EntryName                   string
	APNName                     string
	ProxyAddress                string
	ProxyPort                   int
	MMSC                        string
	MMSProxyAddress             string
	MMSProxyPort                int
	User                        string
	Password                    string
	AuthType                    AuthType
	APNTypeBitmask              Type
	OperatorNumeric             string
	Protocol                    Protocol
	RoamingProtocol             Protocol
	MTUv4                       int
	MTUv6                       int
	CarrierEnabled              bool // apn is enabled
	ProfileID                   int
	NetworkTypeBitmask          NetworkTypeBitmask
	LingeringNetworkTypeBitmask NetworkTypeBitmask
	Persistent                  bool // modem cognitive
	MaxConns                    int
	WaitTime                    int
	MaxConnsTime                int
	MVNOType                    MVNOType
	MVNOMatchData               string
	APNSetID                    int
	CarrierID                   int
	Skip464XLAT                 Skip464XLAT
	AlwaysOn                    bool
	InfrastructureBitmask       Infrastructure
	ESIMBootstrapProvisioning   bool
	// skipped ID, ProfileID, PermanentFailed (those are runtime fields)
}

// Empty returns a new Setting with the default values.
func Empty() Setting {
	return Setting{
		EntryName:                   "",
		APNName:                     "",
		ProxyAddress:                "",
		ProxyPort:                   -1,
		MMSC:                        "",
		MMSProxyAddress:             "",
		MMSProxyPort:                -1,
		User:                        "",
		Password:                    "",
		AuthType:                    AUTH_TYPE_UNKNOWN,
		APNTypeBitmask:              0,
		OperatorNumeric:             "",
		Protocol:                    PROTOCOL_UNKNOWN,
		RoamingProtocol:             PROTOCOL_UNKNOWN,
		MTUv4:                       0,
		MTUv6:                       0,
		CarrierEnabled:              false,
		ProfileID:                   0,
		NetworkTypeBitmask:          0,
		LingeringNetworkTypeBitmask: 0,
		Persistent:                  false,
		MaxConns:                    0,
		WaitTime:                    0,
		MaxConnsTime:                0,
		MVNOType:                    MVNO_TYPE_UNKNOWN,
		MVNOMatchData:               "",
		APNSetID:                    0, // NO_SET_SET
		CarrierID:                   0,
		Skip464XLAT:                 SKIP_464XLAT_DEFAULT,
		AlwaysOn:                    false,
		InfrastructureBitmask:       INFRASTRUCTURE_CELLULAR | INFRASTRUCTURE_SATELLITE,
		ESIMBootstrapProvisioning:   false,
	}
}

// Check checks certain APN fields (same as ApnSetting.build).
func (s Setting) Check() error {
	if s.EntryName == "" {
		return fmt.Errorf("entry name is required")
	}
	if s.APNName == "" {
		return fmt.Errorf("apn name is required")
	}
	if !s.APNTypeBitmask.Valid() {
		return fmt.Errorf("invalid apn type bitmask")
	}
	if s.APNTypeBitmask&TYPE_MMS != 0 {
		if u, err := url.Parse(s.MMSProxyAddress); err == nil && u.Scheme != "" {
			return fmt.Errorf("mms proxy should be a hostname, not a url")
		}
	}
	return nil
}
