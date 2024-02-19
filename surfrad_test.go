package surfrad

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// TestParseLine - Test parsing functionality with various input scenarios.
func TestParseLine(t *testing.T) {
	debug = true
	cases := []struct {
		name     string
		line     string
		expected Data
		wantErr  bool
	}{
		{
			name: "valid data line",
			line: "2024  48  2 17 23 59 23.983  74.37   136.8 0    28.3 0    49.4 0   126.7 0   320.1 0   289.68 0   289.43 0   396.8 0   288.55 0   288.61 0     9.8 0    62.0 0   111.7 0   -76.7 0    35.0 0    15.1 0    29.0 0     5.1 0   106.8 0   903.6 0",
			expected: Data{
				RawTimestamp: RawEntryTime{
					Year: 2024, Month: 2, Day: 17, JDay: 48,
					Hour: 23, Minute: 59, Decimal: 23.983,
				},
				Timestamp:                         time.Date(2024, time.February, 17, 23, 59, 0, 0, time.UTC),
				SolarZenithAngle:                  74.37,
				DownwellingSolar:                  136.8,
				UpwellingSolar:                    28.3,
				DirectNormalSolar:                 49.4,
				DownwellingDiffuseSolar:           126.7,
				DownwellingIR:                     320.1,
				DownwellingIRCaseTemp:             289.68,
				DownwellingIRDomeTemp:             289.43,
				UpwellingIR:                       396.8,
				UpwellingIRCaseTemp:               288.55,
				UpwellingIRDomeTemp:               288.61,
				GlobalUVB:                         9.8,
				PhotosyntheticallyActiveRadiation: 62.0,
				NetSolar:                          111.7,
				NetIR:                             -76.7,
				TotalNetRadiation:                 35.0,
				TemperatureC:                      15.1,
				RelativeHumidity:                  29.0,
				WindSpeedMetersPerSecond:          5.1,
				WindDirectionDegrees:              106.8,
				BarometricPressure:                903.6,
			},
			wantErr: false,
		},
		{
			name:    "incomplete record",
			line:    "1995 1 1 1 0 0 0.0 0.0", // Truncated line
			wantErr: true,
		},
		{
			name: "missing value",
			line: "1995 1 1 1 0 0 0.0 0.0 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1 -9999.9 1",
			expected: Data{
				RawTimestamp: RawEntryTime{
					Year: 1995, Month: 1, Day: 1, JDay: 1, Hour: 0, Minute: 0, Decimal: 0.0,
				},
				Timestamp:        time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
				SolarZenithAngle: 0.0,
				DownwellingSolar: 0.0,
				// QCDWSolar: 1,
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseLine(strings.Fields(tc.line))
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			t.Log(spew.Sdump(got))
			if tc.expected != (Data{}) && !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("ParseLine() got = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestReadData(t *testing.T) {
	f, err := os.OpenFile("testdata/dra24048.dat", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	data, err := ReadData(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spew.Sdump(data))
}
