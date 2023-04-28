package dms_test

import (
	"fmt"
	"testing"

	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccDMSReplicationInstanceDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	replicationInstanceClass := "dms.c4.large"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dms_replication_instance.test"
	dataSourceName := "data.aws_dms_replication_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, dms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckReplicationInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicationInstanceDataSourceConfig_basic(rName, replicationInstanceClass),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReplicationInstanceExists(ctx, dataSourceName),
					resource.TestCheckResourceAttrPair(dataSourceName, "replication_instance_id", resourceName, "replication_instance_id"),
				),
			},
		},
	})
}

func testAccReplicationInstanceDataSourceConfig_basic(rName, replicationInstanceClass string) string {
	return fmt.Sprintf(`
resource "aws_dms_replication_instance" "test" {
  apply_immediately          = true
  replication_instance_class = %q
  replication_instance_id    = %q
}

data "aws_dms_replication_instance" "test" {
  replication_instance_id = aws_dms_replication_instance.test.replication_instance_id
}
`, replicationInstanceClass, rName)
}
