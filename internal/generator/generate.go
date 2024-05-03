package generator

import (
	"embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/mohammadne/hello-world/internal/entities"
	"github.com/mohammadne/hello-world/pkg/cryptography"
)

type Generator struct {
	Protocol      entities.Protocol
	Version       string
	ServerAddress string
	ServerPort    int
	ClientPort    int

	serverOutputDirectory string
	clientOutputDirectory string
	keyPair               *cryptography.KeyPair
	uuid                  string
}

const (
	TemplatesDirectory = "templates/"
)

//go:embed templates
var templates embed.FS

func (generator Generator) Run(outputDirectory string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error findout pwd: %v\n", err)
	}
	outputDirectory = fmt.Sprintf("%s/%s/", wd, outputDirectory)

	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		log.Fatalf("error creating directory %s, %v\n", outputDirectory, err)
	}

	generator.serverOutputDirectory = outputDirectory + "/server/"
	generator.clientOutputDirectory = outputDirectory + "/client/"

	generator.uuid = cryptography.GenerateUUID()

	keyPair, err := cryptography.GenerateCurve25519Keys()
	if err != nil {
		log.Fatalf("error creating Curve25519 key pairs %v\n", err)
	}
	generator.keyPair = keyPair

	for _, machine := range []entities.Machine{entities.Server, entities.Client} {
		path := generator.outputPath(machine)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("error creating directory %s, %v\n", path, err)
		}

		generator.generateDockerfile(machine)
		generator.generateConfig(machine)
		generator.generateDockerCompose(machine)
	}
}

func (generator *Generator) outputPath(machine entities.Machine) string {
	if machine == entities.Server {
		return generator.serverOutputDirectory
	}
	return generator.clientOutputDirectory
}

const (
	DockerfileTemplateFile = "Dockerfile.tmpl"
	DockerfileOutputFile   = "Dockerfile"
)

func (generator *Generator) generateDockerfile(machine entities.Machine) {
	templateFile := TemplatesDirectory + DockerfileTemplateFile
	tmpl, err := template.ParseFS(templates, templateFile)
	if err != nil {
		panic(err)
	}

	var outputFile *os.File
	outputFile, err = os.Create(generator.outputPath(machine) + DockerfileOutputFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	values := struct{ XrayVersion string }{XrayVersion: generator.Version}

	if err = tmpl.Execute(outputFile, values); err != nil {
		panic(err)
	}
}

const (
	ConfigServerTemplateFile = "config-server.tmpl.json"
	ConfigClientTemplateFile = "config-client.tmpl.json"
	ConfigOutputFile         = "config.json"
)

func (generator *Generator) generateConfig(machine entities.Machine) {
	var templateFile string
	if machine == entities.Server {
		templateFile = TemplatesDirectory + ConfigServerTemplateFile
	} else {
		templateFile = TemplatesDirectory + ConfigClientTemplateFile
	}

	tmpl, err := template.ParseFS(templates, templateFile)
	if err != nil {
		panic(err)
	}

	var outputFile *os.File
	outputFile, err = os.Create(generator.outputPath(machine) + ConfigOutputFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	var values any
	if machine == entities.Server {
		values = struct {
			Port       int
			UUID       string
			PrivateKey string
		}{
			Port:       generator.ServerPort,
			UUID:       generator.uuid,
			PrivateKey: generator.keyPair.PrivateKey,
		}
	} else {
		values = struct {
			ClientPort    int
			ServerAddress string
			ServerPort    int
			UUID          string
			PublicKey     string
		}{
			ClientPort:    generator.ClientPort,
			ServerAddress: generator.ServerAddress,
			ServerPort:    generator.ServerPort,
			UUID:          generator.uuid,
			PublicKey:     generator.keyPair.PublicKey,
		}
	}

	if err = tmpl.Execute(outputFile, values); err != nil {
		panic(err)
	}
}

const (
	DockerComposeTemplateFile = "docker-compose.tmpl.yml"
	DockerComposeOutputFile   = "docker-compose.yml"
)

func (generator *Generator) generateDockerCompose(machine entities.Machine) {
	templateFile := TemplatesDirectory + DockerComposeTemplateFile
	tmpl, err := template.ParseFS(templates, templateFile)
	if err != nil {
		panic(err)
	}

	values := struct{ Port int }{Port: generator.ClientPort}
	if machine == entities.Server {
		values.Port = generator.ServerPort
	}

	var outputFile *os.File
	outputFile, err = os.Create(generator.outputPath(machine) + DockerComposeOutputFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if err = tmpl.Execute(outputFile, values); err != nil {
		panic(err)
	}
}
