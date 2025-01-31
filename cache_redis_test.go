package cacheredis

import "testing"

func TestConnect(t *testing.T) {
	if rds := NewRedisCache("192.168.1.37:6379", "", 11, 5000); rds == nil {
		t.Fail()
	}
}
