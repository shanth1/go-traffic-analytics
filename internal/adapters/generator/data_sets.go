package generator

type GeoLocation struct {
	CountryCode string
	CountryName string
	Cities      []string
}

type weightedGeo struct {
	Geo    GeoLocation
	Weight int
}

var worldData = []weightedGeo{
	{GeoLocation{"US", "United States", []string{"New York", "Los Angeles", "Chicago", "San Francisco", "Austin"}}, 40},
	{GeoLocation{"DE", "Germany", []string{"Berlin", "Munich", "Hamburg", "Frankfurt"}}, 15},
	{GeoLocation{"GB", "United Kingdom", []string{"London", "Manchester", "Liverpool"}}, 15},
	{GeoLocation{"FR", "France", []string{"Paris", "Lyon", "Marseille"}}, 10},
	{GeoLocation{"JP", "Japan", []string{"Tokyo", "Osaka", "Kyoto"}}, 5},
	{GeoLocation{"BR", "Brazil", []string{"Sao Paulo", "Rio de Janeiro"}}, 5},
	{GeoLocation{"IN", "India", []string{"Mumbai", "Delhi", "Bangalore"}}, 10},
}

type DeviceProfile struct {
	Device  string
	OS      string
	Browser string
}

var deviceProfiles = []DeviceProfile{
	{"Desktop", "Windows", "Chrome"},
	{"Desktop", "Windows", "Firefox"},
	{"Desktop", "Windows", "Edge"},
	{"Desktop", "macOS", "Safari"},
	{"Desktop", "macOS", "Chrome"},
	{"Desktop", "Linux", "Firefox"},
	{"Mobile", "iOS", "Safari"},
	{"Mobile", "Android", "Chrome"},
	{"Tablet", "iOS", "Safari"},
}

var referers = []string{
	"https://google.com",
	"https://facebook.com",
	"https://twitter.com",
	"https://linkedin.com",
	"https://instagram.com",
	"https://t.co",
	"Direct",
}

var campaignNames = []string{
	"Summer Sale 2024", "Black Friday Prep", "Product Hunt Launch",
	"Newsletter Q1", "Instagram Bio Link", "Tech Blog Promo",
	"Webinar Signup", "E-book Download", "Partner Program",
}

var slugWords = []string{
	"go", "get", "buy", "sale", "top", "secret", "deal", "promo", "discount", "join",
}
