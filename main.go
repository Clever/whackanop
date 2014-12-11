package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Op represents a mongo operation.
type Op struct {
	ID          int    `bson:"opid"`
	Active      bool   `bson:"active"`
	Op          string `bson:"op"`
	SecsRunning int    `bson:"secs_running"`
	Namespace   string `bson:"ns"`
	Query       bson.M `bson:"query"`
}

// OpKiller kills a mongo op. Interface mostly for testing.
type OpKiller interface {
	Kill(op Op) error
}

// MongoOpKiller implements OpKiller on a real mongo database.
type MongoOpKiller struct {
	Session *mgo.Session
}

// Kill uses the $cmd.sys.killop virtual collection to kill an operation.
func (mko MongoOpKiller) Kill(op Op) error {
	return mko.Session.DB("admin").C("$cmd.sys.killop").Find(bson.M{"op": op.ID}).One(nil)
}

// OpFinder finds mongo operations. Interface mostly for testing.
type OpFinder interface {
	Find(query bson.M) ([]Op, error)
}

// MongoOpFinder implements OpFinder on a real mongo database.
type MongoOpFinder struct {
	Session *mgo.Session
}

// Find operations matching a query.
func (mfo MongoOpFinder) Find(query bson.M) ([]Op, error) {
	var result struct {
		Inprog []Op `bson:"inprog"`
	}
	err := mfo.Session.DB("admin").C("$cmd.sys.inprog").Find(query).One(&result)
	return result.Inprog, err
}

// WhackAnOp periodically finds and kills operations.
type WhackAnOp struct {
	OpFinder OpFinder
	OpKiller OpKiller
	Query    bson.M
	Tick     <-chan time.Time
	Debug    bool
	Verbose  bool
}

// Run polls for ops, killing any it finds.
func (wao WhackAnOp) Run() error {
	for _ = range wao.Tick {
		ops, err := wao.OpFinder.Find(wao.Query)
		if err != nil {
			return fmt.Errorf("whackanop: error finding ops %s", err)
		} else if wao.Verbose {
			log.Printf("found %d ops", len(ops))
		}
		for _, op := range ops {
			q, _ := json.Marshal(op.Query)
			log.Printf("opid=%d ns=%s op=%s secs_running=%d query=%s\n", op.ID,
				op.Namespace, op.Op, op.SecsRunning, q)
			if wao.Debug {
				continue
			}
			log.Printf("killing op %d", op.ID)
			if err := wao.OpKiller.Kill(op); err != nil {
				return fmt.Errorf("whackanop: error killing op %s", err)
			}
		}
	}
	return nil
}

func main() {
	flags := flag.NewFlagSet("whackanop", flag.ExitOnError)
	mongourl := flags.String("mongourl", "localhost", "mongo url to connect to")
	interval := flags.Int("interval", 1, "how often, in seconds, to poll mongo for operations")
	querystr := flags.String("query", `{"op": "query", "secs_running": {"$gt": 60}}`, "query sent to db.currentOp()")
	debug := flags.Bool("debug", true, "in debug mode, operations that match the query are logged instead of killed")
	version := flags.Bool("version", false, "print the version and exit")
	verbose := flags.Bool("verbose", false, "more verbose logging")
	flags.Parse(os.Args[1:])

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	var query bson.M
	if err := json.Unmarshal([]byte(*querystr), &query); err != nil {
		log.Fatal(err)
	}

	session, err := mgo.Dial(*mongourl)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	log.Printf("mongourl=%s interval=%d debug=%t query=%#v", *mongourl, *interval, *debug, query)

	wao := WhackAnOp{
		OpFinder: MongoOpFinder{session},
		OpKiller: MongoOpKiller{session},
		Query:    query,
		Tick:     time.Tick(time.Duration(*interval) * time.Second),
		Debug:    *debug,
		Verbose:  *verbose,
	}
	if err := wao.Run(); err != nil {
		log.Fatal(err)
	}
}
