package surfrad

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation int     `json:"elevation"`
}

type Station struct {
	StationName StationName `json:"station_name"`
	LocatedAt   Location    `json:"located_at"`

	Version int `json:"version"`

	Entries []Data `json:"entries"`
}

type RawEntryTime struct {
	Year    int     `json:"year"`
	Month   int     `json:"month"`
	Day     int     `json:"day"`
	JDay    int     `json:"jday"`
	Hour    int     `json:"hour"`
	Minute  int     `json:"minute"`
	Decimal float64 `json:"decimal_time"` // (hour.decimalminutes, e.g., 23.5 = 2330)
}

/*
The variables, their data type, and description are given below:

station_name	character	station name, e. g., Goodwin Creek
latitude		real	latitude in decimal degrees (e. g., 40.80)
longitude		real	longitude in decimal degrees (e. g., 105.12)
elevation		integer	elevation above sea level in meters
year			integer	year, i.e., 1995
jday			integer	Julian day (1 through 365 [or 366])
month			integer	number of the month (1-12)
day			integer	day of the month(1-31)
hour			integer	hour of the day (0-23)
min			integer	minute of the hour (0-59)
dt			real	decimal time (hour.decimalminutes, e.g., 23.5 = 2330)
zen			real	solar zenith angle (degrees)
dw_solar		real	downwelling global solar (Watts m^-2)
uw_solar		real	upwelling global solar (Watts m^-2)
direct_n		real	direct-normal solar (Watts m^-2)
diffuse		real	downwelling diffuse solar (Watts m^-2)
dw_ir			real	downwelling thermal infrared (Watts m^-2)
dw_casetemp		real	downwelling IR case temp. (K)
dw_dometemp		real	downwelling IR dome temp. (K)
uw_ir			real	upwelling thermal infrared (Watts m^-2)
uw_casetemp		real	upwelling IR case temp. (K)
uw_dometemp		real	upwelling IR dome temp. (K)
uvb			real	global UVB (milliWatts m^-2)
par			real	photosynthetically active radiation (Watts m^-2)
netsolar		real	net solar (dw_solar - uw_solar) (Watts m^-2)
netir			real	net infrared (dw_ir - uw_ir) (Watts m^-2)
totalnet		real	net radiation (netsolar+netir) (Watts m^-2)
temp			real	10-meter air temperature (?C)
rh			real	relative humidity (%)
windspd		real	wind speed (ms^-1)
winddir		real	wind direction (degrees, clockwise from north)
pressure		real	station pressure (mb)
*/

type Data struct {
	RawTimestamp RawEntryTime `json:"raw_time_data"`
	Timestamp    time.Time    `json:"timestamp"`

	// Solar Radiation
	SolarZenithAngle                  float64 `json:"solar_zenith_angle,omitempty"`
	DownwellingSolar                  float64 `json:"downwelling_solar,omitempty"`
	UpwellingSolar                    float64 `json:"upwelling_solar,omitempty"`
	DirectNormalSolar                 float64 `json:"direct_normal_solar,omitempty"`
	DownwellingDiffuseSolar           float64 `json:"downwelling_diffuse_solar,omitempty"`
	DownwellingIR                     float64 `json:"downwelling_ir,omitempty"`
	DownwellingIRCaseTemp             float64 `json:"downwelling_ir_case_temp,omitempty"`
	DownwellingIRDomeTemp             float64 `json:"downwelling_ir_dome_temp,omitempty"`
	UpwellingIR                       float64 `json:"upwelling_ir,omitempty"`
	UpwellingIRCaseTemp               float64 `json:"upwelling_ir_case_temp,omitempty"`
	UpwellingIRDomeTemp               float64 `json:"upwelling_ir_dome_temp,omitempty"`
	GlobalUVB                         float64 `json:"global_uvb,omitempty"`
	PhotosyntheticallyActiveRadiation float64 `json:"photosynthetically_active_radiation,omitempty"`
	NetSolar                          float64 `json:"net_solar,omitempty"`
	NetIR                             float64 `json:"net_ir,omitempty"`
	TotalNetRadiation                 float64 `json:"total_net,omitempty"`

	TemperatureC             float64 `json:"temperature,omitempty"` // celcius
	RelativeHumidity         float64 `json:"relative_humidity,omitempty"`
	WindSpeedMetersPerSecond float64 `json:"wind_speed,omitempty"`          // m/s
	WindDirectionDegrees     float64 `json:"wind_direction,omitempty"`      // degrees, clockwise from north
	BarometricPressure       float64 `json:"barometric_pressure,omitempty"` // mb

}

//goland:noinspection GoMixedReceiverTypes
func (s *Station) ParseHeader(headerInfo []string) (error, bool) {
	var errs []error

	if len(headerInfo) < 6 {
		errs = append(errs, fmt.Errorf("invalid header length: %v", headerInfo))
	}

	if len(headerInfo) == 0 {
		return errors.Join(errs...), false
	}

	var err error

	if s.LocatedAt.Latitude, err = strconv.ParseFloat(headerInfo[0], 64); err != nil {
		errs = append(errs, fmt.Errorf("error parsing latitude: %v", err))
	}

	if len(headerInfo) == 1 {
		return errors.Join(errs...), false
	}

	if s.LocatedAt.Longitude, err = strconv.ParseFloat(headerInfo[1], 64); err != nil {
		errs = append(errs, fmt.Errorf("error parsing longitude: %v", err))
	}

	if len(headerInfo) == 2 {
		return errors.Join(errs...), false
	}

	if s.LocatedAt.Elevation, err = strconv.Atoi(headerInfo[2]); err != nil {
		errs = append(errs, fmt.Errorf("error parsing elevation: %v", err))
	}

	if len(headerInfo) < 6 {
		return errors.Join(errs...), false
	}

	if s.Version, err = strconv.Atoi(headerInfo[5]); err != nil {
		errs = append(errs, fmt.Errorf("error parsing version: %v", err))
	}

	// we don't mind if we're missing the version

	return errors.Join(errs...), true
}

//goland:noinspection GoMixedReceiverTypes
func (s Station) Len() int {
	return len(s.Entries)
}

func ReadData(r io.Reader) (Station, error) {
	scanner := bufio.NewScanner(r)

	station := new(Station)
	var errs []error

	if scanner.Scan() {
		station.StationName = StationName(strings.TrimSpace(scanner.Text()))
		if !station.StationName.Valid() {
			errs = append(errs, fmt.Errorf("invalid or unknown station name: %s", station.StationName))
		}
	}

	if scanner.Scan() {
		err, ok := station.ParseHeader(strings.Fields(scanner.Text()))
		if err != nil {
			errs = append(errs, err)
		}
		if !ok {
			return *station, errors.Join(errs...)
		}
	}

	debugPrint(
		"Station Name: %s, Latitude: %.2f, Longitude: %.2f, Elevation: %d, Version: %d\n",
		station.StationName, station.LocatedAt.Latitude,
		station.LocatedAt.Longitude, station.LocatedAt.Elevation,
		station.Version,
	)

	lineNo := 0

	for scanner.Scan() {
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			errs = append(errs, err)
			return *station, errors.Join(errs...)
		}
		lineNo++
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 29 {
			errs = append(errs, fmt.Errorf("incomplete record on line: %d", lineNo))
			continue // skip incomplete records
		}

		record, err := ParseLine(fields)
		if err != nil {
			debugPrint("error parsing line: %v\n", err)
			debugPrint("line: %s\n", line)
			errs = append(errs, err)
			continue
		}

		station.Entries = append(station.Entries, record)

		debugPrint("parsed entry: %v\n", record)
	}

	debugPrint("processed %d entries\n", len(station.Entries))

	return *station, errors.Join(errs...)
}

func (d *Data) ParseTimestamp(fields []string) error {
	if len(fields) < 7 {
		return fmt.Errorf("incomplete timestamp: %v", fields)
	}

	var err error

	// Assuming fields are in the correct order as per the data structure
	if d.RawTimestamp.Year, err = strconv.Atoi(fields[0]); err != nil {
		return err
	}
	if d.RawTimestamp.JDay, err = strconv.Atoi(fields[1]); err != nil {
		return err
	}
	if d.RawTimestamp.Month, err = strconv.Atoi(fields[2]); err != nil {
		return err
	}
	if d.RawTimestamp.Day, err = strconv.Atoi(fields[3]); err != nil {
		return err
	}
	if d.RawTimestamp.Hour, err = strconv.Atoi(fields[4]); err != nil {
		return err
	}
	if d.RawTimestamp.Minute, err = strconv.Atoi(fields[5]); err != nil {
		return err
	}
	d.RawTimestamp.Decimal = parseFloat(fields[6])

	d.Timestamp = time.Date(
		d.RawTimestamp.Year, time.Month(d.RawTimestamp.Month), d.RawTimestamp.Day,
		d.RawTimestamp.Hour, d.RawTimestamp.Minute, 0, 0, time.UTC,
	)

	return err
}

func (d *Data) OmitInvalidOrMissing() {
	// new idea: use reflection to iterate over the fields and set them to 0 if they are -9999.9

	count := reflect.ValueOf(d).Elem().NumField()
	timeType := reflect.TypeOf(time.Time{})
	for i := 0; i < count; i++ {
		field := reflect.ValueOf(d).Elem().Field(i)
		if field.Type().Kind() == reflect.Float64 && field.Float() == -9999.9 {
			field.SetZero()
		}
		if field.Type().Kind() == reflect.Int && field.Int() == -9999 {
			field.SetZero()
		}
		if field.Type().Kind() == timeType.Kind() {
			switch field.Interface().(type) {
			case time.Time:
				if field.Interface().(time.Time).IsZero() {
					field.Set(reflect.Zero(timeType))
				}
			}
		}
	}
}

func ParseLine(fields []string) (Data, error) {
	var data = new(Data)
	var err error

	if err = data.ParseTimestamp(fields); err != nil {
		return *data, err
	}

	for i, field := range fields {
		switch i {
		case 7:
			data.SolarZenithAngle = parseFloat(field)
		case 8:
			data.DownwellingSolar = parseFloat(field)
		case 9:
			// data.QCDWSolar, _ = strconv.Atoi(field)
		case 10:
			data.UpwellingSolar = parseFloat(field)
		case 11:
			// data.QCUWSolar, _ = strconv.Atoi(field)
		case 12:
			data.DirectNormalSolar = parseFloat(field)
		case 13:
			// data.QCDirectN, _ = strconv.Atoi(field)
		case 14:
			data.DownwellingDiffuseSolar = parseFloat(field)
		case 15:
			// data.QCDiffuse, _ = strconv.Atoi(field)
		case 16:
			data.DownwellingIR = parseFloat(field)
		case 17:
			// data.QCDWIR, _ = strconv.Atoi(field)
		case 18:
			data.DownwellingIRCaseTemp = parseFloat(field)
		case 19:
			// data.QCDWCasetemp, _ = strconv.Atoi(field)
		case 20:
			data.DownwellingIRDomeTemp = parseFloat(field)
		case 21:
			// data.QCDWDometemp, _ = strconv.Atoi(field)
		case 22:
			data.UpwellingIR = parseFloat(field)
		case 23:
			// data.QCUWIR, _ = strconv.Atoi(field)
		case 24:
			data.UpwellingIRCaseTemp = parseFloat(field)
		case 25:
			// data.QCUWCasetemp, _ = strconv.Atoi(field)
		case 26:
			data.UpwellingIRDomeTemp = parseFloat(field)
		case 27:
			// data.QCUWDometemp, _ = strconv.Atoi(field)
		case 28:
			data.GlobalUVB = parseFloat(field)
		case 29:
			// data.QCUVB, _ = strconv.Atoi(field)
		case 30:
			data.PhotosyntheticallyActiveRadiation = parseFloat(field)
		case 31:
			// data.QCPAR, _ = strconv.Atoi(field)
		case 32:
			data.NetSolar = parseFloat(field)
		case 33:
			// data.QCNetSolar, _ = strconv.Atoi(field)
		case 34:
			data.NetIR = parseFloat(field)
		case 35:
			// data.QCNetIR, _ = strconv.Atoi(field)
		case 36:
			data.TotalNetRadiation = parseFloat(field)
		case 37:
			// data.QCTotalNet, _ = strconv.Atoi(field)
		case 38:
			data.TemperatureC = parseFloat(field)
		case 39:
			// data.QCTemp, _ = strconv.Atoi(field)
		case 40:
			data.RelativeHumidity = parseFloat(field)
		case 41:
			// data.QCRH, _ = strconv.Atoi(field)
		case 42:
			data.WindSpeedMetersPerSecond = parseFloat(field)
		case 43:
			// data.QCWindSpd, _ = strconv.Atoi(field)
		case 44:
			data.WindDirectionDegrees = parseFloat(field)
		case 45:
			// data.QCWindDir, _ = strconv.Atoi(field)
		case 46:
			data.BarometricPressure = parseFloat(field)
		case 47:
			// data.QCPressure, _ = strconv.Atoi(field)
		default:
			//
		}
	}

	if len(fields) < 47 {
		err = fmt.Errorf("incomplete record: %v", fields)
	}

	data.OmitInvalidOrMissing()

	return *data, err
}

func parseFloat(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return value
}
