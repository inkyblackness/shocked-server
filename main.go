package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/docopt/docopt-go"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"

	core "github.com/inkyblackness/shocked-core"
	"github.com/inkyblackness/shocked-core/release"
	"github.com/inkyblackness/shocked-server/app"
)

func usage() string {
	return app.Title + `

Usage:
	shocked-server --source=<srcdir> --projects=<prjdir> [--swagger=<swdir>]
	shocked-server -h | --help
	shocked-server --version

Options:
	-h --help            Show this screen.
	--version            Show version.
	--source=<srcdir>    A path pointing to the root of a System Shock source directory
	--projects=<prjdir>  A path pointing to a directory containing the projects
	--swagger=<swdir>    An optional path pointing to the Swagger UI resources
`
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, app.Title, false)
	port := 8080
	log.Printf("Arguments: %v", arguments)

	source, srcErr := release.ReleaseFromDir(arguments["--source"].(string))
	if srcErr != nil {
		log.Fatalf("Source is not available: %v", srcErr)
		return
	}
	projects, prjErr := release.NewContainerFromDir(arguments["--projects"].(string))
	if prjErr != nil {
		log.Fatalf("Projects dir is not available: %v", prjErr)
		return
	}

	workspace := core.NewWorkspace(source, projects)
	wsContainer := restful.NewContainer()

	app.NewWorkspaceResource(wsContainer, workspace)

	swDir := arguments["--swagger"]
	if swDir != nil {
		config := swagger.Config{
			WebServices:     wsContainer.RegisteredWebServices(), // you control what services are visible
			WebServicesUrl:  fmt.Sprintf("http://localhost:%d", port),
			ApiPath:         "/apidocs.json",
			ApiVersion:      "0.1",
			SwaggerPath:     "/apidocs/",
			SwaggerFilePath: swDir.(string)}
		swagger.RegisterSwaggerService(config, wsContainer)
	}

	log.Printf("start listening on localhost:%d", port)
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
