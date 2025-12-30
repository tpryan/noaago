package noaago

type ProductType string

const (
	ProductWaterLevel          ProductType = "water_level"
	ProductPredictions         ProductType = "predictions"
	ProductAirTemp             ProductType = "air_temperature"
	ProductWaterTemp           ProductType = "water_temperature"
	ProductWind                ProductType = "wind"
	ProductAirPressure         ProductType = "air_pressure"
	ProductConductivity        ProductType = "conductivity"
	ProductVisibility          ProductType = "visibility"
	ProductHumidity            ProductType = "humidity"
	ProductSalinity            ProductType = "salinity"
	ProductHourlyHeight        ProductType = "hourly_height"
	ProductHighLow             ProductType = "high_low"
	ProductDailyMean           ProductType = "daily_mean"
	ProductMonthlyMean         ProductType = "monthly_mean"
	ProductOneMinuteWaterLevel ProductType = "one_minute_water_level"
)

type Datum string

const (
	DatumMLLW Datum = "MLLW" // Mean Lower Low Water
	DatumMSL  Datum = "MSL"  // Mean Sea Level
	DatumMHW  Datum = "MHW"  // Mean High Water
	DatumMHHW Datum = "MHHW" // Mean Higher High Water
	DatumMLW  Datum = "MLW"  // Mean Low Water
	DatumNAVD Datum = "NAVD" // North American Vertical Datum
	DatumSTND Datum = "STND" // Station Datum
)

type StationType string

const (
	StationTypeWaterLevels StationType = "waterlevels"
	StationTypeCurrents    StationType = "currents"
)

type Units string

const (
	UnitsEnglish Units = "english"
	UnitsMetric  Units = "metric"
)

type TimeZone string

const (
	TimeZoneGMT    TimeZone = "gmt"
	TimeZoneLST    TimeZone = "lst"     // Local Standard Time
	TimeZoneLSTLDT TimeZone = "lst_ldt" // Local Standard/Daylight Time
)
