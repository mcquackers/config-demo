package repo

import (
	"fmt"

	"github.com/mcquackers/config-demo/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func makeUpsertConfigFilter(configLevel entities.ConfigLevel, corporateID, venueID, vendorID string) (bson.M, error) {
	if configLevel == entities.CONFIG_LEVEL_UNSPECIFIED {
		return nil, fmt.Errorf("invalid config level: %d", configLevel)
	}
	filter := bson.M{
		"config_level": configLevel,
	}

	if configLevel >= entities.CONFIG_LEVEL_CORPORATE {
		filter["corporate_id"] = corporateID
	}
	if configLevel >= entities.CONFIG_LEVEL_VENUE {
		filter["venue_id"] = venueID
	}
	if configLevel >= entities.CONFIG_LEVEL_VENDOR {
		filter["vendor_id"] = vendorID
	}

	return filter, nil
}
func makeConfigPipeline(corporateID, venueID, vendorID string, configType entities.ConfigType) mongo.Pipeline {
	switch configType {
	case entities.CONFIG_TYPE_UNSPECIFIED: // || !configType.Valid()
		return mongo.Pipeline{}
	case entities.CONFIG_TYPE_FULL:
		return makeMainConfigPipeline(corporateID, venueID, vendorID)
	default:
		return makeUnderlyingConfigPipeline(corporateID, venueID, vendorID, configType)
	}
}
func makeMainConfigPipeline(corporateID, venueID, vendorID string) mongo.Pipeline {
	return mongo.Pipeline{
		{
			{
				Key: "$match",
				Value: bson.M{
					"corporate_id": corporateID,
					"venue_id":     venueID,
					"vendor_id":    vendorID,
				},
			},
		},

		{
			{
				Key:   "$limit",
				Value: 1,
			},
		},
	}
}
func makeUnderlyingConfigPipeline(corporateID, venueID, vendorID string, configType entities.ConfigType) mongo.Pipeline {
	return append(makeMainConfigPipeline(corporateID, venueID, vendorID), makeReplaceRootStage(configType))
}

func makeReplaceRootStage(configType entities.ConfigType) bson.D {
	return bson.D{
		{
			Key: "$replaceRoot",
			Value: bson.M{
				"newRoot": fmt.Sprintf("$%s", configType.String()),
			},
		},
	}
}

func makeGetActiveConfigPipeline(configLevel entities.ConfigLevel, corporateID, venueID, vendorID string, configType entities.ConfigType) mongo.Pipeline {
	return mongo.Pipeline{
		makeGetActiveConfigMatch(configLevel, corporateID, venueID, vendorID, configType),
		makeGetActiveConfigSort(),
		makeGetActiveConfigLimit(),
		makeReplaceRootStage(configType),
	}
}

func makeGetActiveConfigMatch(configLevel entities.ConfigLevel, corporateID, venueID, vendorID string, configType entities.ConfigType) bson.D {
	configMatch := bson.D{
		{
			Key: "$match",
			Value: bson.D{
				{
					"config_level",
					bson.D{
						{
							"$lte",
							configLevel,
						},
					},
				},
				{
					"$or",
					bson.A{
						makeCorporateConfigQuery(corporateID, configType),
						makeVenueConfigQuery(corporateID, venueID, configType),
						makeVendorConfigQuery(corporateID, venueID, vendorID, configType),
					},
				},
			},
		},
	}

	return configMatch
}

func makeGetActiveConfigSort() bson.D {
	return bson.D{{
		Key: "$sort",
		Value: bson.D{
			{
				"vendor_id", -1,
			},
			{
				"venue_id", -1,
			},
		},
	},
	}
}

func makeGetActiveConfigLimit() bson.D {
	return bson.D{{
		Key:   "$limit",
		Value: 1,
	},
	}
}

func makeCorporateConfigQuery(corporateID string, configType entities.ConfigType) bson.D {
	return bson.D{
		{"corporate_id", corporateID},
		{fmt.Sprintf("%s.meta.enabled", configType.String()), true},
	}
}
func makeVenueConfigQuery(corporateID, venueID string, configType entities.ConfigType) bson.D {
	return bson.D{
		{"corporate_id", corporateID},
		{"venue_id", venueID},
		{fmt.Sprintf("%s.meta.enabled", configType.String()), true},
	}
}
func makeVendorConfigQuery(corporateID, venueID, vendorID string, configType entities.ConfigType) bson.D {
	return bson.D{
		{"corporate_id", corporateID},
		{"venue_id", venueID},
		{"vendor_id", vendorID},
		{fmt.Sprintf("%s.meta.enabled", configType.String()), true},
	}
}
