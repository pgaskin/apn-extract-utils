package apn

import (
	"bytes"
	"fmt"
	"iter"
	"math/bits"
	"slices"
	"strconv"
)

// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/telephony/java/android/telephony/ServiceState.java;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae

type RILRadioTechnology int

const (
	RIL_RADIO_TECHNOLOGY_UNKNOWN RILRadioTechnology = iota
	RIL_RADIO_TECHNOLOGY_GPRS
	RIL_RADIO_TECHNOLOGY_EDGE
	RIL_RADIO_TECHNOLOGY_UMTS
	RIL_RADIO_TECHNOLOGY_IS95A
	RIL_RADIO_TECHNOLOGY_IS95B
	RIL_RADIO_TECHNOLOGY_1xRTT
	RIL_RADIO_TECHNOLOGY_EVDO_0
	RIL_RADIO_TECHNOLOGY_EVDO_A
	RIL_RADIO_TECHNOLOGY_HSDPA
	RIL_RADIO_TECHNOLOGY_HSUPA
	RIL_RADIO_TECHNOLOGY_HSPA
	RIL_RADIO_TECHNOLOGY_EVDO_B
	RIL_RADIO_TECHNOLOGY_EHRPD
	RIL_RADIO_TECHNOLOGY_LTE
	RIL_RADIO_TECHNOLOGY_HSPAP
	RIL_RADIO_TECHNOLOGY_GSM
	RIL_RADIO_TECHNOLOGY_TD_SCDMA
	RIL_RADIO_TECHNOLOGY_IWLAN
	RIL_RADIO_TECHNOLOGY_LTE_CA
	RIL_RADIO_TECHNOLOGY_NR
	NEXT_RIL_RADIO_TECHNOLOGY
)

type BearerBitmask int

const (
	BEARER_BITMASK_GPRS     BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_GPRS - 1))
	BEARER_BITMASK_EDGE     BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_EDGE - 1))
	BEARER_BITMASK_UMTS     BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_UMTS - 1))
	BEARER_BITMASK_IS95A    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_IS95A - 1))
	BEARER_BITMASK_IS95B    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_IS95B - 1))
	BEARER_BITMASK_1xRTT    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_1xRTT - 1))
	BEARER_BITMASK_EVDO_0   BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_EVDO_0 - 1))
	BEARER_BITMASK_EVDO_A   BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_EVDO_A - 1))
	BEARER_BITMASK_HSDPA    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_HSDPA - 1))
	BEARER_BITMASK_HSUPA    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_HSUPA - 1))
	BEARER_BITMASK_HSPA     BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_HSPA - 1))
	BEARER_BITMASK_EVDO_B   BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_EVDO_B - 1))
	BEARER_BITMASK_EHRPD    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_EHRPD - 1))
	BEARER_BITMASK_LTE      BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_LTE - 1))
	BEARER_BITMASK_HSPAP    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_HSPAP - 1))
	BEARER_BITMASK_GSM      BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_GSM - 1))
	BEARER_BITMASK_TD_SCDMA BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_TD_SCDMA - 1))
	BEARER_BITMASK_IWLAN    BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_IWLAN - 1))
	BEARER_BITMASK_LTE_CA   BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_LTE_CA - 1))
	BEARER_BITMASK_NR       BearerBitmask = (1 << (RIL_RADIO_TECHNOLOGY_NR - 1))
)

func (x RILRadioTechnology) Valid() bool {
	return x > 0 && x < NEXT_RIL_RADIO_TECHNOLOGY
}

func (x RILRadioTechnology) ToNetworkType() NetworkType {
	switch x {
	case RIL_RADIO_TECHNOLOGY_GPRS:
		return NETWORK_TYPE_GPRS
	case RIL_RADIO_TECHNOLOGY_EDGE:
		return NETWORK_TYPE_EDGE
	case RIL_RADIO_TECHNOLOGY_UMTS:
		return NETWORK_TYPE_UMTS
	case RIL_RADIO_TECHNOLOGY_HSDPA:
		return NETWORK_TYPE_HSDPA
	case RIL_RADIO_TECHNOLOGY_HSUPA:
		return NETWORK_TYPE_HSUPA
	case RIL_RADIO_TECHNOLOGY_HSPA:
		return NETWORK_TYPE_HSPA
	case RIL_RADIO_TECHNOLOGY_IS95A:
		fallthrough
	case RIL_RADIO_TECHNOLOGY_IS95B:
		return NETWORK_TYPE_CDMA
	case RIL_RADIO_TECHNOLOGY_1xRTT:
		return NETWORK_TYPE_1xRTT
	case RIL_RADIO_TECHNOLOGY_EVDO_0:
		return NETWORK_TYPE_EVDO_0
	case RIL_RADIO_TECHNOLOGY_EVDO_A:
		return NETWORK_TYPE_EVDO_A
	case RIL_RADIO_TECHNOLOGY_EVDO_B:
		return NETWORK_TYPE_EVDO_B
	case RIL_RADIO_TECHNOLOGY_EHRPD:
		return NETWORK_TYPE_EHRPD
	case RIL_RADIO_TECHNOLOGY_LTE:
		return NETWORK_TYPE_LTE
	case RIL_RADIO_TECHNOLOGY_HSPAP:
		return NETWORK_TYPE_HSPAP
	case RIL_RADIO_TECHNOLOGY_GSM:
		return NETWORK_TYPE_GSM
	case RIL_RADIO_TECHNOLOGY_TD_SCDMA:
		return NETWORK_TYPE_TD_SCDMA
	case RIL_RADIO_TECHNOLOGY_IWLAN:
		return NETWORK_TYPE_IWLAN
	case RIL_RADIO_TECHNOLOGY_LTE_CA:
		return NETWORK_TYPE_LTE_CA
	case RIL_RADIO_TECHNOLOGY_NR:
		return NETWORK_TYPE_NR
	default:
		return NETWORK_TYPE_UNKNOWN
	}
}

func MakeBearerBitmask(t ...RILRadioTechnology) BearerBitmask {
	return MakeBearerBitmaskSeq(slices.Values(t))
}

func MakeBearerBitmaskSeq(seq iter.Seq[RILRadioTechnology]) BearerBitmask {
	var b BearerBitmask
	for x := range seq {
		b |= (1 << (x - 1))
	}
	return b
}

func (x BearerBitmask) Valid() bool {
	return x.String() != ""
}

func (x BearerBitmask) String() string {
	b, _ := x.MarshalText()
	return string(b)
}

func (x BearerBitmask) Seq() iter.Seq[RILRadioTechnology] {
	return func(yield func(RILRadioTechnology) bool) {
		for i := 0; i < bits.UintSize; i++ {
			if x&1 != 0 {
				if !yield(RILRadioTechnology(i + 1)) {
					return
				}
			}
			if x >>= 1; x == 0 {
				return
			}
		}
	}
}

func (x *BearerBitmask) UnmarshalText(b []byte) error {
	var err error
	*x = MakeBearerBitmaskSeq(func(yield func(RILRadioTechnology) bool) {
		if len(b) != 0 {
			for _, t := range bytes.Split(b, []byte{'|'}) {
				v, err1 := strconv.ParseInt(string(t), 10, 0)
				if err1 != nil || !RILRadioTechnology(v).Valid() {
					err = fmt.Errorf("invalid bearer %q", string(t))
					return
				}
				if !yield(RILRadioTechnology(v)) {
					return
				}
			}
		}
	})
	return err
}

func (x BearerBitmask) MarshalText() ([]byte, error) {
	var b []byte
	for t := range x.Seq() {
		if !t.Valid() {
			return nil, fmt.Errorf("invalid bearer %d", t)
		}
		if len(b) != 0 {
			b = append(b, '|')
		}
		b = strconv.AppendInt(b, int64(t), 10)
	}
	return b, nil
}

func bitmaskHasTech[T, U ~int](x T, t U) bool {
	if x == 0 {
		return true
	}
	if t >= 1 {
		return ((x & (1 << (t - 1))) != 0)
	}
	return false
}

func ConvertNetworkTypeBitmaskToBearerBitmask(ntb NetworkTypeBitmask) BearerBitmask {
	var bb BearerBitmask
	if ntb != 0 {
		for bearer := range NEXT_RIL_RADIO_TECHNOLOGY {
			if bitmaskHasTech(ntb, bearer.ToNetworkType()) {
				bb |= (1 << (bearer - 1))
			}
		}
	}
	return bb
}
func ConvertBearerBitmaskToNetworkTypeBitmask(bb BearerBitmask) NetworkTypeBitmask {
	var ntb NetworkTypeBitmask
	if bb != 0 {
		for bearer := range NEXT_RIL_RADIO_TECHNOLOGY {
			if bitmaskHasTech(bb, bearer) {
				if t := bearer.ToNetworkType(); t >= 1 {
					ntb |= (1 << (t - 1))
				}
			}
		}
	}
	return ntb
}
