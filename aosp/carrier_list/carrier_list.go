package carrier_list

// https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:tools/carrier_settings/proto/carrier_list.proto

//go:generate protoc --go_out=.  --go_opt=Mcarrier_list.proto=github.com/pgaskin/apn-extract-utils/aosp/carrier_list --go_opt=paths=source_relative carrier_list.proto
