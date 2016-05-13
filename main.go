package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

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
	shocked-server --source=<srcdir> --projects=<prjdir> [--swagger=<swdir>] [--client=<clientdir>] [--address=<addr>]
	shocked-server -h | --help
	shocked-server --version

Options:
	-h --help             Show this screen.
	--version             Show version.
	--source=<srcdir>     A path pointing to the root of a System Shock source directory
	--projects=<prjdir>   A path pointing to a directory containing the projects
	--swagger=<swdir>     An optional path pointing to the Swagger UI resources
	--client=<clientdir>  An optional path pointing to the client directory
	--address=<addr>      The ip:port combination to listen on. Default: "localhost:8080".
`
}

func serveClient(container *restful.Container, localPath string) {
	rootDir := localPath

	handleRequest := func(req *restful.Request, resp *restful.Response) {
		actual := path.Join(rootDir, req.PathParameter("subpath"))
		http.ServeFile(resp.ResponseWriter, req.Request, actual)
	}

	ws := new(restful.WebService)
	ws.Route(ws.GET("/client/{subpath:*}").To(handleRequest))
	container.Add(ws)
	log.Printf("Client added from ", localPath)
}

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, app.Title, false)
	addressArg := arguments["--address"]
	address := "localhost:8080"
	log.Printf("Arguments: %v", arguments)

	if addressArg != nil {
		address = addressArg.(string)
	}

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

	clientDir := arguments["--client"]
	if clientDir != nil {
		serveClient(wsContainer, clientDir.(string))
	}

	swDir := arguments["--swagger"]
	if swDir != nil {
		config := swagger.Config{
			WebServices:     wsContainer.RegisteredWebServices(), // you control what services are visible
			WebServicesUrl:  fmt.Sprintf("http://%s", address),
			ApiPath:         "/apidocs.json",
			ApiVersion:      "0.1",
			SwaggerPath:     "/apidocs/",
			SwaggerFilePath: swDir.(string)}
		swagger.RegisterSwaggerService(config, wsContainer)
	}

	log.Printf("start listening on <%s>", address)
	server := &http.Server{Addr: address, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
