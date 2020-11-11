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

func (r MDBRepo) SetConfig(ctx context.Context, corporateID, venueID, vendorID string, config entities.ValidatedConfig) (entities.ValidatedConfig, error) {
	//configType == MAIN is invalid for SET; seems too dangerous
	filter := bson.M{
		"corporate_id": corporateID,
		"venue_id":     venueID,
		"vendor_id":    vendorID,
	}

	replaceOpts := options.FindOneAndUpdate().SetUpsert(true) //.SetReturnDocument(options.After)
	result := r.configCollection.FindOneAndUpdate(ctx, filter, bson.M{"$set": bson.M{config.GetConfigType().String(): config}}, replaceOpts)
	if err := result.Err(); err != nil && err != mongo.ErrNoDocuments {
		return nil, result.Err()
	}

	return r.GetConfig(ctx, corporateID, venueID, vendorID, config.GetConfigType())
}

func (r *MDBRepo) GetConfig(ctx context.Context, corporateID, venueID, vendorID string, configType entities.ConfigType) (entities.ValidatedConfig, error) {
	csr, err := r.configCollection.Aggregate(ctx, makeConfigPipeline(corporateID, venueID, vendorID, configType))
	if err != nil {
		return nil, err
	}
	defer csr.Close(ctx)

	return getConfigFromCursor(ctx, csr, configType)
}

//Must be updated when new config type is added
func getConfigFromCursor(ctx context.Context, cursor *mongo.Cursor, configType entities.ConfigType) (entities.ValidatedConfig, error) {
	if !cursor.Next(ctx) {
		return nil, fmt.Errorf("no config found")
	}

	var err error
	switch configType {
	case entities.CONFIG_TYPE_MAIN:
		config := entities.MainConfig{}
		err = cursor.Decode(&config)
		if err != nil {
			return nil, err
		}

		return &config, nil
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
