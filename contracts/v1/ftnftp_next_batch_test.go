package v1

import "testing"

func TestFTNFTPObjectValidRejectsNegativeSize(t *testing.T) {
	o := FTNFTPObject{ID:"o1", TenantID:"t1", Bucket:"b", Key:"k", Size:-1}
	if o.Valid() { t.Fatal("negative object size must be invalid") }
}
