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
	demoConfig := &entities.CloudCartConfig{
		ConfigMeta: entities.ConfigMeta{
			ChangedAt:  time.Now(),
			ChangedBy:  "Sherlock Holmes",
		},
		EnableCalculateReductionsAndTaxes: true,
		EnableValidateCartSums:            true,
		EnableValidatePrices:              false,
	}

	fmt.Println("Upsert new config")
	returnConf, err := repo.SetConfig(context.Background(), demoCorpID, demoVenueID, demoVendorID, demoConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("%+v\n", returnConf)
	fmt.Println("=================================")

	fmt.Println("Retrieve MAIN config")
	fullConf, err := repo.GetConfig(context.Background(), demoCorpID, demoVenueID, demoVendorID, entities.CONFIG_TYPE_MAIN)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("=================================")
	fmt.Printf("%+v\n", fullConf)
	fmt.Println("=================================")

	demoConfig2 := &entities.OtherConfig{
		ConfigMeta: entities.ConfigMeta{
			ChangedBy: "Watson",
			ChangedAt: time.Now(),
		},
		ADifferentValue: "hello",
		AFloat: 192.43,
	}

	fmt.Println("Set new config value on existing main config")
	returnConf, err = repo.SetConfig(context.Background(), demoCorpID, demoVenueID, demoVendorID, demoConfig2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("=================================")
	fmt.Printf("%+v\n", returnConf)
	fmt.Println("=================================")


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
