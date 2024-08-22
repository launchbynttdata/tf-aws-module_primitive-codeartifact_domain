package testimpl

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	base            = "../../examples/"
	testVarFileName = "/test.tfvars"
)

var standardTags = map[string]string{
	"provisioner": "Terraform",
}

func TestCodeArtifact(t *testing.T, ctx types.TestContext) {
	t.Parallel()
	stage := test_structure.RunTestStage

	files, err := os.ReadDir(base)
	assert.NoError(t, err)
	basePath, _ := filepath.Abs(base)
	for _, file := range files {
		dir := filepath.Join(basePath, file.Name())
		if file.IsDir() {
			defer stage(t, "teardown_codeartifact", func() { tearDownCodeArtifact(t, dir) })
			stage(t, "setup_codeartifact", func() { setupCodeArtifactTest(t, dir) })
			stage(t, "test_codeartifact", func() { testCodeArtifact(t, dir) })
		}
	}
}

func setupCodeArtifactTest(t *testing.T, dir string) {

	terraformOptions := &terraform.Options{
		TerraformDir: dir,
		VarFiles:     []string{dir + testVarFileName},
		NoColor:      true,
		Logger:       logger.Discard,
	}

	test_structure.SaveTerraformOptions(t, dir, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

}

func testCodeArtifact(t *testing.T, dir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, dir)
	terraformOptions.Logger = logger.Discard

	expectedPatternARN := "^arn:aws:codeartifact:[a-z0-9-]+:[0-9]{12}:[a-z0-9-]+"

	actualARN := terraform.Output(t, terraformOptions, "arn")
	assert.NotEmpty(t, actualARN, "ARN is empty")
	assert.Regexp(t, expectedPatternARN, actualARN, "ARN does not match expected pattern")

	client := GetAWSCodeartifactClient(t)
	tfvarsFullPath := dir + testVarFileName

	expectedDomainName, err := terraform.GetVariableAsStringFromVarFileE(t, tfvarsFullPath, "domain")
	assert.NoError(t, err)

	input := &codeartifact.DescribeDomainInput{
		Domain: aws.String(expectedDomainName),
	}

	result, err := client.DescribeDomain(context.TODO(), input)
	assert.NoError(t, err, "The expected code artifact domain was not found")

	domain := result.Domain

	actualName := domain.Name
	assert.Equal(t, expectedDomainName, *actualName, "Domain Name does not match")
	checkTagsMatch(t, tfvarsFullPath, actualARN, client)
}

func checkTagsMatch(t *testing.T, tfvarsFullPath string, actualARN string, client *codeartifact.Client) {
	expectedTags, err := terraform.GetVariableAsMapFromVarFileE(t, tfvarsFullPath, "tags")
	assert.NoError(t, err)

	input := &codeartifact.ListTagsForResourceInput{
		ResourceArn: aws.String(actualARN),
	}
	result, err := client.ListTagsForResource(context.TODO(), input)
	assert.NoError(t, err, "Failed to retrieve tags from AWS")
	// convert AWS Tag[] to map so we can compare
	actualTags := map[string]string{}
	for _, tag := range result.Tags {
		actualTags[*tag.Key] = *tag.Value
	}

	// add the standard tags and the resource_name tag to the expected tags
	for k, v := range standardTags {
		expectedTags[k] = v
	}

	assert.True(t, reflect.DeepEqual(actualTags, expectedTags), fmt.Sprintf("tags did not match, expected: %v\nactual: %v", expectedTags, actualTags))
}

func tearDownCodeArtifact(t *testing.T, dir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, dir)
	terraformOptions.Logger = logger.Discard
	terraform.Destroy(t, terraformOptions)
}

func GetAWSCodeartifactClient(t *testing.T) *codeartifact.Client {
	ecrClient := codeartifact.NewFromConfig(GetAWSConfig(t))
	return ecrClient
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
