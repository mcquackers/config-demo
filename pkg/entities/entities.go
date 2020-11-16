package entities

import (
	"time"
)

type ConfigType int
type ConfigLevel int

//Must be updated when adding new config, part of full ENUM pattern
const (
	CONFIG_TYPE_UNSPECIFIED = iota
	CONFIG_TYPE_FULL
	CONFIG_TYPE_DEMO_CONFIG
	CONFIG_TYPE_OTHER_EXAMPLE
)

const (
	CONFIG_LEVEL_UNSPECIFIED = iota
	CONFIG_LEVEL_CORPORATE
	CONFIG_LEVEL_VENUE
	CONFIG_LEVEL_VENDOR
)

//This association isn't executed as part of the demo,  but it would live in the Service layer
//If a request attempts to set a configuration on a level that isn't associated to that type, reject it as a bad request
//The idea is that certain configurations _cannot_ be overridden at certain levels.  Some configs will be corporate only,
//venue only, vendor only, etc.
//This shouldn't affect the configuration retrieval logic/code at all.  If a request for a venue-associated configuration is made
//with CONFIG_LEVEL_VENDOR, it should seek the first active configuration above it.

//Make no mistake, this _is_ an array, not a map.  Was very confused that this compiled, but I'm grateful for the explicitness
var CONFIG_LEVEL_TO_TYPE_ASSOCIATIONS = [4][]ConfigType{
	CONFIG_LEVEL_UNSPECIFIED: {},
	CONFIG_LEVEL_CORPORATE: {
		CONFIG_TYPE_DEMO_CONFIG,
		CONFIG_TYPE_OTHER_EXAMPLE, //New Configuration types should always be appended; this means the []ConfigType array
		//will always be sorted, meaning we can execute a sort.Search (https://golang.org/pkg/sort/#Search)
		//to get the expected index of the config type.  We then directly access it and compare it to the
		//Requested configuration type.  If equal, accept, if not equal, reject as bad request

		//  configLevel == CONFIG_LEVEL_CORPORATE ; configType == CONFIG_TYPE_OTHER_EXAMPLE
		//ex:  if index := sort.Search(len(CONFIG_LEVEL_TO_TYPE_ASSOCIATIONS[configLevel]), func(i int) bool {
		  			//return i >= configType
		//}); len(CONFIG_LEVEL_TO_TYPE_ASSOCIATIONS[configLevel]) > index && CONFIG_LEVEL_TO_TYPE_ASSOCIATIONS[configLevel][index] == configType
		//A bit verbose, easy to refactor
		//This is much more performant than a map
	},
	CONFIG_LEVEL_VENUE: {
		CONFIG_TYPE_DEMO_CONFIG,
	},
	CONFIG_LEVEL_VENDOR: {
		CONFIG_TYPE_OTHER_EXAMPLE,
	},
}

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

type CorporateConfig struct {
	CorporateID string          `bson:"corporate_id"`
	CloudCart   CloudCartConfig `bson:"cloud_cart"`
	OtherConfig OtherConfig     `bson:"other_example"`
}
type VenueConfig struct {
	CorporateID string          `bson:"corporate_id"`
	VenueID     string          `bson:"venue_id,omitempty"`
	CloudCart   CloudCartConfig `bson:"cloud_cart"`
	OtherConfig OtherConfig     `bson:"other_example"`
}
type VendorConfig struct {
	CorporateID string          `bson:"corporate_id"`
	VenueID     string          `bson:"venue_id,omitempty"`
	VendorID    string          `bson:"vendor_id,omitempty"`
	CloudCart   CloudCartConfig `bson:"cloud_cart"`
	OtherConfig OtherConfig     `bson:"other_example"`
}

func (c *CorporateConfig) Validate() error {
	return nil
}

func (c *CorporateConfig) GetConfigType() ConfigType {
	return CONFIG_TYPE_FULL
}
func (c *VenueConfig) Validate() error {
	return nil
}

func (c *VenueConfig) GetConfigType() ConfigType {
	return CONFIG_TYPE_FULL
}
func (c *VendorConfig) Validate() error {
	return nil
}

func (c *VendorConfig) GetConfigType() ConfigType {
	return CONFIG_TYPE_FULL
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
