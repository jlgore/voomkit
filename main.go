package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/netip"
	"os"
  "math/rand"
  "time"
	"os/signal"
	"syscall"
	"github.com/brutella/hap"
	"github.com/brutella/hap/characteristic"
	"github.com/jurjevic/pixoo64"
	"github.com/xxandev/homekit"
)

type Config struct{ homekit.AccessoryConfig }

type DConfig struct{
  DivoomAddr string
}

var (
	debug  bool
	config Config
  dconfig DConfig
)

func init() {
	log.SetOutput(os.Stdout) // log.SetOutput(ioutil.Discard)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)

	flag.BoolVar(&debug, "d", false, "hap debug log activate")
	flag.StringVar(&config.Name, "n", "pixoo64", "homekit accessory name")
	flag.StringVar(&config.SN, "sn", "pixoo64", "homekit accessory serial number")
	flag.StringVar(&config.Host, "h", "", "homekit host, example: 192.168.1.xxx")
	flag.UintVar(&config.Port, "p", 10726, "homekit port, example: 10101, 10102...")
	flag.StringVar(&config.Pin, "pin", "19378246", "homekit pin, example: 82143697, 13974682")
  flag.StringVar(&dconfig.DivoomAddr, "DivoomAddr", "192.168.88.53", "pixoo64 IP address: 192.168.XX.XX")
	flag.Parse()

	homekit.OnLog(debug)
}

func main() {

  addr, err := netip.ParseAddr(dconfig.DivoomAddr)
  if err != nil {
    log.Fatal("Invalid IP Address")
  }

  device, err := pixoo64.NewDevice(addr)
  if err != nil {
    log.Fatal("Failed to add Pixoo", err)
  }


	acc := homekit.NewAccessoryTelevision(config.GetInfo("pixoo64"))
	_ = acc.AddInputSource(1, "Random Visualizer", characteristic.InputSourceTypeHdmi)
	_ = acc.AddInputSource(2, "Favorite Cloud Channel", characteristic.InputSourceTypeHdmi)
	_ = acc.AddInputSource(3, "Recommended Cloud Channel", characteristic.InputSourceTypeHdmi)
	_ = acc.AddInputSource(4, "Subscribed Artist Cloud Channel", characteristic.InputSourceTypeHdmi)
	_ = acc.AddInputSource(5, "Album Cloud Channel", characteristic.InputSourceTypeOther)
	log.SetPrefix(fmt.Sprintf("[%T] <%v> ", acc, acc.GetSN()))
	storage := hap.NewFsStore(fmt.Sprintf("./%s", acc.GetSN()))
	server, err := hap.NewServer(storage, acc.GetAccessory())
	if err != nil {
		log.Fatalf("error create hap server: %v\n", err)
	}
	log.Printf("hap server create successful.\n")

	acc.Television.Active.OnValueRemoteUpdate(func(v int) {
    if v == 0 {
     device.Display().Switch(0)
   } else if v == 1 {
     device.Display().Switch(1)
   } else {
     log.Fatal("Homekit sent invalid interger to Device")
   }
    log.Printf("remote update active: %[1]T - %[1]v\n", v)
	})

	acc.Television.ActiveIdentifier.OnValueRemoteUpdate(func(v int) {
    switch {
    case v == 1:
      rand.Seed(time.Now().UnixNano())
      max := 10
      device.Channel().ActivateVisualizer(rand.Intn(max -1 + 1) + 1)
    case v == 2:
      device.Channel().ActivateCloudIndex(1)
    case v == 3:
      device.Channel().ActivateCloudIndex(0)
    case v == 4:
      device.Display().Reset()
      device.Display().ResetText()
      device.Channel().ActivateCloudIndex(2)
    case v == 5:
      device.Display().Reset()
      device.Display().ResetText()
      device.Channel().ActivateCloudIndex(3)
    } 

		log.Printf("remote update active identifier: %[1]T - %[1]v\n", v)
  })
	acc.Television.ConfiguredName.OnValueRemoteUpdate(func(v string) {
		log.Printf("remote update configured name: %[1]T - %[1]v\n", v)
	})
	acc.Television.SleepDiscoveryMode.OnValueRemoteUpdate(func(v int) {
		log.Printf("remote update sleep discovery mode: %[1]T - %[1]v\n", v)
	})
	acc.Television.Brightness.OnValueRemoteUpdate(func(v int) {
		log.Printf("remote update brightness: %[1]T - %[1]v\n", v)
	})
	acc.Television.ClosedCaptions.OnValueRemoteUpdate(func(v int) {
		log.Printf("remote update closed captions: %[1]T - %[1]v\n", v)
	})
	acc.Television.DisplayOrder.OnValueRemoteUpdate(func(v []byte) {
		log.Printf("remote update display order: %[1]T - %[1]v\n", v)
	})
	acc.Television.CurrentMediaState.OnValueRemoteUpdate(func(v int) {
		log.Printf("remote update current media state: %[1]T - %[1]v\n", v)
	})
	acc.Television.TargetMediaState.OnValueRemoteUpdate(func(v int) {
		log.Printf("remote update target media state: %[1]T - %[1]v\n", v)
	})
	acc.Television.PowerModeSelection.OnValueRemoteUpdate(func(v int) {
    log.Printf("remote update power mode selection: %[1]T - %[1]v\n", v)
	})
	acc.Television.PictureMode.OnValueRemoteUpdate(func(v int) {
		log.Printf("remote update picture mode: %[1]T - %[1]v\n", v)
	})

	acc.Television.RemoteKey.OnValueRemoteUpdate(func(v int) {
		switch v {
		case characteristic.RemoteKeyRewind:
			log.Printf("remote update key rewind: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyFastForward:
			log.Printf("remote update key fast forward: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyExit:
			log.Printf("remote update key exit: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyPlayPause:
			log.Printf("remote update key play/pause: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyInfo:
			log.Printf("remote update key info: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyNextTrack:
			log.Printf("remote update key next: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyPrevTrack:
			log.Printf("remote update key prev: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyArrowUp:
			log.Printf("remote update key up: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyArrowDown:
			log.Printf("remote update key down: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyArrowLeft:
			log.Printf("remote update key left: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyArrowRight:
			log.Printf("remote update key right: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeySelect:
			log.Printf("remote update key select: %[1]T - %[1]v\n", v)
		case characteristic.RemoteKeyBack:
			log.Printf("remote update key back: %[1]T - %[1]v\n", v)
		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sig
		log.Println("program stop.")
		signal.Stop(sig)
		cancel()
	}()
	homekit.SetServer(server, config.GetAddress(), config.GetPin())
	log.Printf("hap server starting set, address: %v, pin: %v.\n", server.Addr, server.Pin)
	log.Fatalf("hap server: %v\n", server.ListenAndServe(ctx))
}
