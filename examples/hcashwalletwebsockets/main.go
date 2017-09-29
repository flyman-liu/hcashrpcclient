// Copyright (c) 2014-2015 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/HcashOrg/hcashrpcclient"
	"github.com/HcashOrg/hcashutil"
)

func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the hcashrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := hcashrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance hcashutil.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}

	// Connect to local hcashwallet RPC server using websockets.
	certHomeDir := hcashutil.AppDataDir("hcashwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &hcashrpcclient.ConnConfig{
		Host:         "localhost:9110",
		Endpoint:     "ws",
		User:         "yourrpcuser",
		Pass:         "yourrpcpass",
		Certificates: certs,
	}
	client, err := hcashrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	unspent, err := client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Num unspent outputs (utxos): %d", len(unspent))
	if len(unspent) > 0 {
		log.Printf("First utxo:\n%v", spew.Sdump(unspent[0]))
	}

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*10, func() {
		log.Println("Client shutting down...")
		client.Shutdown()
		log.Println("Client shutdown complete.")
	})

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}