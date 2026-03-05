package cacheredis

import "testing"

func TestConnect(t *testing.T) {
	if rds := NewRedisCache("192.168.1.37:6379", "", 11, 5000); rds == nil {
		t.Fail()
	}
}

func TestSetNotExists(t *testing.T) {
	rds := NewRedisCache("192.168.1.37:6379", "", 11, 5000)
	if rds == nil {
		t.Fail()
	}

	// Set the value.
	if err := rds.SetCond("TestKey", []byte("value1"), 30000, true); err != nil {
		t.Logf("Error: %s", err)
	}

	// Set the value. If the key exists, it should return an error
	if err := rds.SetCond("TestKey", []byte("value2"), 30000, true); err != nil {
		t.Logf("Error: %s", err)
	}
}
