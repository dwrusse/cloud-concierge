package identifyCloudActors

import (
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestConvertProviderResourceActionsToJSON(t *testing.T) {
	ica := IdentifyCloudActors{}

	// Given
	inputProviderResourceActions := terraformValueObjects.ProviderResourceActions{
		"google": terraformValueObjects.DivisionResourceActions{
			"dragondrop-dev": map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions{
				"resource_1": {
					Creator: terraformValueObjects.CloudActorTimeStamp{
						Actor:     terraformValueObjects.CloudActor("creator@gmail.com"),
						Timestamp: terraformValueObjects.Timestamp("time_1"),
					},
					Modifier: terraformValueObjects.CloudActorTimeStamp{
						Actor:     terraformValueObjects.CloudActor("modifier@gmail.com"),
						Timestamp: terraformValueObjects.Timestamp("time_2"),
					},
				},
				"resource_2": {
					Creator: terraformValueObjects.CloudActorTimeStamp{
						Actor:     terraformValueObjects.CloudActor("el_creator@gmail.com"),
						Timestamp: terraformValueObjects.Timestamp("time_3"),
					},
				},
				"resource_3": {},
			},
		},
	}

	// When
	jsonBytes, err := ica.convertProviderResourceActionsToJSON(inputProviderResourceActions)
	if err != nil {
		t.Errorf("received unexpected error within ica.convertProviderResourceActionsToJSON:%v", err)
	}

	// Then
	expectedJSON := `{"google":{"dragondrop-dev":{"resource_1":{"creation":{"actor":"creator@gmail.com","timestamp":"time_1"},"modified":{"actor":"modifier@gmail.com","timestamp":"time_2"}},"resource_2":{"creation":{"actor":"el_creator@gmail.com","timestamp":"time_3"}}}}}`
	if expectedJSON != string(jsonBytes) {
		t.Errorf("got:\n%v\nexpected:\n%v", string(jsonBytes), expectedJSON)
	}
}
