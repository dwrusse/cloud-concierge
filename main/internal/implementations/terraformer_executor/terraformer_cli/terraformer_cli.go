package terraformerCLI

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// TerraformImportMigrationGeneratorParams is the struct that wraps the params to run terraform import statement
type TerraformImportMigrationGeneratorParams struct {
	Provider       string
	Division       terraformValueObjects.Division
	Resources      []string
	AdditionalArgs []string
	Regions        []string
	IsCompact      bool
}

// TerraformerCLI interface is an abstraction on the methods needed within the
// terraformer package.
type TerraformerCLI interface {

	// Import runs the `terraformer import` command to import all resources for the specified provider, division,
	// and specified credentials.
	Import(params TerraformImportMigrationGeneratorParams) (terraformValueObjects.Path, error)

	// UpdateState runs the `terraform state replace-provider` command to upgrade the state file generated
	// to version 4.
	UpdateState(provider string, location string) error
}

// Config is the struct that contains parameters considered to import the resources
type Config struct {
	// ResourcesWhiteList represents the list of resource names that will be exclusively considered for inclusion in the import statement.
	ResourcesWhiteList terraformValueObjects.ResourceNameList

	// ResourcesBlackList represents the list of resource names that will be excluded from consideration for inclusion in the import statement.
	ResourcesBlackList terraformValueObjects.ResourceNameList
}

// terraformerCLI implements the TerraformerCLI interface.
type terraformerCLI struct {
	// config is the struct that contains parameters considered to import the resources such the black and white resources list
	config Config
}

// newTerraformerCLI creates a new instance of the terraformerCLI struct.
func newTerraformerCLI(config Config) TerraformerCLI {
	return &terraformerCLI{config: config}
}

// Import runs the `terraformer import` command.
func (tfrCLI *terraformerCLI) Import(params TerraformImportMigrationGeneratorParams) (terraformValueObjects.Path, error) {
	divisionOutput := fmt.Sprintf("--path-output=./%s-%v", params.Provider, params.Division)

	importProvider := getActualImportProvider(params.Provider)
	mainArgs := []string{
		"import", importProvider,
		fmt.Sprintf("--compact=%s", strconv.FormatBool(params.IsCompact)),
		divisionOutput,
		"--path-pattern={output}",
	}

	if len(params.Regions) > 0 {
		regions := strings.Join(params.Regions, ",")
		mainArgs = append(mainArgs, fmt.Sprintf("--regions=%s", regions))
	}

	if len(tfrCLI.config.ResourcesBlackList) > 0 {
		resourceGroups := tfrCLI.getGroupListByResourceNames(tfrCLI.config.ResourcesBlackList)

		if len(resourceGroups) > 0 {
			excludes := strings.Join(resourceGroups, ",")
			mainArgs = append(mainArgs, fmt.Sprintf("--excludes=%s", excludes))
			mainArgs = append(mainArgs, "--resources=*")
		}
	} else if len(tfrCLI.config.ResourcesWhiteList) > 0 {
		resourceGroups := tfrCLI.getGroupListByResourceNames(tfrCLI.config.ResourcesWhiteList)

		if len(resourceGroups) > 0 {
			resources := strings.Join(tfrCLI.getGroupListByResourceNames(tfrCLI.config.ResourcesWhiteList), ",")
			mainArgs = append(mainArgs, fmt.Sprintf("--resources=%s", resources))
		}
	} else {
		mainArgs = append(mainArgs, "--resources=*")
	}

	args := append(mainArgs, params.AdditionalArgs...)
	log.Infof("Terraformer ARGS: %s", args)
	err := executeCommand("terraformer", args...)

	if err != nil {
		return "", fmt.Errorf("[Import] Error in running 'terraformer import': %v", err)
	}
	return terraformValueObjects.Path(fmt.Sprintf("./%s-%v/", params.Provider, params.Division)), nil
}

func getActualImportProvider(provider string) string {
	if provider == "azurerm" {
		return "azure"
	}

	return provider
}

func (tfrCLI *terraformerCLI) UpdateState(provider string, location string) error {
	// Specify the location of the state file, as well as the from and to provider plug in values.
	stateFlag := fmt.Sprintf("-state=%s/terraform.tfstate", location)
	fromProvider := fmt.Sprintf("registry.terraform.io/-/%s", provider)
	toProvider := fmt.Sprintf("hashicorp/%s", provider)

	args := []string{"state", "replace-provider", "-auto-approve", stateFlag, fromProvider, toProvider}

	err := executeCommand("terraform", args...)
	if err != nil {
		return fmt.Errorf("[UpdateState] Error in running 'terraform state replace-provider': %v", err)
	}
	return nil
}

func (tfrCLI *terraformerCLI) getGroupListByResourceNames(list terraformValueObjects.ResourceNameList) []string {
	resourceGroups := make([]string, 0)

	for _, resourceName := range list {
		resourceGroup := googleResourceGroups[resourceName]
		if resourceGroup == "" {
			resourceGroup = awsResourceGroups[resourceName]
		}
		if resourceGroup == "" {
			resourceGroup = azureResourceGroups[resourceName]
		}
		resourceGroups = append(resourceGroups, resourceGroup)
	}

	return resourceGroups
}

// executeCommand wraps os.exec.Command with capturing of std output and errors.
func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	// Setting up logging objects
	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("%v\n\n%v", err, stderr.String()+out.String())
	}
	fmt.Printf("\n%s Output:\n\n%v\n", command, out.String())
	return nil
}
