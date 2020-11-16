package repo

import (
	"context"
	"fmt"

	"github.com/mcquackers/config-demo/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MDBRepo struct {
	client           *mongo.Client
	configCollection *mongo.Collection
}

func NewMDBRepo(client *mongo.Client) MDBRepo {
	return MDBRepo{
		client:           client,
		configCollection: client.Database("config-demo").Collection("configs"),
	}
}

func (r MDBRepo) SetConfig(ctx context.Context, configLevel entities.ConfigLevel, corporateID, venueID, vendorID string, config entities.ValidatedConfig) (entities.ValidatedConfig, error) {
	//configType == MAIN is invalid for SET; seems too dangerous
	filter, err := makeUpsertConfigFilter(configLevel, corporateID, venueID, vendorID)
	if err != nil {
		return nil, err
	}
	fmt.Println(filter)

	replaceOpts := options.FindOneAndUpdate().SetUpsert(true) //.SetReturnDocument(options.After)
	result := r.configCollection.FindOneAndUpdate(ctx, filter, bson.M{"$set": bson.M{config.GetConfigType().String(): config}}, replaceOpts)
	if err := result.Err(); err != nil && err != mongo.ErrNoDocuments {
		return nil, result.Err()
	}

	return r.GetSpecificConfig(ctx, configLevel, corporateID, venueID, vendorID, config.GetConfigType())
}

func (r *MDBRepo) GetSpecificConfig(ctx context.Context, configLevel entities.ConfigLevel, corporateID, venueID, vendorID string, configType entities.ConfigType) (entities.ValidatedConfig, error) {
	csr, err := r.configCollection.Aggregate(ctx, makeConfigPipeline(corporateID, venueID, vendorID, configType))
	if err != nil {
		return nil, err
	}
	defer csr.Close(ctx)

	return getConfigFromCursor(ctx, csr, configLevel, configType)
}

func (r *MDBRepo) GetActiveConfig(ctx context.Context, configLevel entities.ConfigLevel, corporateID, venueID, vendorID string, configType entities.ConfigType) (entities.ValidatedConfig, error) {
	csr, err := r.configCollection.Aggregate(ctx, makeGetActiveConfigPipeline(configLevel, corporateID, venueID, vendorID, configType))
	if err != nil {
		return nil, err
	}

	defer csr.Close(ctx)

	return getConfigFromCursor(ctx, csr, configLevel, configType)
}

//Must be updated when new config type is added
func getConfigFromCursor(ctx context.Context, cursor *mongo.Cursor, configLevel entities.ConfigLevel, configType entities.ConfigType) (entities.ValidatedConfig, error) {
	if !cursor.Next(ctx) {
		return emptyConfigForType(configLevel, configType), nil
	}

	var err error
	switch configType {
	case entities.CONFIG_TYPE_FULL:
		switch configLevel{
		case entities.CONFIG_LEVEL_CORPORATE:
			config := entities.CorporateConfig{}
			err = cursor.Decode(&config)
			if err != nil {
				return nil, err
			}
			return &config, nil
		case entities.CONFIG_LEVEL_VENUE:
			config := entities.VenueConfig{}
			err = cursor.Decode(&config)
			if err != nil {
				return nil, err
			}
			return &config, nil
		case entities.CONFIG_LEVEL_VENDOR:
			config := entities.VendorConfig{}
			err = cursor.Decode(&config)
			if err != nil {
				return nil, err
			}
			return &config, nil
		default: return nil, fmt.Errorf("unspecified config level: %d", configLevel)
		}
	case entities.CONFIG_TYPE_DEMO_CONFIG:
		config := entities.CloudCartConfig{}
		err = cursor.Decode(&config)
		if err != nil {
			return nil, err
		}
		return &config, nil

	case entities.CONFIG_TYPE_OTHER_EXAMPLE:
		config := entities.OtherConfig{}
		err = cursor.Decode(&config)
		if err != nil {
			return nil, err
		}
		return &config, nil
	}

	return nil, fmt.Errorf("unsupported config type")
}

func emptyConfigForType(configLevel entities.ConfigLevel, configType entities.ConfigType) entities.ValidatedConfig {
	switch configType {
	case entities.CONFIG_TYPE_FULL:
		switch configLevel{
		case entities.CONFIG_LEVEL_CORPORATE:
			config := entities.CorporateConfig{}
			return &config
		case entities.CONFIG_LEVEL_VENUE:
			config := entities.VenueConfig{}
			return &config
		case entities.CONFIG_LEVEL_VENDOR:
			config := entities.VendorConfig{}
			return &config
		default: return nil
		}
	case entities.CONFIG_TYPE_DEMO_CONFIG:
		config := entities.CloudCartConfig{}
		return &config

	case entities.CONFIG_TYPE_OTHER_EXAMPLE:
		config := entities.OtherConfig{}
		return &config
	}
	return nil
}
