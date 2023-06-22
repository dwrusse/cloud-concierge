package terraformerExecutor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/driftmitigation/hclcreate"
	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
	terraformerCli "github.com/dragondrop-cloud/driftmitigation/implementations/terraformer_executor/terraformer_cli"
	"github.com/dragondrop-cloud/driftmitigation/interfaces"
)

func TestCreateIsolatedTerraformerExecutor(t *testing.T) {
	// Given
	ctx := context.Background()
	hclConfig := hclcreate.Config{}
	executorConfig := terraformerCli.TerraformerExecutorConfig{}
	cliConfig := terraformerCli.Config{}
	terraformerExecutorProvider := "isolated"
	terraformerExecutorFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	terraformerExecutor, err := terraformerExecutorFactory.Instantiate(ctx, terraformerExecutorProvider, dragonDrop, divisionToProvider, hclConfig, executorConfig, cliConfig)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, terraformerExecutor)
}
