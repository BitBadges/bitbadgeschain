package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	cmtstore "github.com/cometbft/cometbft/proto/tendermint/store"
	proto "github.com/cosmos/gogoproto/proto"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	// Define command line flags
	sourcePtr := flag.String("source", "", "Path to the source snapshot directory (required)")
	targetPtr := flag.String("target", "", "Path to the target node directory (required)")

	// Parse flags
	flag.Parse()

	// Validate required flags
	if *sourcePtr == "" || *targetPtr == "" {
		fmt.Println("Error: Both source and target paths are required")
		fmt.Println("Usage: go run migrate.go -source /path/to/source -target /path/to/target")
		os.Exit(1)
	}

	// Use the provided paths
	sourceDBPath := *sourcePtr
	targetDBPath := *targetPtr

	txIndexSnapshotDB, err := leveldb.OpenFile(sourceDBPath+"/tx_index.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	blockstoreSnapshotDB, err := leveldb.OpenFile(sourceDBPath+"/blockstore.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	targetTxIndexDB, err := leveldb.OpenFile(targetDBPath+"/tx_index.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	targetBlockstoreDB, err := leveldb.OpenFile(targetDBPath+"/blockstore.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer txIndexSnapshotDB.Close()
	defer blockstoreSnapshotDB.Close()
	defer targetTxIndexDB.Close()
	defer targetBlockstoreDB.Close()

	//Iterate over the tx_index.db and copy the keys and values to the targetDB
	//If we already have the key in the targetDB, skip it and do not overwrite
	iter := txIndexSnapshotDB.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if _, err := targetTxIndexDB.Get(key, nil); err == nil {
			continue
		}

		targetTxIndexDB.Put(key, value, nil)
	}

	iter = blockstoreSnapshotDB.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if _, err := targetBlockstoreDB.Get(key, nil); err == nil {
			allowedPrefixes := []string{"BH", "SC", "C", "P", "H"}
			parts := strings.Split(string(key), ":")
			prefix := parts[0]
			if slices.Contains(allowedPrefixes, prefix) {
				//These are fine if we already have them
				continue
			} else {
				//We need to update the starting height of what we have (key == "blockStore")
				bytes, err := targetBlockstoreDB.Get([]byte(key), nil)
				if err != nil {
					panic(err)
				}

				//Set the base to 1 (since we have height 1 on now from the snapshot)
				var bsj cmtstore.BlockStoreState
				if err := proto.Unmarshal(bytes, &bsj); err != nil {
					panic(fmt.Sprintf("Could not unmarshal bytes: %X", bytes))
				}

				bsj.Base = 1

				bytes, err = proto.Marshal(&bsj)
				if err != nil {
					panic(err)
				}

				targetBlockstoreDB.Put([]byte(key), bytes, nil)

				continue
			}
		}

		targetBlockstoreDB.Put(key, value, nil)
	}
}
