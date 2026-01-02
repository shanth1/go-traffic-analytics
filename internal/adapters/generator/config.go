package generator

type SimulationProfile string

const (
	ProfileBalanced SimulationProfile = "balanced"
	ProfileSpike    SimulationProfile = "spike"
	ProfileGrowth   SimulationProfile = "growth"
)

type IntRange struct {
	Min int
	Max int
}

type Config struct {
	Seed int64

	UsersCount       int
	CampaignsPerUser IntRange
	LinksPerCampaign IntRange
	ClicksPerLink    IntRange

	HistoryDays int

	ViralLinkProbability float64

	BatchSize int

	PasswordDefault string
}

func DefaultConfig() Config {
	return Config{
		Seed:                 123,
		UsersCount:           3,
		CampaignsPerUser:     IntRange{Min: 2, Max: 5},
		LinksPerCampaign:     IntRange{Min: 10, Max: 20},
		ClicksPerLink:        IntRange{Min: 100, Max: 1000},
		HistoryDays:          90,
		ViralLinkProbability: 0.03,
		BatchSize:            100,
		PasswordDefault:      "password",
	}
}
