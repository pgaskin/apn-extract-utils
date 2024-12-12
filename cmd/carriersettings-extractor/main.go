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
	"github.com/pgaskin/apn-extract-utils/source/carriersettings"
	"github.com/pgaskin/xmlwriter"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// TODO: make this actually a proper command
// TODO: refactor main logic into new source/carriersettings package
// TODO: compare output with lineage carriersettings-extractor
// TODO: look at logic in https://cs.android.com/android/platform/superproject/main/+/main:tools/carrier_settings/java/CarrierConfigConverterV2.java;bpv=0
// https://android.googlesource.com/platform/packages/apps/CarrierConfig/+/master/src/com/android/carrierconfig/DefaultCarrierConfigService.java
// https://cs.android.com/android/platform/superproject/main/+/main:frameworks/opt/telephony/src/java/com/android/internal/telephony/CarrierResolver.java;drc=be5b10f9022f6e4aeab9c39f50c1e6ac27e19eae;l=1018
// TODO: rewrite this
// TODO: rewrite the carrier id matching logic to actually match properly (the different fields have different levels of precedence, and matching is more than just string matching)

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
		filterNameSuffix              = "" //"_ca"
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

	tier2Settings, err := openProto[*carrier_settings.MultiCarrierSettings](pixelCarrierSettings, "others.pb")
	if err != nil {
		panic(err)
	}

	if debugDumpText {
		buf, _ := txt.Marshal(tier2Settings)
		os.WriteFile(filepath.Join("dbg", "others.textpb"), buf, 0666)
	}

	allSettings := map[string]*carrier_settings.CarrierSettings{} // [canonicalName]
	for _, cs := range tier2Settings.Setting {
		if filterNameSuffix != "" && !strings.HasSuffix(cs.GetCanonicalName(), filterNameSuffix) {
			continue
		}
		slog.Debug("loaded tier 2 carrier settings", "canonical_name", cs.GetCanonicalName())
		allSettings[*cs.CanonicalName] = cs
	}
	slog.Info("loaded tier 2 carrier settings", "total", len(tier2Settings.Setting))

	var tier1Settings int
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
		tier1Settings++
		slog.Debug("loaded settings", "name", name, "canonical_name", carrierSettings.GetCanonicalName())
		if _, ok := allSettings[*carrierSettings.CanonicalName]; ok {
			slog.Warn("replacing tier 2 settings", "name", name, "canonical_name", carrierSettings.GetCanonicalName())
		}
		allSettings[*carrierSettings.CanonicalName] = carrierSettings
		return nil
	}); err != nil {
		panic(err)
	}
	slog.Info("loaded tier 1 carrier settings", "total", tier1Settings)

	carrierMap := map[string][]*carrier_list.CarrierMap{} // [canonicalName]
	for _, c := range carrierList.Entry {
		canonicalName := c.GetCanonicalName()
		if filterNameSuffix != "" && !strings.HasSuffix(canonicalName, filterNameSuffix) {
			continue
		}
		carrierMap[canonicalName] = append(carrierMap[canonicalName], c)
	}
	for _, canonicalName := range slices.Sorted(maps.Keys(allSettings)) {
		if len(carrierMap[canonicalName]) == 0 {
			slog.Error("failed to find carrier_list entry for carrier, dropping", "canonical_name", canonicalName)
			delete(allSettings, canonicalName)
			continue
		}
	}
	slog.Info("mapped carrier_settings to carrier_list entries")

	carrierMapID := map[string][]*carrierid.CarrierId{} // [canonicalName]
	carrierIdMatchedExact := map[int]string{}
	for _, canonicalName := range slices.Sorted(maps.Keys(allSettings)) {
		carrier := carrierMap[canonicalName]
		for _, cs := range carrier {
			for _, wantMatch := range cs.CarrierId {
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
					// TODO: improve this, maybe filter by all instead of one
					other, ok := carrierIdMatchedExact[i]
					if !ok || other != canonicalName {
						carrierIdMatchedExact[i] = canonicalName
					}
					if ok && other != canonicalName {
						slog.Warn("multiple carriersettings carriers matched a single carrierId exactly", "canonical_name", canonicalName, "other_canonical_name", other, "carrier_id", *carrierId.CarrierId[i].CanonicalId)
					}
					if !slices.Contains(carrierMapID[canonicalName], carrierId.CarrierId[i]) {
						carrierMapID[canonicalName] = append(carrierMapID[canonicalName], carrierId.CarrierId[i])
					}
				}
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
		for _, cs := range carrier {
			for _, wantMatch := range cs.CarrierId {
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
					// TODO: improve this
					i := possibleMatches[0]
					other, ok := carrierIdMatchedPLMNOnly[i]
					if !ok || other != canonicalName {
						carrierIdMatchedPLMNOnly[i] = canonicalName
						slog.Warn("added a carrierId match for a carriersettings carrier by the plmn only", "canonical_name", canonicalName, "carrier_id", *carrierId.CarrierId[i].CanonicalId, "name", carrierId.CarrierId[i].GetCarrierName)
					}
					if ok && other != canonicalName {
						slog.Warn("multiple carriersettings carriers which idn't match a single carrierId exactly matched a non-exactly-matched carrierId by the just the PLMN", "canonical_name", canonicalName, "other_canonical_name", other, "carrier_id", *carrierId.CarrierId[i].CanonicalId)
					}
					if !slices.Contains(carrierMapID[canonicalName], carrierId.CarrierId[i]) {
						carrierMapID[canonicalName] = append(carrierMapID[canonicalName], carrierId.CarrierId[i])
					}
				}
			}
		}
	}
	slog.Info("mapped carrier_list entries to carrierId (plus plmn-only matches of remaining carrierId entries)", "have", len(carrierMapID), "missing", len(carrierMap)-len(carrierMapID))

	type ConvertedAPN struct {
		Comment       string
		CanonicalName string
		Setting       apn.Setting
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
			s, err := carriersettings.ConvertAPN(src)
			if err != nil {
				slog.Error("failed to convert apn, skipping", "error", err)
				continue
			}
			if s.EntryName == "" {
				slog.Error("missing apn name, skipping")
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
					for _, cs := range carrier {
						for _, c := range cs.CarrierId {
							tmp, err := carriersettings.WithAPNCarrier(s, c)
							if err != nil {
								panic(err)
							}
							if canonicalID != -1 {
								tmp.CarrierID = canonicalID
							}
							apns = append(apns, ConvertedAPN{
								Comment:       canonicalName,
								CanonicalName: canonicalName,
								Setting:       tmp,
							})
						}
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
