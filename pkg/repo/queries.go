package repo

import (
	"fmt"

	"github.com/mcquackers/config-demo/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func makeConfigPipeline(corporateID, venueID, vendorID string, configType entities.ConfigType) mongo.Pipeline {
	switch configType {
	case entities.CONFIG_TYPE_UNSPECIFIED: // || !configType.Valid()
		return mongo.Pipeline{}
	case entities.CONFIG_TYPE_MAIN:
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
	return append(makeMainConfigPipeline(corporateID, venueID, vendorID), bson.D{
		{
			Key: "$replaceRoot",
			Value: bson.D{{
				Key:   "newRoot",
				Value: fmt.Sprintf("$%s", configType.String()),
			}},
		},
	})
}
