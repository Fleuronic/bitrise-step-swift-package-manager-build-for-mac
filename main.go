package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/kitasuke/go-swift/swift"
	"github.com/kitasuke/go-swift/utility"
)

const (
	BuildPathEnvKey      = "build_path"
	BuildTestsEnvKey     = "build_tests"
	DisableSandboxEnvKey = "disable_sandbox"
)

// ConfigModel ...
type ConfigModel struct {
	// Project Parameters
	BuildPath string

	// Build Run Configs
	buildTests     string
	disableSandbox string
}

func (configs ConfigModel) print() {
	fmt.Println()

	log.Infof("Project Parameters:")
	log.Printf("- BuildPath: %s", configs.BuildPath)

	fmt.Println()
	log.Infof("Build Run Configs:")
	log.Printf("- BuildTests: %s", configs.buildTests)
	log.Printf("- DisableSandbox: %s", configs.disableSandbox)
}

func createConfigsModelFromEnvs() ConfigModel {
	return ConfigModel{
		// Project Parameters
		BuildPath: os.Getenv(BuildPathEnvKey),

		// Test Run Configs
		buildTests:     os.Getenv(BuildTestsEnvKey),
		disableSandbox: os.Getenv(DisableSandboxEnvKey),
	}
}

func (configs ConfigModel) validate() error {
	if err := validateRequiredInputWithOptions(configs.buildTests, BuildTestsEnvKey, []string{"yes", "no"}); err != nil {
		return err
	}

	if err := validateRequiredInputWithOptions(configs.disableSandbox, DisableSandboxEnvKey, []string{"yes", "no"}); err != nil {
		return err
	}

	return nil
}

//--------------------
// Functions
//--------------------

func validateRequiredInput(value, key string) error {
	if value == "" {
		return fmt.Errorf("Missing required input: %s", key)
	}
	return nil
}

func validateRequiredInputWithOptions(value, key string, options []string) error {
	if err := validateRequiredInput(value, key); err != nil {
		return err
	}

	found := false
	for _, option := range options {
		if option == value {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Invalid input: (%s) value: (%s), valid options: %s", key, value, strings.Join(options, ", "))
	}

	return nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

//--------------------
// Main
//--------------------

func main() {
	configs := createConfigsModelFromEnvs()
	configs.print()
	if err := configs.validate(); err != nil {
		failf("Issue with input: %s", err)
	}

	fmt.Println()
	log.Infof("Other Configs:")

	buildTests := configs.buildTests == "yes"
	disableSandbox := configs.disableSandbox == "yes"

	swiftVersion, err := utility.GetSwiftVersion()
	if err != nil {
		failf("Failed to get the version of swift! Error: %s", err)
	}

	log.Printf("* swift_version: %s (%s)", swiftVersion.Version, swiftVersion.Target)

	fmt.Println()

	// setup CommandModel for test
	buildCommandModel := swift.NewBuildCommand()
	buildCommandModel.SetBuildPath(configs.BuildPath)
	buildCommandModel.SetBuildTests(buildTests)
	buildCommandModel.SetDisableSandbox(disableSandbox)

	log.Infof("$ %s\n", buildCommandModel.PrintableCmd())

	if err := buildCommandModel.Run(); err != nil {
		failf("Build failed, error: %s", err)
	}
}
