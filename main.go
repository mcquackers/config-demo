package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mcquackers/config-demo/pkg/entities"
	"github.com/mcquackers/config-demo/pkg/repo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	client := mustSetUpClient()
	repo := repo.NewMDBRepo(client)
	fmt.Println("client set up")

	demoCorpID := "1"
	demoVenueID := "2"
	demoVendorID := "3"
	demoConfigVendor := &entities.CloudCartConfig{
		ConfigMeta: entities.ConfigMeta{
			Enabled:   false,
			ChangedAt: time.Now(),
			ChangedBy: "VENDOR",
		},
		EnableCalculateReductionsAndTaxes: true,
		EnableValidateCartSums:            true,
		EnableValidatePrices:              false,
	}

	demoConfigVenue := &entities.CloudCartConfig{
		ConfigMeta: entities.ConfigMeta{
			Enabled:   true,
			ChangedAt: time.Now(),
			ChangedBy: "VENUE",
		},
		EnableCalculateReductionsAndTaxes: true,
		EnableValidateCartSums:            true,
		EnableValidatePrices:              false,
	}

	demoConfigCorporate := &entities.CloudCartConfig{
		ConfigMeta: entities.ConfigMeta{
			Enabled:   true,
			ChangedAt: time.Now(),
			ChangedBy: "CORPORATE",
		},
		EnableCalculateReductionsAndTaxes: true,
		EnableValidateCartSums:            true,
		EnableValidatePrices:              false,
	}

	fmt.Println("Retrieve unset configuration")
	returnConf, err := repo.GetSpecificConfig(context.Background(), entities.CONFIG_LEVEL_CORPORATE, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_DEMO_CONFIG)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", returnConf)
	fmt.Println("=================================")


	fmt.Println("Upsert new config vendor level - inactive")
	returnConf, err = repo.SetConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, demoConfigVendor)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", returnConf)
	fmt.Println("=================================")

	fmt.Println("Retrieve MAIN config")
	fullConf, err := repo.GetSpecificConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_FULL)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("=================================")
	fmt.Printf("%+v\n", fullConf)
	fmt.Println("=================================")

	demoConfig2 := &entities.OtherConfig{
		ConfigMeta: entities.ConfigMeta{
			Enabled:   true,
			ChangedBy: "Watson",
			ChangedAt: time.Now(),
		},
		ADifferentValue: "hello",
		AFloat:          192.43,
	}

	fmt.Println("Set new config value on existing main config")
	returnConf, err = repo.SetConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, demoConfig2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("=================================")
	fmt.Printf("%+v\n", returnConf)
	fmt.Println("=================================")

	fmt.Println("Retrieve MAIN config")
	fullConf, err = repo.GetSpecificConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_FULL)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("=================================")
	fmt.Printf("%+v\n", fullConf)
	fmt.Println("=================================")

	fmt.Println("Retrieve OtherConfig")
	otherConf, err := repo.GetSpecificConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_OTHER_EXAMPLE)
	fmt.Println("=================================")
	fmt.Printf("%+v\n", otherConf)
	fmt.Println("=================================")

	fmt.Println("retrieve active configuration starting with vendor - no active expected")
	ccConf, err := repo.GetActiveConfig(context.Background(), entities.CONFIG_LEVEL_VENUE, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_DEMO_CONFIG)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", ccConf)
	fmt.Println("=================================")

	fmt.Println("Set Active Venue level config")
	_,err = repo.SetConfig(context.Background(), entities.CONFIG_LEVEL_VENUE, demoCorpID, demoVenueID, demoVendorID, demoConfigVenue)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Attempt to retrieve active vendor level demo config; expect venue level config")
	ccConf, err = repo.GetActiveConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_DEMO_CONFIG)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", ccConf)
	fmt.Println("=================================")

	fmt.Println("Set corporate level demo config - active")
	_,_ = repo.SetConfig(context.Background(), entities.CONFIG_LEVEL_CORPORATE, demoCorpID, demoVenueID, demoVendorID, demoConfigCorporate)


	fmt.Println("Attempt to retrieve active vendor level demo config; expect venue level config")
	dc, err := repo.GetActiveConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_DEMO_CONFIG)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", dc)
	fmt.Println("=================================")

	fmt.Println("Disable venue level demo config")
	demoConfigVenue.ConfigMeta.Enabled = false
	ccConf, _ = repo.SetConfig(context.Background(), entities.CONFIG_LEVEL_VENUE, demoCorpID, demoVenueID, demoVendorID, demoConfigVenue)
	fmt.Println("Venue level demo config")
	fmt.Println("=================================")
	fmt.Printf("%+v\n", ccConf)
	fmt.Println("=================================")

	fmt.Println("Attempt to retrieve active vendor level demo config; expect corporate level config")
	dc, err = repo.GetActiveConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_DEMO_CONFIG)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", dc)
	fmt.Println("=================================")
	uc, err := repo.GetActiveConfig(context.Background(), entities.CONFIG_LEVEL_VENDOR, "5", "6", "7", entities.CONFIG_TYPE_OTHER_EXAMPLE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", uc)


	client.Disconnect(context.Background())
}

func mustSetUpClient() *mongo.Client {
	mdbConnectionOpts := options.Client().
		SetConnectTimeout(5 * time.Second).
		SetHosts([]string{"localhost:27017"}).
		SetReplicaSet("testRepl")

	mdbClient, err := mongo.NewClient(mdbConnectionOpts)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(255)
	}

	err = mdbClient.Connect(context.Background())
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(255)
	}

	err = mdbClient.Ping(context.Background(), readpref.PrimaryPreferred())
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(255)
	}

	return mdbClient
}
