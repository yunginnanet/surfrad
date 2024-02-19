package surfrad

import (
	"testing"
)

func TestValidateStationID(t *testing.T) {
	cases := []struct {
		name     string
		sid      StationID
		expected bool
	}{
		{"Valid StationID Bondville", StationIDBondville, true},
		{"Valid StationID Fort Peck", StationIDFortPeck, true},
		{"Invalid StationID", StationID{'x', 'y', 'z'}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateStationID(tc.sid)
			if result != tc.expected {
				t.Errorf("ValidateStationID(%q) == %t, expected %t", tc.sid, result, tc.expected)
			}
		})
	}
}

func TestValidateStationName(t *testing.T) {
	cases := []struct {
		name     string
		sn       StationName
		expected bool
	}{
		{"Valid StationName Bondville", StationBondville, true},
		{"Valid StationName Fort Peck", StationFortPeck, true},
		{"Invalid StationName", StationName("Nowhere, Narnia"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateStationName(tc.sn)
			if result != tc.expected {
				t.Errorf("ValidateStationName(%q) == %t, expected %t", tc.sn, result, tc.expected)
			}
		})
	}
}

func TestGetStationID(t *testing.T) {
	cases := []struct {
		name       string
		sn         StationName
		expectedID StationID
		expectedOK bool
	}{
		{"Get ID for Bondville", StationBondville, StationIDBondville, true},
		{"Get ID for Invalid", StationName("Nowhere, Narnia"), StationID{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, ok := GetStationID(tc.sn)
			if !ok && tc.expectedOK {
				t.Errorf("GetStationID(%q) was expected to succeed but did not", tc.sn)
			}
			if ok && id != tc.expectedID {
				t.Errorf("GetStationID(%q) == %q, expected %q", tc.sn, id, tc.expectedID)
			}
		})
	}
}

func TestGetStationName(t *testing.T) {
	cases := []struct {
		name         string
		sid          StationID
		expectedName StationName
		expectedOK   bool
	}{
		{"Get Name for Bondville", StationIDBondville, StationBondville, true},
		{"Get Name for Invalid", StationID{'x', 'y', 'z'}, "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name, ok := GetStationName(tc.sid)
			if !ok && tc.expectedOK {
				t.Errorf("GetStationName(%q) was expected to succeed but did not", tc.sid)
			}
			if ok && name != tc.expectedName {
				t.Errorf("GetStationName(%q) == %q, expected %q", tc.sid, name, tc.expectedName)
			}
		})
	}
}
