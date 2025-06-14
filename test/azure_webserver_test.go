package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "bd7bd383-5df6-4911-a820-52a7dc9cdb5b"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "niyo0041",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	t.Run("Confirm VM exists", func(t *testing.T) {
		assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
	})

	t.Run("Confirm NIC exists", func(t *testing.T) {
		nicName := terraform.Output(t, terraformOptions, "nic_name")
		assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))
	})

	t.Run("Confirm NIC is connected to the VM", func(t *testing.T) {
		nicName := terraform.Output(t, terraformOptions, "nic_name")
		vmNicsList := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
		assert.Contains(t, vmNicsList, nicName, "The NIC should be connected to the VM")
	})

	t.Run("Confirm correct Ubuntu version", func(t *testing.T) {
		expectedSKU := "22_04-lts-gen2"
		vmImageStruct := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedSKU, vmImageStruct.SKU, "The VM should be running the correct Ubuntu version")
	})

}
