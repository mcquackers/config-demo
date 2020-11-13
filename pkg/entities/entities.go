package entities

import (
	"time"
)

type ConfigType int

//Must be updated when adding new config, part of full ENUM pattern
const (
	CONFIG_TYPE_UNSPECIFIED = iota
	CONFIG_TYPE_MAIN
	CONFIG_TYPE_DEMO_CONFIG
	CONFIG_TYPE_OTHER_EXAMPLE
)

type ValidatedConfig interface {
	Validate() error
	GetConfigType() ConfigType
	//UnmarshalBSONValue(bsontype.Type, []byte) error
}

//Must be updated when adding new config; irrelevant once full ENUM pattern is in place
func (ct ConfigType) String() string {
	switch ct {
	case CONFIG_TYPE_DEMO_CONFIG:
		return "cloud_cart"
	case CONFIG_TYPE_OTHER_EXAMPLE:
		return "other_example"
	default:
		return ""
	}
}

type MainConfig struct {
	CorporateID string          `bson:"corporate_id"`
	VenueID     string          `bson:"venue_id,omitempty"`
	VendorID    string          `bson:"vendor_id,omitempty"`
	CloudCart   CloudCartConfig `bson:"cloud_cart"`
	OtherConfig OtherConfig     `bson:"other_example"`
}

func (c *MainConfig) Validate() error {
	return nil
}

func (c *MainConfig) GetConfigType() ConfigType {
	return CONFIG_TYPE_MAIN
}

type ConfigMeta struct {
	Enabled   bool      `bson:"enabled"`
	ChangedBy string    `bson:"changed_by"`
	ChangedAt time.Time `bson:"changed_at"`
}

type CloudCartConfig struct {
	ConfigMeta                        `bson:"meta"`
	EnableCalculateReductionsAndTaxes bool `bson:"enable_calculate_reductions_and_taxes"`
	EnableValidatePrices              bool `bson:"enable_validate_prices"`
	EnableValidateCartSums            bool `bson:"enable_validate_cart_sums"`
}

func (c *CloudCartConfig) Validate() error {
	return nil
}

func (c *CloudCartConfig) GetConfigType() ConfigType {
	return CONFIG_TYPE_DEMO_CONFIG
}

type OtherConfig struct {
	ConfigMeta      `bson:"meta"`
	ADifferentValue string  `bson:"a_different_value"`
	AFloat          float32 `bson:"a_float"`
}

func (c *OtherConfig) Validate() error {
	return nil
}

func (c *OtherConfig) GetConfigType() ConfigType {
	return CONFIG_TYPE_OTHER_EXAMPLE
}

//func (c *CloudCartConfig) UnmarshalBSONValue(_ bsontype.Type, raw []byte) error {
//return bson.Unmarshal(raw, c) //Caused recursive loop; under the hood, Unmarshal looks for this interface method
//}
