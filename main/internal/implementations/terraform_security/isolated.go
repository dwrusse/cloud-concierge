package terraformSecurity

import "context"

// IsolatedTerraformSecurity is an isolated implementation from interfaces.TerraformSecurity
type IsolatedTerraformSecurity struct {
}

// NewIsolatedTerraformSecurity generates an instance from IsolatedTerraformSecurity
func NewIsolatedTerraformSecurity() *IsolatedTerraformSecurity {
	return &IsolatedTerraformSecurity{}
}

// ExecuteScan is called from the main job flow to mock the output files from tfsec to show
// to the user in the PR
func (i *IsolatedTerraformSecurity) ExecuteScan(ctx context.Context) error {
	return nil
}
