package riotapi

type Region string

const (
	Brasil            Region = "br1"
	EuropeNordicEast  Region = "eun1"
	EuropeWest        Region = "euw1"
	Japan             Region = "jp1"
	Korea             Region = "kr"
	LatinAmericaNorth Region = "la1"
	LatinAmericaSouth Region = "la2"
	MiddleEast        Region = "me1"
	NorthAmerica      Region = "na1"
	Oceania           Region = "oc1"
	PBE               Region = "pbe1"
	Russia            Region = "ru"
	SouthEastAsia     Region = "sg2"
	Turkey            Region = "tr1"
	Taiwan            Region = "tw2"
	Vietnam           Region = "vn2"
)

var Regions = []Region{
	Brasil,
	EuropeNordicEast,
	EuropeWest,
	Japan,
	Korea,
	LatinAmericaNorth,
	LatinAmericaSouth,
	MiddleEast,
	NorthAmerica,
	Oceania,
	PBE,
	Russia,
	SouthEastAsia,
	Turkey,
	Taiwan,
	Vietnam,
}

type Area string

const (
	Americas Area = "americas"
	Asia     Area = "asia"
	Europe   Area = "europe"
	SEA      Area = "sea"
)

var regionToArea = map[Region]Area{
	Brasil:            Americas,
	EuropeNordicEast:  Europe,
	EuropeWest:        Europe,
	Japan:             Asia,
	Korea:             Asia,
	LatinAmericaNorth: Americas,
	LatinAmericaSouth: Americas,
	MiddleEast:        Europe,
	NorthAmerica:      Americas,
	Oceania:           SEA,
	Russia:            Europe,
	SouthEastAsia:     SEA,
	Turkey:            Europe,
	Taiwan:            SEA,
	Vietnam:           SEA,
}
