//go:generate go-bindata -prefix ../../migrations/ -pkg migrations -o ../../internal/migrations/migrations_gen.go ../../migrations/

package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/Seaman-hub/numberserver/api/ns"
	"github.com/Seaman-hub/numberserver/internal/api"
	"github.com/Seaman-hub/numberserver/internal/common"
	"github.com/Seaman-hub/numberserver/internal/storage"
)

func init() {
	grpclog.SetLogger(log.StandardLogger())
}

var version string // set by the compiler

func run(c *cli.Context) error {
	tasks := []func(*cli.Context) error{
		printStartMessage,
		// setEtcdConnection,
		startAPIServer,
	}

	for _, t := range tasks {
		if err := t(c); err != nil {
			log.Fatal(err)
		}
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping numberserver")
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}

func printStartMessage(c *cli.Context) error {
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "not implemented",
	}).Info("starting Number Server")
	return nil
}

func setEtcdConnection(c *cli.Context) error {
	log.Info("connecting to etcd")
	db, err := storage.OpenDatabase(c)
	if err != nil {
		log.Fatalln("Failed to open storage:", err)
	}
	common.DB = db
	return nil
}

func startAPIServer(c *cli.Context) error {
	log.WithFields(log.Fields{
		"bind": c.String("bind"),
	}).Info("starting Number api server")

	var opts []grpc.ServerOption
	gs := grpc.NewServer(opts...)
	nsAPI := api.NewNumberServerAPI(uint(c.Int("id-initial")))
	ns.RegisterNumberServerServer(gs, nsAPI)

	ln, err := net.Listen("tcp", c.String("bind"))
	if err != nil {
		return errors.Wrap(err, "start api listener error")
	}
	go gs.Serve(ln)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "numberserver"
	app.Usage = "number-server for SDN networks"
	app.Version = version
	app.Copyright = "See http://github.com/Seaman-hub/numberserver for copyright information"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "bind",
			Usage:  "ip:port to bind the api server",
			Value:  "0.0.0.0:8000",
			EnvVar: "BIND",
		},
		cli.StringFlag{
			Name:   "storage-server",
			Value:  "localhost:2379",
			Usage:  "connect to storage server",
			EnvVar: "STORAGE_SERVER",
		},
		cli.IntFlag{
			Name:   "id-initial",
			Value:  1,
			Usage:  "initial id value to initialize ID pool",
			EnvVar: "ID_INITIAL",
		},
	}

	app.Run(os.Args)
}
