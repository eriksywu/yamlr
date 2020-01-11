package main

import (
	"flag"
	"fmt"
	"gitlab.com/erikwu09/yamlr/app"
	"gitlab.com/erikwu09/yamlr/memoryrepo"
	"gitlab.com/erikwu09/yamlr/server"
	"gitlab.com/erikwu09/yamlr/validation"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting service")
	var port = flag.Int("port", 8082, "port")
	logger.Println(fmt.Sprintf("listening on port %d", *port))
	repo, err := memoryrepo.GetMemoryRepository()
	if err != nil {
		logger.Fatal(err)
	}
	mgr, _ := app.BuildMetadataManager(repo, validation.SimpleValidator{}, logger)
	service := server.BuildYamlApp(*port, &mgr, logger)
	service.Run()
}
