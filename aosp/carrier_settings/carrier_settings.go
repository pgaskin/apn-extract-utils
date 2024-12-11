package carrier_settings

// https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:tools/carrier_settings/proto/carrier_settings.proto

//go:generate protoc --go_out=.  --go_opt=Mcarrier_settings.proto=github.com/pgaskin/apn-extract-utils/aosp/carrier_settings --go_opt=paths=source_relative carrier_settings.proto
