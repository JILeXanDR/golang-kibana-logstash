package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"os"
	"github.com/mcuadros/go-version"
)

type Number struct {
	ID    primitive.ObjectID
	Value int `json:"value" bson:"value"`
}

type db struct {
	host    string
	name    string
	timeout time.Duration
	ctx     context.Context
	client  *mongo.Client
}

func main() {

	db := newDB("mongodb://localhost:30000,localhost:30001,localhost:30002/{db}?replicaSet=mongo-network", "testing_new_mgo_driver", 3*time.Second)
	if err := db.connect(); err != nil {
		panic(err)
	}
	defer db.client.Disconnect(db.ctx)

	if err := db.client.Database(db.name).Drop(db.ctx); err != nil {
		panic(err)
	}

	ver, err := getServerVersion(db.ctx, db.client)
	log.Println(os.Getenv("TOPOLOGY"))
	if err != nil || version.CompareSimple(ver, "4.0") < 0 || os.Getenv("TOPOLOGY") != "replica_set" {
		log.Panicf("server does not support transactions")
	}

	if err := TransactionsExamples(db.ctx, db.client); err != nil {
		panic(err)
	}
}

//func withTransaction(ctx context.Context, client *mongo.Client) {
//
//	sess, err := db.client.StartSession()
//	if err != nil {
//		panic(err)
//	}
//	defer sess.EndSession(db.ctx)
//
//	if err := db.client.UseSession(db.ctx, func(sctx mongo.SessionContext) error {
//		opts := options.Transaction().SetReadConcern(readconcern.Snapshot()).SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
//		if err := sctx.StartTransaction(opts); err != nil {
//			return err
//		}
//		res1, _ := db.insertNumber(Number{Value: 1})
//		res2, _ := db.insertNumber(Number{Value: 2})
//		res3, _ := db.insertNumber(Number{Value: 3})
//		log.Println(res1, res2, res3)
//		sctx.AbortTransaction(sctx)
//		//return errors.New("!error!")
//		//sctx.CommitTransaction(sctx)
//		return errors.New("fake error")
//	}); err != nil {
//		panic(err)
//	}
//}

func (db *db) connect() error {
	ctx, _ := context.WithTimeout(context.Background(), db.timeout)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.host))
	if err != nil {
		return err
	}

	db.ctx = ctx
	db.client = client

	return nil
}

func (db *db) insertNumber(data Number) (Number, error) {
	collection := db.client.Database(db.name).Collection("numbers")
	res, err := collection.InsertOne(db.ctx, &data)
	if err != nil {
		return Number{}, err
	}

	data.ID = res.InsertedID.(primitive.ObjectID)

	return data, nil
}

func newDB(host, name string, timeout time.Duration) *db {
	return &db{
		host:    host,
		name:    name,
		timeout: timeout,
	}
}

func getServerVersion(ctx context.Context, client *mongo.Client) (string, error) {
	serverStatus, err := client.Database("admin").RunCommand(
		ctx,
		bsonx.Doc{{"serverStatus", bsonx.Int32(1)}},
	).DecodeBytes()
	if err != nil {
		return "", err
	}

	ver, err := serverStatus.LookupErr("version")
	if err != nil {
		return "", err
	}

	return ver.StringValue(), nil
}

func TransactionsExamples(ctx context.Context, client *mongo.Client) error {
	dbName := "qwdqwdqwd"
	_, err := client.Database(dbName).Collection("employees").InsertOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	_, err = client.Database(dbName).Collection("employees").DeleteOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	_, err = client.Database(dbName).Collection("events").InsertOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	_, err = client.Database(dbName).Collection("events").DeleteOne(ctx, bson.D{{"pi", 3.14159}})
	if err != nil {
		return err
	}
	// Start Transactions Retry Example 3

	runTransactionWithRetry := func(sctx mongo.SessionContext, txnFn func(mongo.SessionContext) error) error {
		for {
			err := txnFn(sctx) // Performs transaction.
			if err == nil {
				return nil
			}

			log.Println("Transaction aborted. Caught exception during transaction.")

			// If transient error, retry the whole transaction
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				log.Println("TransientTransactionError, retrying transaction...")
				continue
			}
			return err
		}
	}

	commitWithRetry := func(sctx mongo.SessionContext) error {
		for {
			err := sctx.CommitTransaction(sctx)
			switch e := err.(type) {
			case nil:
				log.Println("Transaction committed.")
				return nil
			case mongo.CommandError:
				// Can retry commit
				if e.HasErrorLabel("UnknownTransactionCommitResult") {
					log.Println("UnknownTransactionCommitResult, retrying commit operation...")
					continue
				}
				log.Println("Error during commit...")
				return e
			default:
				log.Println("Error during commit...")
				return e
			}
		}
	}

	// Updates two collections in a transaction.
	updateEmployeeInfo := func(sctx mongo.SessionContext) error {
		employees := client.Database(dbName).Collection("employees")
		events := client.Database(dbName).Collection("events")

		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}

		_, err = employees.UpdateOne(sctx, bson.D{{"employee", 3}}, bson.D{{"$set", bson.D{{"status", "Inactive"}}}})
		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Println("caught exception during transaction, aborting.")
			return err
		}
		_, err = events.InsertOne(sctx, bson.D{{"employee", 3}, {"status", bson.D{{"new", "Inactive"}, {"old", "Active"}}}})
		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Println("caught exception during transaction, aborting.")
			return err
		}

		return commitWithRetry(sctx)
	}

	return client.UseSessionWithOptions(
		ctx, options.Session().SetDefaultReadPreference(readpref.Primary()),
		func(sctx mongo.SessionContext) error {
			return runTransactionWithRetry(sctx, updateEmployeeInfo)
		},
	)
}
