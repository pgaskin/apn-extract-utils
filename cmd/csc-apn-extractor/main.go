package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// TODO: implement apn extractor for samsung csc optics
	// each customer.xml contains multiple Network elements with a name and id (like carrierId, but samsung's version)
	// the Connections element contains a ProfileHandle element for each NetworkName referencing a Profile element for each APN type (browser=default,mms,ims,xcap) (like a apn_set_id)
	// the Profile elements contain the apn config
	var (
		optics = os.DirFS("/home/pg/srctest/carriersettings/S711WVLS5CXJ1_S711WVLS5CXJ1_MQB86705389/csc/optics")
	)

	srcs, err := fs.Glob(optics, "configs/carriers/*/conf/customer.xml")
	if err != nil {
		slog.Error("failed to find cscs", "error", err)
		os.Exit(1)
	}
	if len(srcs) == 0 {
		slog.Error("no cscs found")
		os.Exit(1)
	}

	for i, src := range srcs {
		slog := slog.With("i", i)
		if err := func() (err error) {
			defer func() {
				if x := recover(); x != nil {
					if e, ok := x.(error); ok {
						err = e
					} else {
						err = errors.New(fmt.Sprint(x))
					}
					var buf [2048]byte
					n := runtime.Stack(buf[:], false)
					err = fmt.Errorf("%w\n%s", err, string(buf[:n]))
				}
			}() // TODO: more validation and error checking so this isn't necessary

			buf, err := fs.ReadFile(optics, src)
			if err != nil {
				return err
			}
			doc := etree.NewDocument()
			if err := doc.ReadFromBytes(buf); err != nil {
				return err
			}

			type GeneralInfo struct {
				Version string
				Country string
				Region  string
				CSC     string
			}
			generalInfo := GeneralInfo{
				Version: doc.FindElement("/CustomerData/GeneralInfo/CSCEdition").Text(),
				Country: doc.FindElement("/CustomerData/GeneralInfo/CountryISO").Text(),
				Region:  doc.FindElement("/CustomerData/GeneralInfo/Region").Text(),
				CSC:     doc.FindElement("/CustomerData/GeneralInfo/SalesCode").Text(),
			}
			if exp := path.Base(path.Dir(path.Dir(src))); generalInfo.CSC != exp {
				return fmt.Errorf("expected csc %q, got %q", exp, generalInfo.CSC)
			}
			slog.Info("loaded csc", "version", generalInfo.Version, "country", generalInfo.Country, "region", generalInfo.Region, "csc", generalInfo.CSC)

			// TODO: is there a better way to do this?
			if generalInfo.CSC == "XAC" {
				slog.Warn("skipping generic csc")
				return nil
			}

			type NetworkInfo struct {
				MCCMNC string
				Name   string // TODO: is this unique like google's canonicalName?
				NWID   string // essentially samsung's equivalent of carrier_id

				Subset struct {
					Type string
					Code string
				}
			}
			var networkInfos []NetworkInfo
			for _, x := range doc.FindElements("/CustomerData/GeneralInfo/NetworkInfo") {
				if false {
					var b bytes.Buffer
					x.WriteTo(&b, &etree.WriteSettings{})
					b.WriteByte('\n')
					b.WriteTo(os.Stdout)
				}
				networkInfo := NetworkInfo{
					MCCMNC: x.SelectElement("MCCMNC").Text(),
					Name:   x.SelectElement("NetworkName").Text(),
					NWID:   x.SelectElement("NWID").Text(),
				}
				if y := x.SelectElement("SubsetCode"); y != nil {
					networkInfo.Subset.Code = y.Text()
					networkInfo.Subset.Type = textOrEmpty(x.SelectElement("CodeType"))
				}
				networkInfos = append(networkInfos, networkInfo)
				slog.Info("... network", "mccmnc", networkInfo.MCCMNC, "name", networkInfo.Name, "nwid", networkInfo.NWID, "subset_type", networkInfo.Subset.Type, "subset_code", networkInfo.Subset.Code)
			}

			type Profile struct {
				NetworkName          string
				ProfileName          string
				ProfileType          string
				IpVersion            string
				Editable             string
				EnableStatus         string
				URL                  string
				Auth                 string
				Bearer               string
				Protocol             string
				MTUSize              string
				RoamingIpVersion     string
				Selectable           string
				HiddenStatus         string
				Proxy_EnableFlag     string
				Proxy_ServAddr       string
				Proxy_Port           string
				ProfileId            string
				PSparam_APN          string
				PSparam_TrafficClass string
			}
			var profiles []Profile
			for _, x := range doc.FindElements("/CustomerData/Settings/Connections/Profile") {
				if false {
					var b bytes.Buffer
					x.WriteTo(&b, &etree.WriteSettings{})
					b.WriteByte('\n')
					b.WriteTo(os.Stdout)
				}
				var profile Profile
				for p, c := range func(yield func(string, *etree.Element) bool) {
					for _, c := range x.ChildElements() {
						if c.Tag == "PSparam" || c.Tag == "Proxy" {
							for _, cc := range c.ChildElements() {
								if !yield(c.Tag+"_", cc) {
									return
								}
							}
							continue
						}
						if len(c.ChildElements()) != 0 {
							panic(fmt.Errorf("%q has children", c.Tag))
						}
						if !yield("", c) {
							return
						}
					}
				} {
					v := reflect.ValueOf(&profile).Elem().FieldByName(p + c.Tag)
					if v == (reflect.Value{}) {
						panic(fmt.Errorf("unknown field %q=%q", p+c.Tag, c.Text()))
					}
					v.Set(reflect.ValueOf(c.Text()))
				}
				profiles = append(profiles, profile)
			}

			type ProfileHandle struct {
				NetworkName          string
				ProfBrowser          string
				ProfMMS              string
				ProfIMS              string
				ProfXCAP             string
				ProfEpdgXCAP         string
				ProfEpdgMMS          string
				ProfIntSharing       string
				ProfEmergencyIMSCall string
			}
			var profileHandles []ProfileHandle
			for _, x := range doc.FindElements("/CustomerData/Settings/Connections/ProfileHandle") {
				profileHandle := ProfileHandle{
					NetworkName:          textOrEmpty(x.SelectElement("NetworkName")),
					ProfBrowser:          textOrEmpty(x.SelectElement("ProfBrowser")),
					ProfMMS:              textOrEmpty(x.SelectElement("ProfMMS")),
					ProfIMS:              textOrEmpty(x.SelectElement("ProfIMS")),
					ProfXCAP:             textOrEmpty(x.SelectElement("ProfXCAP")),
					ProfEpdgXCAP:         textOrEmpty(x.SelectElement("ProfEpdgXCAP")),
					ProfEpdgMMS:          textOrEmpty(x.SelectElement("ProfEpdgMMS")),
					ProfIntSharing:       textOrEmpty(x.SelectElement("ProfIntSharing")),
					ProfEmergencyIMSCall: textOrEmpty(x.SelectElement("ProfEmergencyIMSCall")),
				}
				if exp, act := mustInt(x.SelectElement("NbNetProfile").Text()), countNonEmptyUnique(
					profileHandle.ProfBrowser,
					profileHandle.ProfMMS,
					profileHandle.ProfIMS,
					profileHandle.ProfXCAP,
					profileHandle.ProfEpdgXCAP,
					profileHandle.ProfEpdgMMS,
					profileHandle.ProfIntSharing,
					profileHandle.ProfEmergencyIMSCall,
				); exp != act {
					if true {
						var b bytes.Buffer
						x.WriteTo(&b, &etree.WriteSettings{})
						b.WriteByte('\n')
						b.WriteTo(os.Stdout)
					}
					return fmt.Errorf("expected %d profile references, got %d (do we not handle some?)", exp, act)
				}
				profileHandles = append(profileHandles, profileHandle)
			}

			for _, ph := range profileHandles {
				slog.Info("... ... handle", "name", ph.NetworkName)
				ref := map[string]Profile{}
				for i := range reflect.TypeOf(ph).NumField() {
					if kind, ok := strings.CutPrefix(reflect.TypeOf(ph).FieldByIndex([]int{i}).Name, "Prof"); ok {
						if name := reflect.ValueOf(ph).FieldByIndex([]int{i}).String(); name != "" {
							idx := slices.IndexFunc(profiles, func(e Profile) bool {
								return e.ProfileName == name
							})
							if idx == -1 {
								if true {
									e := json.NewEncoder(os.Stdout)
									e.SetIndent("", "  ")
									e.Encode(profiles)
								}
								return fmt.Errorf("failed to resolve profile %s:%q for %#v", kind, name, ph)
							}
							ref[name] = profiles[idx]
						} else {
							slog.Warn("skipping blank profile ref", "kind", kind)
						}
					}
				}
				if true {
					e := json.NewEncoder(os.Stdout)
					e.SetIndent("", "  ")
					e.Encode(ref)
				}

				// TODO: convert to apns
				//
				// bearer is comma-separated ps,iwlan,lte (and others?)
				// bearer ps (packet-switched), assume 1|2|3|9|10|11|13|14|15|18|20?
				//
				// TODO: maybe see if I can find the APK which handles CSCs so I can reverse it without guessing at the params?
				//
				// yes, it's in SecTelephonyProvider.apk/com.android.telephony.TelephonyProvider.parse...
				// skimming it, it seems like they have a function which outputs standard apns-conf xml, so I might see if it works
				// it writes output to /data/user_de/0/com.android.providers.telephony/databases/apninfo.xml?
				// it matches using salescode though, so I'll have to see if it's possible to map that to an aosp carrier_id in a sane way
			}

			return nil
		}(); err != nil {
			slog.Error("failed to process csc", "path", src, "error", err)
			os.Exit(1)
		}
	}
}

func textOrEmpty(e *etree.Element) string {
	if e == nil {
		return ""
	}
	return e.Text()
}

func countNonEmptyUnique(s ...string) int {
	m := map[string]struct{}{}
	for _, s := range s {
		if s != "" {
			m[s] = struct{}{}
		}
	}
	return len(m)
}

func countNonEmpty(s ...string) int {
	var n int
	for _, s := range s {
		if s != "" {
			n++
		}
	}
	return n
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}
