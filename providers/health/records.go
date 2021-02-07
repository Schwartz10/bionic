package health

import "encoding/xml"

func (p *health) parseRecord(_ *Data, decoder *xml.Decoder, start *xml.StartElement) error {
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