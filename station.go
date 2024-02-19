package surfrad

/*
"bon" is the station identifier for Bondville, Illinois
"fpk" is the station identifier for Fort Peck, Montana
"gwn" is the station identifier for Goodwin Creek, Mississippi
"tbl" is the station identifier for Table Mountain, Colorado
"dra" is the station identifier for Desert Rock, Nevada
"psu" is the station identifier for Penn State, Pennsylvania
"sxf" is the station identifier for Sioux Falls, South Dakota
*/

type (
	StationID   [3]rune
	StationName string
)

const (
	StationBondville     StationName = "Bondville"
	StationFortPeck      StationName = "Fort Peck"
	StationGoodwinCreek  StationName = "Goodwin Creek"
	StationTableMountain StationName = "Table Mountain"
	StationDesertRock    StationName = "Desert Rock"
	StationPennState     StationName = "Penn State"
	StationSiouxFalls    StationName = "Sioux Falls"
)

func (sn StationName) String() string {
	return string(sn)
}

func (sn StationName) Valid() bool {
	return ValidateStationName(sn)
}

func (sid StationID) String() string {
	return string(sid[:])
}

func (sid StationID) Valid() bool {
	return ValidateStationID(sid)
}

var (
	StationIDBondville     StationID = [3]rune{'b', 'o', 'n'}
	StationIDFortPeck      StationID = [3]rune{'f', 'p', 'k'}
	StationIDGoodwinCreek  StationID = [3]rune{'g', 'w', 'n'}
	StationIDTableMountain StationID = [3]rune{'t', 'b', 'l'}
	StationIDDesertRock    StationID = [3]rune{'d', 'r', 'a'}
	StationIDPennState     StationID = [3]rune{'p', 's', 'u'}
	StationIDSiouxFalls    StationID = [3]rune{'s', 'x', 'f'}

	StationIDToName = map[StationID]StationName{
		StationIDBondville:     StationBondville,
		StationIDFortPeck:      StationFortPeck,
		StationIDGoodwinCreek:  StationGoodwinCreek,
		StationIDTableMountain: StationTableMountain,
		StationIDDesertRock:    StationDesertRock,
		StationIDPennState:     StationPennState,
		StationIDSiouxFalls:    StationSiouxFalls,
	}
	NameToStationID = map[StationName]StationID{
		StationBondville:     StationIDBondville,
		StationFortPeck:      StationIDFortPeck,
		StationGoodwinCreek:  StationIDGoodwinCreek,
		StationTableMountain: StationIDTableMountain,
		StationDesertRock:    StationIDDesertRock,
		StationPennState:     StationIDPennState,
		StationSiouxFalls:    StationIDSiouxFalls,
	}
)

func ValidateStationID(sid StationID) bool {
	_, ok := StationIDToName[sid]
	return ok
}

func ValidateStationName(sn StationName) bool {
	_, ok := NameToStationID[sn]
	return ok
}

func GetStationID(sn StationName) (StationID, bool) {
	sid, ok := NameToStationID[sn]
	return sid, ok
}

func GetStationName(sid StationID) (StationName, bool) {
	sn, ok := StationIDToName[sid]
	return sn, ok
}
