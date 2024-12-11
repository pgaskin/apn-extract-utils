package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/pgaskin/apn-extract-utils/aosp/apn"
	"github.com/pgaskin/apn-extract-utils/aosp/apnsconf"
	"github.com/pgaskin/apn-extract-utils/aosp/carrier_list"
	"github.com/pgaskin/apn-extract-utils/aosp/carrier_settings"
	"github.com/pgaskin/apn-extract-utils/aosp/carrierid"
	"github.com/pgaskin/xmlwriter"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// TODO: make this actually a proper command
// TODO: refactor main logic into new source/carriersettings package
// TODO: compare output with lineage carriersettings-extractor

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	var (
		src                           = os.DirFS("/data/android/lineage")
		pixelCarrierSettings, _       = fs.Sub(src, "vendor/google/caiman/proprietary/product/etc/CarrierSettings")
		onlyCarrierIDMatch            = false
		expandAdditionalFromCarrierID = true
		debugDumpText                 = true
		filterNameSuffix              = "_ca"
	)

	txt := prototext.MarshalOptions{
		EmitUnknown:  true,
		Indent:       "  ",
		AllowPartial: true,
	}

	if debugDumpText {
		os.MkdirAll("dbg", 0777)
	}

	carrierId, err := openProto[*carrierid.CarrierList](src, "packages/providers/TelephonyProvider/assets/sdk34_carrier_id/carrier_list.pb")
	if err != nil {
		panic(err)
	}
	slog.Info("loaded carrier identification", "total", len(carrierId.CarrierId))

	if debugDumpText {
		buf, _ := txt.Marshal(carrierId)
		os.WriteFile(filepath.Join("dbg", "carrierId.textpb"), buf, 0666)
	}

	carrierList, err := openProto[*carrier_list.CarrierList](pixelCarrierSettings, "carrier_list.pb")
	if err != nil {
		panic(err)
	}
	slog.Info("loaded carrier list", "total", len(carrierList.Entry))

	if debugDumpText {
		buf, _ := txt.Marshal(carrierList)
		os.WriteFile(filepath.Join("dbg", "carrier_list.textpb"), buf, 0666)
	}

	genericSettings, err := openProto[*carrier_settings.MultiCarrierSettings](pixelCarrierSettings, "others.pb")
	if err != nil {
		panic(err)
	}

	if debugDumpText {
		buf, _ := txt.Marshal(genericSettings)
		os.WriteFile(filepath.Join("dbg", "others.textpb"), buf, 0666)
	}

	allSettings := map[string]*carrier_settings.CarrierSettings{} // [canonicalName]
	for _, cs := range genericSettings.Setting {
		if filterNameSuffix != "" && !strings.HasSuffix(cs.GetCanonicalName(), filterNameSuffix) {
			continue
		}
		slog.Debug("loaded generic settings", "canonical_name", cs.GetCanonicalName())
		allSettings[*cs.CanonicalName] = cs
	}
	slog.Info("loaded generic settings", "total", len(genericSettings.Setting))

	var specificSettings int
	if err := fs.WalkDir(pixelCarrierSettings, ".", func(name string, d fs.DirEntry, err error) error {
		if path.Ext(name) != ".pb" || name == "carrier_list.pb" || name == "others.pb" {
			return nil
		}
		carrierSettings, err := openProto[*carrier_settings.CarrierSettings](pixelCarrierSettings, name)
		if err != nil {
			return err
		}
		if filterNameSuffix != "" && !strings.HasSuffix(carrierSettings.GetCanonicalName(), filterNameSuffix) {
			return nil
		}
		if debugDumpText {
			buf, _ := txt.Marshal(carrierSettings)
			os.WriteFile(filepath.Join("dbg", strings.TrimSuffix(filepath.FromSlash(name), ".pb")+".textpb"), buf, 0666)
		}
		specificSettings++
		slog.Debug("loaded settings", "name", name, "canonical_name", carrierSettings.GetCanonicalName())
		if _, ok := allSettings[*carrierSettings.CanonicalName]; ok {
			slog.Warn("replacing generic settings", "name", name, "canonical_name", carrierSettings.GetCanonicalName())
		}
		allSettings[*carrierSettings.CanonicalName] = carrierSettings
		return nil
	}); err != nil {
		panic(err)
	}
	slog.Info("loaded specific settings", "total", specificSettings)

	carrierMap := map[string]*carrier_list.CarrierMap{} // [canonicalName]
	for _, canonicalName := range slices.Sorted(maps.Keys(allSettings)) {
		i := slices.IndexFunc(carrierList.Entry, func(c *carrier_list.CarrierMap) bool {
			return c.GetCanonicalName() == canonicalName
		})
		if i == -1 {
			slog.Error("failed to find carrier_list entry for carrier, dropping", "canonical_name", canonicalName)
			delete(allSettings, canonicalName)
			continue
		}
		carrierMap[canonicalName] = carrierList.Entry[i]
	}
	slog.Info("mapped carrier_settings to carrier_list entries")

	carrierMapID := map[string][]*carrierid.CarrierId{} // [canonicalName]
	carrierIdMatchedExact := map[int]string{}
	for _, canonicalName := range slices.Sorted(maps.Keys(allSettings)) {
		carrier := carrierMap[canonicalName]
		for _, wantMatch := range carrier.CarrierId {
			i := slices.IndexFunc(carrierId.CarrierId, func(c *carrierid.CarrierId) bool {
				return slices.ContainsFunc(c.CarrierAttribute, func(a *carrierid.CarrierAttribute) bool {
					if !slices.Contains(a.MccmncTuple, *wantMatch.MccMnc) {
						return false
					}
					if wantMatch.MvnoData == nil {
						return len(a.ImsiPrefixXpattern) == 0 &&
							len(a.Spn) == 0 &&
							len(a.Plmn) == 0 &&
							len(a.Gid1) == 0 &&
							len(a.Gid2) == 0 &&
							len(a.PreferredApn) == 0 &&
							len(a.IccidPrefix) == 0 &&
							len(a.PrivilegeAccessRule) == 0
					}
					switch wantData := wantMatch.MvnoData.(type) {
					case *carrier_list.CarrierId_Spn:
						return len(a.ImsiPrefixXpattern) == 0 &&
							slices.ContainsFunc(a.Spn, func(e string) bool {
								return strings.EqualFold(e, wantData.Spn)
							}) &&
							len(a.Plmn) == 0 &&
							len(a.Gid1) == 0 &&
							len(a.Gid2) == 0 &&
							len(a.PreferredApn) == 0 &&
							len(a.IccidPrefix) == 0 &&
							len(a.PrivilegeAccessRule) == 0
					case *carrier_list.CarrierId_Imsi:
						return slices.ContainsFunc(a.ImsiPrefixXpattern, func(matchPattern string) bool {
							if len(matchPattern) < len(wantData.Imsi) {
								return false
							}
							for i, wantDigit := range []byte(wantData.Imsi) {
								matchDigit := matchPattern[i]
								switch {
								case wantDigit == matchDigit:
								case (wantDigit == 'x' || wantDigit == 'X') && (matchDigit == 'x' || matchDigit == 'X'):
								case matchDigit == 'x' || matchDigit == 'X':
								default:
									return false
								}
							}
							return true
						}) &&
							len(a.Spn) == 0 &&
							len(a.Plmn) == 0 &&
							len(a.Gid1) == 0 &&
							len(a.Gid2) == 0 &&
							len(a.PreferredApn) == 0 &&
							len(a.IccidPrefix) == 0 &&
							len(a.PrivilegeAccessRule) == 0
					case *carrier_list.CarrierId_Gid1:
						return len(a.ImsiPrefixXpattern) == 0 &&
							len(a.Spn) == 0 &&
							len(a.Plmn) == 0 &&
							slices.ContainsFunc(a.Gid1, func(e string) bool {
								return strings.EqualFold(e, wantData.Gid1)
							}) &&
							len(a.Gid2) == 0 &&
							len(a.PreferredApn) == 0 &&
							len(a.IccidPrefix) == 0 &&
							len(a.PrivilegeAccessRule) == 0
					default:
						panic("unhandled mnvo data type")
					}
				})
			})
			if i != -1 {
				carrierMapID[canonicalName] = append(carrierMapID[canonicalName], carrierId.CarrierId[i])
				if other, ok := carrierIdMatchedExact[i]; ok {
					slog.Warn("multiple carriersettings carriers matched a single carrierId exactly", "canonical_name", canonicalName, "other_canonical_name", other, "carrier_id", *carrierId.CarrierId[i].CanonicalId)
				}
				carrierIdMatchedExact[i] = canonicalName
			}
		}
		if n := len(carrierMapID[canonicalName]); n == 0 {
			slog.Warn("no carrierId match for carriersettings carrier", "canonical_name", canonicalName)
			continue
		} else if n > 1 {
			slog.Warn("more than one carrierId match for carriersettings carrier", "canonical_name", canonicalName)
			continue
		}
	}
	slog.Info("mapped carrier_list entries to carrierId (exact matches)", "have", len(carrierMapID), "missing", len(carrierMap)-len(carrierMapID))
	carrierIdMatchedPLMNOnly := map[int]string{}
	for _, canonicalName := range slices.Sorted(maps.Keys(allSettings)) {
		carrier := carrierMap[canonicalName]
		if _, ok := carrierMapID[canonicalName]; ok {
			continue
		}
		for _, wantMatch := range carrier.CarrierId {
			var possibleMatches []int
			for i, c := range carrierId.CarrierId {
				if _, ok := carrierIdMatchedExact[i]; ok {
					continue
				}
				if slices.ContainsFunc(c.CarrierAttribute, func(a *carrierid.CarrierAttribute) bool {
					return slices.Contains(a.MccmncTuple, *wantMatch.MccMnc)
				}) {
					possibleMatches = append(possibleMatches, i)
				}
			}
			if len(possibleMatches) == 1 {
				i := possibleMatches[0]
				carrierMapID[canonicalName] = append(carrierMapID[canonicalName], carrierId.CarrierId[i])
				if other, ok := carrierIdMatchedPLMNOnly[i]; ok {
					slog.Warn("multiple carriersettings carriers which idn't match a single carrierId exactly matched a non-exactly-matched carrierId by the just the PLMN", "canonical_name", canonicalName, "other_canonical_name", other, "carrier_id", *carrierId.CarrierId[i].CanonicalId)
				}
				carrierIdMatchedPLMNOnly[i] = canonicalName
				slog.Warn("added a carrierId match for a carriersettings carrier by the plmn only", "canonical_name", canonicalName, "carrier_id", *carrierId.CarrierId[i].CanonicalId, "name", carrierId.CarrierId[i].GetCarrierName)
			}
		}
	}
	slog.Info("mapped carrier_list entries to carrierId (plus plmn-only matches of remaining carrierId entries)", "have", len(carrierMapID), "missing", len(carrierMap)-len(carrierMapID))

	type ConvertedAPN struct {
		Comment       string
		CanonicalName string
		Setting       apn.Setting
		UserVisible   bool
		UserEditable  bool
	}
	var apns []ConvertedAPN
	for _, canonicalName := range slices.Sorted(maps.Keys(allSettings)) {
		carrier := carrierMap[canonicalName]
		carrierIDs := carrierMapID[canonicalName]
		carrierSettings := allSettings[canonicalName]

		slog := slog.With("canonical_name", canonicalName)

		if carrierSettings.Apns == nil {
			slog.Debug("no apns for carrier")
			continue
		}

		for _, src := range carrierSettings.Apns.Apn {
			s := apn.Empty()
			s.CarrierEnabled = true
			s.InfrastructureBitmask = 0

			if err := func() error {
				s.EntryName = src.GetName()
				s.APNName = src.GetValue()

				if s.EntryName == "" {
					return fmt.Errorf("empty apn name")
				}

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
						s.APNTypeBitmask |= apn.TYPE_XCAP // TODO: is this correct?
					case carrier_settings.ApnItem_RCS:
						s.APNTypeBitmask |= apn.TYPE_RCS
					default:
						panic("unhandled apn type")
					}
				}

				// yes, it's named bearerbitmask, but it's actually the network type bitmask
				var bb apn.BearerBitmask
				if v := src.GetBearerBitmask(); v != "0" {
					if err := bb.UnmarshalText([]byte(v)); err != nil {
						return err
					}
					s.NetworkTypeBitmask = apn.ConvertBearerBitmaskToNetworkTypeBitmask(bb)
					if bb1 := apn.ConvertNetworkTypeBitmaskToBearerBitmask(s.NetworkTypeBitmask); bb1 != bb {
						slog.Warn("lossy bearer bitmask conversion", "from", bb, "to", s.NetworkTypeBitmask, "back", bb1)
					}
				}

				if v := src.GetServer(); v != "" {
					slog.Warn("no mapping for server param", "value", v)
				}
				s.ProxyAddress = src.GetProxy()
				if v := src.GetPort(); v != "" {
					v, err := strconv.ParseInt(v, 10, 0)
					if err != nil {
						return fmt.Errorf("port: %w", err)
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
						return fmt.Errorf("port: %w", err)
					}
					s.MMSProxyPort = int(v)
				}

				switch src.GetProtocol() {
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
					panic("unhandled protocol")
				}

				switch src.GetRoamingProtocol() {
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
					panic("unhandled protocol")
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

				switch src.GetSkip_464Xlat() {
				case carrier_settings.ApnItem_SKIP_464XLAT_DEFAULT:
					s.Skip464XLAT = apn.SKIP_464XLAT_DEFAULT
				case carrier_settings.ApnItem_SKIP_464XLAT_DISABLE:
					s.Skip464XLAT = apn.SKIP_464XLAT_DISABLE
				case carrier_settings.ApnItem_SKIP_464XLAT_ENABLE:
					s.Skip464XLAT = apn.SKIP_464XLAT_ENABLE
				default:
					panic("unhandled skip 464xlat value")
				}

				return nil
			}(); err != nil {
				slog.Error("failed to convert apn, skipping", "error", err)
				continue
			}

			if err := s.Check(); err != nil {
				slog.Warn("check failed for apn", "error", err)
			}

			if onlyCarrierIDMatch {
				for _, c := range carrierIDs {
					tmp := s
					tmp.CarrierID = int(*c.CanonicalId)

					apns = append(apns, ConvertedAPN{
						Comment:       canonicalName,
						CanonicalName: canonicalName,
						Setting:       tmp,
						UserVisible:   src.GetUserVisible(),
						UserEditable:  src.GetUserEditable(),
					})
				}
			} else {
				var canonicalIDs []int
				for _, cm := range carrierIDs {
					canonicalIDs = append(canonicalIDs, int(*cm.CanonicalId))
				}
				if len(canonicalIDs) == 0 {
					canonicalIDs = append(canonicalIDs, -1)
				}
				for _, canonicalID := range canonicalIDs {
					for _, c := range carrier.CarrierId {
						if c.MccMnc == nil {
							slog.Warn("skipping carrier match without mccmnc")
							continue
						}
						tmp := s
						if canonicalID != -1 {
							tmp.CarrierID = canonicalID
						}
						tmp.OperatorNumeric = *c.MccMnc

						if c.MvnoData != nil {
							switch d := c.GetMvnoData().(type) {
							case *carrier_list.CarrierId_Spn:
								tmp.MVNOType = apn.MVNO_TYPE_SPN
								tmp.MVNOMatchData = d.Spn
							case *carrier_list.CarrierId_Imsi:
								tmp.MVNOType = apn.MVNO_TYPE_IMSI
								tmp.MVNOMatchData = d.Imsi
							case *carrier_list.CarrierId_Gid1:
								tmp.MVNOType = apn.MVNO_TYPE_GID
								tmp.MVNOMatchData = d.Gid1
							default:
								panic("unhandled mnvo data type")
							}
						}

						apns = append(apns, ConvertedAPN{
							Comment:       canonicalName,
							CanonicalName: canonicalName,
							Setting:       tmp,
							UserVisible:   src.GetUserVisible(),
							UserEditable:  src.GetUserEditable(),
						})
					}
				}
				if expandAdditionalFromCarrierID {
					for _, cm := range carrierIDs {
						_ = cm
						// TODO: expand from cm.CarrierAttribute? will need to warn if there is more than one match condition, since that can't be expressed
					}
				}
			}
		}
	}
	slog.Info("converted apns", "total", len(apns))

	w := xmlwriter.New(os.Stdout)
	w.Indent("  ")
	w.Start(nil, "apns", xmlwriter.NS("").Bind(""))
	w.Attr(nil, "version", strconv.Itoa(apnsconf.Version))
	var last string
	for _, c := range apns {
		if last == "" || last != c.Comment {
			last = c.Comment
			w.BlankLine()
			w.Comment(true, " "+c.Comment+" ")
		}
		w.Start(nil, "apn")
		var err error
		for k, v := range apnsconf.XMLAttrSeq(c.Setting, &err) {
			w.Attr(nil, k, v)
		}
		if !c.UserVisible {
			w.Attr(nil, "user_visible", "false")
		}
		if !c.UserEditable {
			w.Attr(nil, "user_editable", "false")
		}
		if err != nil {
			panic(err)
		}
		w.End(true)
	}
	w.End(false)
	if err := w.Close(); err != nil {
		panic(err)
	}
}

func openProto[T proto.Message](fsys fs.FS, fn string) (T, error) {
	var z T
	msg := reflect.New(reflect.TypeOf(z).Elem()).Interface().(T)
	buf, err := fs.ReadFile(fsys, fn)
	if err != nil {
		return z, fmt.Errorf("read %T from %q: %w", msg, fn, err)
	}
	if err := proto.Unmarshal(buf, msg); err != nil {
		return z, fmt.Errorf("read %T from %q: %w", msg, fn, err)
	}
	return msg, nil
}
