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
		Seed:                 12345,
		UsersCount:           10,
		HistoryDays:          30,
		CampaignsPerUser:     struct{ Min, Max int }{Min: 1, Max: 5},
		LinksPerCampaign:     struct{ Min, Max int }{Min: 2, Max: 10},
		ClicksPerLink:        struct{ Min, Max int }{Min: 50, Max: 500},
		ViralLinkProbability: 0.05,
		BatchSize:            1000,
		PasswordDefault:      "password123",
	}
}
