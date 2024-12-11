package main

func main() {
	// TODO: implement apn extractor for samsung csc optics
	// each customer.xml contains multiple Network elements with a name and id (like carrierId, but samsung's version)
	// the Connections element contains a ProfileHandle element for each NetworkName referencing a Profile element for each APN type (browser=default,mms,ims,xcap) (like a apn_set_id)
	// the Profile elements contain the apn config
}
