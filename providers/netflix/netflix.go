package netflix

import (
	"github.com/shekhirin/bionic-cli/providers/provider"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"path"
)

const name = "netflix"
const tablePrefix = "netflix_"

type netflix struct {
	provider.Database
}

func New(db *gorm.DB) provider.Provider {
	return &netflix{
		Database: provider.NewDatabase(db),
	}
}

func (netflix) Name() string {
	return name
}

func (netflix) TablePrefix() string {
	return tablePrefix
}

func (netflix) Models() []schema.Tabler {
	return []schema.Tabler{
		&ViewingAction{},
		&SubscriptionHistoryItem{},
		&ClickstreamAction{},
		&IndicatedPreference{},
		&InteractiveTitle{},
		&MyListItem{},
		&PlaybackRelatedEvent{},
		&Playtrace{},
		&Rating{},
		&SearchHistoryItem{},
		&Device{},
		&IPAddress{},
		&BillingHistoryItem{},
	}
}

func (p *netflix) ImportFns(inputPath string) ([]provider.ImportFn, error) {
	if !provider.IsPathDir(inputPath) {
		return nil, provider.ErrInputPathShouldBeDirectory
	}

	return []provider.ImportFn{
		provider.NewImportFn(
			"Viewing Activity",
			p.importViewingActivity,
			path.Join(inputPath, "Content_Interaction", "ViewingActivity.csv"),
		),
		provider.NewImportFn(
			"Subscriptions History",
			p.importSubscriptionsHistory,
			path.Join(inputPath, "Account", "SubscriptionHistory.csv"),
		),
		provider.NewImportFn(
			"Clickstream",
			p.importClickstream,
			path.Join(inputPath, "Clickstream", "Clickstream.csv"),
		),
		provider.NewImportFn(
			"Indicated Preferences",
			p.importIndicatedPreferences,
			path.Join(inputPath, "Content_Interaction", "IndicatedPreferences.csv"),
		),
		provider.NewImportFn(
			"Interactive Titles",
			p.importInteractiveTitles,
			path.Join(inputPath, "Content_Interaction", "InteractiveTitles.csv"),
		),
		provider.NewImportFn(
			"My List",
			p.importMyList,
			path.Join(inputPath, "Content_Interaction", "MyList.csv"),
		),
		provider.NewImportFn(
			"Playback Related Events",
			p.importPlaybackRelatedEvents,
			path.Join(inputPath, "Content_Interaction", "PlaybackRelatedEvents.csv"),
		),
		provider.NewImportFn(
			"Ratings",
			p.importRatings,
			path.Join(inputPath, "Content_Interaction", "Ratings.csv"),
		),
		provider.NewImportFn(
			"Search History",
			p.importSearchHistory,
			path.Join(inputPath, "Content_Interaction", "SearchHistory.csv"),
		),
		provider.NewImportFn(
			"Devices",
			p.importDevices,
			path.Join(inputPath, "Devices", "Devices.csv"),
		),
		provider.NewImportFn(
			"IP Addresses",
			p.importIPAddresses,
			path.Join(inputPath, "ip_Addresses", "ipAddresses.csv"),
		),
		provider.NewImportFn(
			"Billing History",
			p.importBillingHistory,
			path.Join(inputPath, "Payment_And_Billing", "BillingHistory.csv"),
		),
	}, nil
}
