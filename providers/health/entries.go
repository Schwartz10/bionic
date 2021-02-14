package health

import (
	"encoding/xml"
	"github.com/BionicTeam/bionic/types"
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Type            string         `xml:"type,attr" gorm:"uniqueIndex:health_entries_key"`
	SourceName      string         `xml:"sourceName,attr"`
	SourceVersion   string         `xml:"sourceVersion,attr"`
	Unit            string         `xml:"unit,attr"`
	CreationDate    types.DateTime `xml:"creationDate,attr" gorm:"uniqueIndex:health_entries_key"`
	StartDate       types.DateTime `xml:"startDate,attr"`
	EndDate         types.DateTime `xml:"endDate,attr"`
	Value           string         `xml:"value,attr"`
	DeviceID        *int
	Device          *Device          `xml:"device,attr"`
	MetadataEntries []MetadataEntry  `xml:"MetadataEntry" gorm:"polymorphic:Parent"`
	BeatsPerMinutes []BeatsPerMinute `xml:"HeartRateVariabilityMetadataList"`
}

func (Entry) TableName() string {
	return tablePrefix + "entries"
}

func (e Entry) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"type":          e.Type,
		"creation_date": e.CreationDate,
	}
}

type BeatsPerMinute struct {
	gorm.Model
	EntryID uint   `gorm:"uniqueIndex:health_beats_per_minutes_key"`
	BPM     int    `xml:"bpm,attr"`
	Time    string `xml:"time,attr" gorm:"uniqueIndex:health_beats_per_minutes_key"`
}

func (BeatsPerMinute) TableName() string {
	return tablePrefix + "beats_per_minutes"
}

func (bpm BeatsPerMinute) Constraints() map[string]interface{} {
	return map[string]interface{}{
		"entry_id": bpm.EntryID,
		"time":     bpm.Time,
	}
}

func (e *Entry) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias Entry

	var data struct {
		Alias
		HeartRateVariabilityMetadataList struct {
			InstantaneousBeatsPerMinute []BeatsPerMinute `xml:"InstantaneousBeatsPerMinute"`
		} `xml:"HeartRateVariabilityMetadataList"`
	}

	if err := decoder.DecodeElement(&data, &start); err != nil {
		return err
	}

	*e = Entry(data.Alias)

	e.BeatsPerMinutes = data.HeartRateVariabilityMetadataList.InstantaneousBeatsPerMinute

	return nil
}

func (p *health) parseRecord(_ *DataExport, decoder *xml.Decoder, start *xml.StartElement) error {
	var entry Entry

	if err := decoder.DecodeElement(&entry, start); err != nil {
		return err
	}

	err := p.DB().
		Find(&entry, entry.Constraints()).
		Error
	if err != nil {
		return err
	}

	if entry.Device != nil {
		err = p.DB().
			FirstOrCreate(entry.Device, entry.Device.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range entry.MetadataEntries {
		metadataEntry := &entry.MetadataEntries[i]

		metadataEntry.ParentID = entry.ID
		metadataEntry.ParentType = entry.TableName()

		err = p.DB().
			FirstOrCreate(metadataEntry, metadataEntry.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	for i := range entry.BeatsPerMinutes {
		beatsPerMinute := &entry.BeatsPerMinutes[i]

		beatsPerMinute.EntryID = entry.ID

		err = p.DB().
			FirstOrCreate(beatsPerMinute, beatsPerMinute.Constraints()).
			Error
		if err != nil {
			return err
		}
	}

	return p.DB().
		FirstOrCreate(&entry, entry.Constraints()).
		Error
}