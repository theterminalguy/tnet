package repository

var TimeZoneRepo = map[string][]string{
	"GMT": {
		"GMT",
		"Greenwich Mean Time",
		"GMT",
	},
	"UTC": {
		"UTC",
		"Universal Coordinated Time",
		"GMT",
	},
	"WAT": {
		"WAT",
		"West African Standard Time",
		"GMT+1:00",
	},
	"ECT": {
		"ECT",
		"European Central Time",
		"GMT+1:00",
	},
	"EET": {
		"EET",
		"Eastern European Time",
		"GMT+2:00",
	},
	"ART": {
		"ART",
		"(Arabic) Egypt Standard Time",
		"GMT+2:00",
	},
	"EAT": {
		"EAT",
		"Eastern African Time",
		"GMT+3:00",
	},
	"MET": {
		"MET",
		"Middle East Time",
		"GMT+3:30",
	},
	"NET": {
		"NET",
		"Near East Time",
		"GMT+4:00",
	},
	"PLT": {
		"PLT",
		"Pakistan Lahore Time",
		"GMT+5:00",
	},
	"IST": {
		"IST",
		"India Standard Time",
		"GMT+5:30",
	},
	"BST": {
		"BST",
		"Bangladesh Standard Time",
		"GMT+6:00",
	},
	"VST": {
		"VST",
		"Vietnam Standard Time",
		"GMT+7:00",
	},
	"CTT": {
		"CTT",
		"China Taiwan Time",
		"GMT+8:00",
	},
	"JST": {
		"JST",
		"Japan Standard Time",
		"GMT+9:00",
	},
	"ACT": {
		"ACT",
		"Australia Central Time",
		"GMT+9:30",
	},
	"AET": {
		"AET",
		"Australia Eastern Time",
		"GMT+10:00",
	},
	"SST": {
		"SST",
		"Solomon Standard Time",
		"GMT+11:00",
	},
	"NST": {
		"NST",
		"New Zealand Standard Time",
		"GMT+12:00",
	},
	"MIT": {
		"MIT",
		"Midway Islands Time",
		"GMT-11:00",
	},
	"HST": {
		"HST",
		"Hawaii Standard Time",
		"GMT-10:00",
	},
	"AST": {
		"AST",
		"Alaska Standard Time",
		"GMT-9:00",
	},
	"PST": {
		"PST",
		"Pacific Standard Time",
		"GMT-8:00",
	},
	"PNT": {
		"PNT",
		"Phoenix Standard Time",
		"GMT-7:00",
	},
	"MST": {
		"MST",
		"Mountain Standard Time",
		"GMT-7:00",
	},
	"CST": {
		"CST",
		"Central Standard Time",
		"GMT-6:00",
	},
	"EST": {
		"EST",
		"Eastern Standard Time",
		"GMT-5:00",
	},
	"IET": {
		"IET",
		"Indiana Eastern Standard Time",
		"GMT-5:00",
	},
	"PRT": {
		"PRT",
		"Puerto Rico and US Virgin Islands Time",
		"GMT-4:00",
	},
	"CNT": {
		"CNT",
		"Canada Newfoundland Time",
		"GMT-3:00",
	},
	"AGT": {
		"AGT",
		"Argentina Standard Time",
		"GMT-3:00",
	},
	"BET": {
		"BET",
		"Brazil Eastern Time",
		"GMT-3:00",
	},
	"CAT": {
		"CAT",
		"Central African Time",
		"GMT-1:00",
	},
	"NA": {
		"NA",
		"NA",
		"NA",
	},
}

type TimeZoneData struct {
	ID   int               `json:"id"`
	Data TimezoneDataModel `json:"data"`
}

type TimezoneDataModel struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	RelativeToGMT string `json:"relative_to_gmt"`
}

func ReturnTimezoneData() []TimeZoneData {
	return []TimeZoneData{
		{
			ID: 1,
			Data: TimezoneDataModel{
				"GMT",
				"Greenwich Mean Time",
				"GMT",
			},
		},
		{
			ID: 2,
			Data: TimezoneDataModel{
				"UTC",
				"Universal Coordinated Time",
				"GMT",
			},
		},
		{
			ID: 3,
			Data: TimezoneDataModel{
				"WAT",
				"West African Standard Time",
				"GMT+1:00",
			},
		},
		{
			ID: 4,
			Data: TimezoneDataModel{
				"ECT",
				"European Central Time",
				"GMT+1:00",
			},
		},
		{
			ID: 5,
			Data: TimezoneDataModel{
				"EET",
				"Eastern European Time",
				"GMT+2:00",
			},
		},
		{
			ID: 6,
			Data: TimezoneDataModel{
				"ART",
				"(Arabic) Egypt Standard Time",
				"GMT+2:00",
			},
		},
		{
			ID: 7,
			Data: TimezoneDataModel{
				"EAT",
				"Eastern African Time",
				"GMT+3:00",
			},
		},
		{
			ID: 8,
			Data: TimezoneDataModel{
				"MET",
				"Middle East Time",
				"GMT+3:30",
			},
		},
		{
			ID: 9,
			Data: TimezoneDataModel{
				"NET",
				"Near East Time",
				"GMT+4:00",
			},
		},
		{
			ID: 10,
			Data: TimezoneDataModel{
				"PLT",
				"Pakistan Lahore Time",
				"GMT+5:00",
			},
		},
		{
			ID: 11,
			Data: TimezoneDataModel{
				"IST",
				"India Standard Time",
				"GMT+5:30",
			},
		},
		{
			ID: 12,
			Data: TimezoneDataModel{
				"BST",
				"Bangladesh Standard Time",
				"GMT+6:00",
			},
		},
		{
			ID: 13,
			Data: TimezoneDataModel{
				"VST",
				"Vietnam Standard Time",
				"GMT+7:00",
			},
		},
		{
			ID: 14,
			Data: TimezoneDataModel{
				"CTT",
				"China Taiwan Time",
				"GMT+8:00",
			},
		},
		{
			ID: 15,
			Data: TimezoneDataModel{
				"JST",
				"Japan Standard Time",
				"GMT+9:00",
			},
		},
		{
			ID: 16,
			Data: TimezoneDataModel{
				"ACT",
				"Australia Central Time",
				"GMT+9:30",
			},
		}, {
			ID: 17,
			Data: TimezoneDataModel{
				"AET",
				"Australia Eastern Time",
				"GMT+10:00",
			},
		},
		{
			ID: 18,
			Data: TimezoneDataModel{
				"SST",
				"Solomon Standard Time",
				"GMT+11:00",
			},
		},
		{
			ID: 19,
			Data: TimezoneDataModel{
				"NST",
				"New Zealand Standard Time",
				"GMT+12:00",
			},
		},
		{
			ID: 20,
			Data: TimezoneDataModel{
				"MIT",
				"Midway Islands Time",
				"GMT-11:00",
			},
		},
		{
			ID: 21,
			Data: TimezoneDataModel{
				"HST",
				"Hawaii Standard Time",
				"GMT-10:00",
			},
		},
		{
			ID: 22,
			Data: TimezoneDataModel{
				"AST",
				"Alaska Standard Time",
				"GMT-9:00",
			},
		},
		{
			ID: 23,
			Data: TimezoneDataModel{
				"PST",
				"Pacific Standard Time",
				"GMT-8:00",
			},
		},
		{
			ID: 24,
			Data: TimezoneDataModel{
				"PNT",
				"Phoenix Standard Time",
				"GMT-7:00",
			},
		},
		{
			ID: 25,
			Data: TimezoneDataModel{
				"MST",
				"Mountain Standard Time",
				"GMT-7:00",
			},
		},
		{
			ID: 26,
			Data: TimezoneDataModel{
				"CST",
				"Central Standard Time",
				"GMT-6:00",
			},
		},
		{
			ID: 27,
			Data: TimezoneDataModel{
				"EST",
				"Eastern Standard Time",
				"GMT-5:00",
			},
		},
		{
			ID: 28,
			Data: TimezoneDataModel{
				"IET",
				"Indiana Eastern Standard Time",
				"GMT-5:00",
			},
		},
		{
			ID: 29,
			Data: TimezoneDataModel{
				"PRT",
				"Puerto Rico and US Virgin Islands Time",
				"GMT-4:00",
			},
		},
		{
			ID: 30,
			Data: TimezoneDataModel{
				"CNT",
				"Canada Newfoundland Time",
				"GMT-3:00",
			},
		},
		{
			ID: 31,
			Data: TimezoneDataModel{
				"AGT",
				"Argentina Standard Time",
				"GMT-3:00",
			},
		},
		{
			ID: 32,
			Data: TimezoneDataModel{
				"BET",
				"Brazil Eastern Time",
				"GMT-3:00",
			},
		},
		{
			ID: 33,
			Data: TimezoneDataModel{
				"CAT",
				"Central African Time",
				"GMT-1:00",
			},
		},
	}
}
