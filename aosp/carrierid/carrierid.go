package carrierid

// https://cs.android.com/android/platform/superproject/+/android14-qpr3-release:packages/providers/TelephonyProvider/proto/src/carrierId.proto;l=1?q=carrierId.proto&sq=&ss=android%2Fplatform%2Fsuperproject

//go:generate protoc --go_out=.  --go_opt=McarrierId.proto=github.com/pgaskin/apn-extract-utils/aosp/carrierid --go_opt=paths=source_relative carrierId.proto
