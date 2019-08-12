package hdm

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)
import "github.com/pilebones/go-udev/netlink"
import "github.com/kr/pretty"

func toto() {
	matcher := netlink.RuleDefinitions{
		Rules: []netlink.RuleDefinition{
			{
				//Action: "",
				Env: map[string]string{
					"SUBSYSTEM" : "block",
					"DEVTYPE" : "disk",
					//  ID_BUS=ata
					//	ID_TYPE=disk
				},
			},
		},
	}

	conn := new(netlink.UEventConn)
	if err := conn.Connect(netlink.UdevEvent); err != nil {
		log.Fatalln("Unable to connect to Netlink Kobject UEvent socket")
	}
	defer conn.Close()


	queue := make(chan netlink.UEvent)
	errors := make(chan error)
	quit := conn.Monitor(queue, errors, &matcher)

	// Signal handler to quit properly monitor mode
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-signals
		log.Println("Exiting monitor mode...")
		quit <- struct{}{}
		os.Exit(0)
	}()

	// Handling message from queue
	for {
		select {
		case uevent := <-queue:
			//uevent.

			log.Printf("Handle %s\n", pretty.Sprint(uevent))
		case err := <-errors:
			log.Printf("ERROR: %v", err)
		}
	}



	// watch disk events -> add/remove disks
	// watch files events -> run sync
	// scan for backup -> timing sync
	//

}