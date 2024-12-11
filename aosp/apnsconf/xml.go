package apnsconf

import (
	"fmt"
	"iter"
	"strconv"

	"github.com/pgaskin/apn-extract-utils/aosp/apn"
)

const Version = 8

// https://cs.android.com/android/platform/superproject/main/+/main:packages/providers/TelephonyProvider/src/com/android/providers/telephony/TelephonyProvider.java;l=2716;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae (getRow)
// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/telephony/java/android/telephony/data/ApnSetting.java;l=1466;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae (toContentValues, makeApnSetting)
// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/telephony/java/android/telephony/ServiceState.java;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae (conversion functions)
// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/base/core/java/android/provider/Telephony.java;l=3108;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae
// https://github.com/LineageOS/android_vendor_lineage/blob/56ec683ee675eefa2fb618c06e8e29d47f2fffdb/tools/apns-conf.xsd (for confirmation)

func XMLAttrSeq(s apn.Setting, err *error) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		*err = func() error {
			// mcc/mnc/mvno_type/mvno_match_data will be replaced entirely with carrier_id matching in the future
			if s.OperatorNumeric != "" {
				if n := len(s.OperatorNumeric); n != 5 && n != 6 {
					return fmt.Errorf("invalid operator mccmnc length %d", n)
				}
				for _, c := range s.OperatorNumeric {
					if c < '0' && c > '9' {
						return fmt.Errorf("invalid operator mccmnc %q", s.OperatorNumeric)
					}
				}
				if !yield("mcc", s.OperatorNumeric[:3]) {
					return nil
				}
				if !yield("mnc", s.OperatorNumeric[3:]) {
					return nil
				}
			}
			if s.EntryName == "" {
				return fmt.Errorf("entry name is required")
			}
			if !yield("carrier", s.EntryName) {
				return nil
			}
			if v := s.APNName; v != "" {
				if !yield("apn", v) {
					return nil
				}
			}
			if v := s.User; v != "" {
				if !yield("user", v) {
					return nil
				}
			}
			// N/A: string: server
			if v := s.Password; v != "" {
				if !yield("password", v) {
					return nil
				}
			}
			if v := s.ProxyAddress; v != "" {
				if !yield("proxy", v) {
					return nil
				}
			}
			if v := s.ProxyPort; v > 0 {
				if !yield("port", strconv.Itoa(int(s.ProxyPort))) {
					return nil
				}
			}
			if v := s.MMSProxyAddress; v != "" {
				if !yield("mmsproxy", v) {
					return nil
				}
			}
			if v := s.MMSProxyPort; v > 0 {
				if !yield("mmsport", strconv.Itoa(int(s.ProxyPort))) {
					return nil
				}
			}
			if v := s.MMSC; v != "" {
				if !yield("mmsc", v) {
					return nil
				}
			}
			if v := s.APNTypeBitmask; v != 0 {
				if b, err := v.MarshalText(); err != nil {
					return fmt.Errorf("invalid apn type bitmask: %w", err)
				} else if !yield("type", string(b)) {
					return nil
				}
			}
			if v := s.Protocol; v != apn.PROTOCOL_UNKNOWN {
				if b, err := v.MarshalText(); err != nil {
					return fmt.Errorf("invalid protocol: %w", err)
				} else if !yield("protocol", string(b)) {
					return nil
				}
			}
			if v := s.RoamingProtocol; v != apn.PROTOCOL_UNKNOWN {
				if b, err := v.MarshalText(); err != nil {
					return fmt.Errorf("invalid roaming protocol: %w", err)
				} else if !yield("roaming_protocol", string(b)) {
					return nil
				}
			}
			if v := s.AuthType; v != apn.AUTH_TYPE_UNKNOWN {
				if !yield("authtype", strconv.Itoa(int(v))) {
					return nil
				}
			}
			// N/A: int: bearer (no longer supported, replaced with network_type_bitmask)
			if v := s.ProfileID; v != 0 {
				if !yield("profile_id", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.MaxConns; v != 0 {
				if !yield("max_conns", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.WaitTime; v != 0 {
				if !yield("wait_time", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.MaxConnsTime; v != 0 {
				if !yield("max_conns_time", strconv.Itoa(v)) {
					return nil
				}
			}
			// int: mtu (deprecated, use mtu_v4 or mtu_v6 instead)
			if v := s.MTUv4; v > 0 {
				if !yield("mtu_v4", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.MTUv6; v > 0 {
				if !yield("mtu_v6", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.APNSetID; v != 0 { // NO_SET_SET
				if !yield("apn_set_id", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.CarrierID; v != 0 {
				if !yield("carrier_id", strconv.Itoa(v)) {
					return nil
				}
			}
			if v := s.Skip464XLAT; v != apn.SKIP_464XLAT_DEFAULT {
				if !yield("skip_464xlat", strconv.Itoa(int(v))) {
					return nil
				}
			}
			if v := s.CarrierEnabled; v != true {
				if !yield("carrier_enabled", "false") {
					return nil
				}
			}
			if v := s.Persistent; v != false {
				if !yield("modem_cognitive", "true") {
					return nil
				}
			}
			// N/A: bool: user_visible
			// N/A: bool: user_editable
			if v := s.AlwaysOn; v != false {
				if !yield("always_on", "true") {
					return nil
				}
			}
			if v := s.ESIMBootstrapProvisioning; v != false {
				if !yield("esim_bootstrap_provisioning", "true") {
					return nil
				}
			}
			if v := s.InfrastructureBitmask; v != 0 { // don't include if not explicitly set
				if b, err := v.MarshalText(); err != nil {
					return fmt.Errorf("invalid infrastructure bitmask: %w", err)
				} else if !yield("infrastructure_bitmask", string(b)) {
					return nil
				}
			}
			if s.NetworkTypeBitmask != 0 {
				if !s.NetworkTypeBitmask.Valid() {
					return fmt.Errorf("invalid network type bitmask")
				}
				bearerBitmask := apn.ConvertNetworkTypeBitmaskToBearerBitmask(s.NetworkTypeBitmask)
				bearerBitmaskBack := apn.ConvertBearerBitmaskToNetworkTypeBitmask(bearerBitmask)
				// if not present, ApnSetting.makeApnSetting will create the network_type_bitmask from the bearer_bitmask
				// in newer versions of android, the sample apns-conf.xml only includes network_type_bitmask, but we'll prefer using the bearer_bitmask for compatibility if it effectively equals the network bitmask
				if bearerBitmaskBack != s.NetworkTypeBitmask {
					if v := s.NetworkTypeBitmask; v != 0 {
						if b, err := v.MarshalText(); err != nil {
							return fmt.Errorf("invalid network type bitmask: %w", err)
						} else if !yield("network_type_bitmask", string(b)) {
							return nil
						}
					}
				}
				if v := bearerBitmask; v != 0 {
					if b, err := v.MarshalText(); err != nil {
						return fmt.Errorf("invalid bearer bitmask: %w", err)
					} else if !yield("bearer_bitmask", string(b)) {
						return nil
					}
				}
			}
			if v := s.LingeringNetworkTypeBitmask; v != 0 {
				if b, err := v.MarshalText(); err != nil {
					return fmt.Errorf("invalid lingering network type bitmask: %w", err)
				} else if !yield("lingering_network_type_bitmask", string(b)) {
					return nil
				}
			}
			if v := s.MVNOType; v != apn.MVNO_TYPE_UNKNOWN {
				if b, err := v.MarshalText(); err != nil {
					return fmt.Errorf("invalid mvno type: %w", err)
				} else if !yield("mvno_type", string(b)) {
					return nil
				}
				if s.MVNOMatchData == "" {
					return fmt.Errorf("no mvno match data provided even though match type was %s", v)
				}
				if !yield("mvno_match_data", string(s.MVNOMatchData)) {
					return nil
				}
			}
			return nil
		}()
	}
}
