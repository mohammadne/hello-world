package cmd

import (
	"errors"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/mohammadne/hello-world/internal/entities"
	"github.com/mohammadne/hello-world/internal/generator"
	"github.com/mohammadne/hello-world/pkg/validator"
	"github.com/spf13/cobra"
)

const (
	defaultXrayVersion     = "v1.18.11"
	defaultOutputDirectory = "outputs"
	defaultServerPort      = 443
	defaultClientPort      = 10809
)

type Executer struct {
	outputDirectory string

	serverAddress string
	serverPort    int
	clientPort    int

	// ask from the user
	version  string
	protocol entities.Protocol
}

func (executer Executer) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "executer",
		Short: "execute the program and generate the configuration files based on given prompts",
		Run: func(_ *cobra.Command, _ []string) {
			executer.selectProtocol()
			executer.promptServerDomainOrIP()

			generator.Generator{
				Protocol:      executer.protocol,
				Version:       executer.version,
				ServerAddress: executer.serverAddress,
				ServerPort:    executer.serverPort,
				ClientPort:    executer.clientPort,
			}.Run(executer.outputDirectory)
		},
	}

	cmd.Flags().StringVar(&executer.version, "xray-version", defaultXrayVersion, "The xray-core version")
	cmd.Flags().StringVar(&executer.outputDirectory, "output-directory", defaultOutputDirectory, "The output directory for generated files")
	cmd.Flags().IntVar(&executer.serverPort, "server-port", defaultServerPort, "The port server exposed from")
	cmd.Flags().IntVar(&executer.serverPort, "client-port", defaultClientPort, "The port client exposed from")

	return cmd
}

func (executer *Executer) selectProtocol() {
	prompt := promptui.Select{
		Label: "Select which Xray protocol you want for your configuration",
		Items: []entities.Protocol{
			entities.Reality,
		},
	}

	_, rawProtocol, err := prompt.Run()
	if err != nil {
		log.Fatalf("protocol prompt has been failed %v\n", err)
	} else if entities.ValidateProtocol(rawProtocol) != nil {
		log.Fatalf("invalid protocol value has been given %s \n%v\n", rawProtocol, err)
	}

	executer.protocol = entities.Protocol(rawProtocol)
}

func (executer *Executer) promptServerDomainOrIP() {
	validate := func(input string) error {
		if validator.ValidateIP(input) != nil && validator.ValidateDomain(input) != nil {
			return errors.New("Invalid Domain or IP has been given")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter Domain or IPv4 address of the server",
		Validate: validate,
	}

	address, err := prompt.Run()
	if err != nil {
		log.Fatalf("protocol prompt has been failed %v\n", err)
	}

	executer.serverAddress = address
}
