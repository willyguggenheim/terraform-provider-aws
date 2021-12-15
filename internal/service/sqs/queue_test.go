package sqs_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	awspolicy "github.com/hashicorp/awspolicyequivalence"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfsqs "github.com/hashicorp/terraform-provider-aws/internal/service/sqs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccSQSQueue_basic(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNameConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "sqs", rName),
					resource.TestCheckResourceAttr(resourceName, "content_based_deduplication", "false"),
					resource.TestCheckResourceAttr(resourceName, "deduplication_scope", ""),
					resource.TestCheckResourceAttr(resourceName, "delay_seconds", strconv.Itoa(tfsqs.DefaultQueueDelaySeconds)),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "false"),
					resource.TestCheckResourceAttr(resourceName, "fifo_throughput_limit", ""),
					resource.TestCheckResourceAttr(resourceName, "kms_data_key_reuse_period_seconds", strconv.Itoa(tfsqs.DefaultQueueKMSDataKeyReusePeriodSeconds)),
					resource.TestCheckResourceAttr(resourceName, "kms_master_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "max_message_size", strconv.Itoa(tfsqs.DefaultQueueMaximumMessageSize)),
					resource.TestCheckResourceAttr(resourceName, "message_retention_seconds", strconv.Itoa(tfsqs.DefaultQueueMessageRetentionPeriod)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "policy", ""),
					resource.TestCheckResourceAttr(resourceName, "receive_wait_time_seconds", strconv.Itoa(tfsqs.DefaultQueueReceiveMessageWaitTimeSeconds)),
					resource.TestCheckResourceAttr(resourceName, "redrive_policy", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "url", resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", strconv.Itoa(tfsqs.DefaultQueueVisibilityTimeout)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_disappears(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNameConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					acctest.CheckResourceDisappears(acctest.Provider, tfsqs.ResourceQueue(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSQSQueue_Name_generated(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQueueNameGeneratedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					create.TestCheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_NameGenerated_fifoQueue(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		Providers:    acctest.Providers,
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQueueNameGeneratedFIFOQueueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					create.TestCheckResourceAttrNameWithSuffixGenerated(resourceName, "name", tfsqs.FIFOQueueNameSuffix),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_namePrefix(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQueueNamePrefixConfig("tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					create.TestCheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_NamePrefix_fifoQueue(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQueueNamePrefixFIFOQueueConfig("tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					create.TestCheckResourceAttrNameWithSuffixFromPrefix(resourceName, "name", "tf-acc-test-prefix-", tfsqs.FIFOQueueNameSuffix),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_tags(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTags1Config(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTags2Config(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTags1Config(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccSQSQueue_update(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNameConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "sqs", rName),
					resource.TestCheckResourceAttr(resourceName, "content_based_deduplication", "false"),
					resource.TestCheckResourceAttr(resourceName, "deduplication_scope", ""),
					resource.TestCheckResourceAttr(resourceName, "delay_seconds", strconv.Itoa(tfsqs.DefaultQueueDelaySeconds)),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "false"),
					resource.TestCheckResourceAttr(resourceName, "fifo_throughput_limit", ""),
					resource.TestCheckResourceAttr(resourceName, "kms_data_key_reuse_period_seconds", strconv.Itoa(tfsqs.DefaultQueueKMSDataKeyReusePeriodSeconds)),
					resource.TestCheckResourceAttr(resourceName, "kms_master_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "max_message_size", strconv.Itoa(tfsqs.DefaultQueueMaximumMessageSize)),
					resource.TestCheckResourceAttr(resourceName, "message_retention_seconds", strconv.Itoa(tfsqs.DefaultQueueMessageRetentionPeriod)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "policy", ""),
					resource.TestCheckResourceAttr(resourceName, "receive_wait_time_seconds", strconv.Itoa(tfsqs.DefaultQueueReceiveMessageWaitTimeSeconds)),
					resource.TestCheckResourceAttr(resourceName, "redrive_policy", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", strconv.Itoa(tfsqs.DefaultQueueVisibilityTimeout)),
				),
			},
			{
				Config: testAccUpdatedConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "sqs", rName),
					resource.TestCheckResourceAttr(resourceName, "content_based_deduplication", "false"),
					resource.TestCheckResourceAttr(resourceName, "deduplication_scope", ""),
					resource.TestCheckResourceAttr(resourceName, "delay_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "false"),
					resource.TestCheckResourceAttr(resourceName, "fifo_throughput_limit", ""),
					resource.TestCheckResourceAttr(resourceName, "kms_data_key_reuse_period_seconds", strconv.Itoa(tfsqs.DefaultQueueKMSDataKeyReusePeriodSeconds)),
					resource.TestCheckResourceAttr(resourceName, "kms_master_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "max_message_size", "2048"),
					resource.TestCheckResourceAttr(resourceName, "message_retention_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "policy", ""),
					resource.TestCheckResourceAttr(resourceName, "receive_wait_time_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "redrive_policy", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_Policy_basic(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	expectedPolicy := `
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement":[{
    "Sid": "Stmt1451501026839",
    "Effect": "Allow",
    "Principal":"*",
    "Action":"sqs:SendMessage",
    "Resource":"arn:%[1]s:sqs:%[2]s:%[3]s:%[4]s",
    "Condition":{
      "ArnEquals":{"aws:SourceArn":"arn:%[1]s:sns:%[2]s:%[3]s:%[4]s"}
    }
  }]
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					testAccCheckQueuePolicyAttribute(&queueAttributes, rName, expectedPolicy),
					resource.TestCheckResourceAttr(resourceName, "delay_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "max_message_size", "2048"),
					resource.TestCheckResourceAttr(resourceName, "message_retention_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "receive_wait_time_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_Policy_ignoreEquivalent(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	expectedPolicy := `
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement":[{
    "Sid": "SID1993561419",
    "Effect": "Allow",
    "Principal":"*",
    "Action":[
      "sqs:SendMessage",
      "sqs:DeleteMessage",
      "sqs:ListQueues"
    ],
    "Resource":"arn:%[1]s:sqs:%[2]s:%[3]s:%[4]s",
    "Condition":{
      "ArnEquals":{"aws:SourceArn":"arn:%[1]s:sns:%[2]s:%[3]s:%[4]s"}
    }
  }]
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQueuePolicyEquivalentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					testAccCheckQueuePolicyAttribute(&queueAttributes, rName, expectedPolicy),
					resource.TestCheckResourceAttr(resourceName, "delay_seconds", "90"),
					resource.TestCheckResourceAttr(resourceName, "max_message_size", "2048"),
					resource.TestCheckResourceAttr(resourceName, "message_retention_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "receive_wait_time_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "60"),
				),
			},
			{
				Config:   testAccQueuePolicyNewEquivalentConfig(rName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSQSQueue_recentlyDeleted(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNameConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					acctest.CheckResourceDisappears(acctest.Provider, tfsqs.ResourceQueue(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccNameConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
				),
			},
		},
	})
}

func TestAccSQSQueue_redrivePolicy(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedrivePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "delay_seconds", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "redrive_policy"),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "300"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_fifoQueue(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := fmt.Sprintf("%s.fifo", sdkacctest.RandomWithPrefix(acctest.ResourcePrefix))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFIFOQueueConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "deduplication_scope", "queue"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "true"),
					resource.TestCheckResourceAttr(resourceName, "fifo_throughput_limit", "perQueue"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_FIFOQueue_expectNameError(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccFIFOQueueConfig(rName),
				ExpectError: regexp.MustCompile(`invalid queue name:`),
			},
		},
	})
}

func TestAccSQSQueue_FIFOQueue_contentBasedDeduplication(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := fmt.Sprintf("%s.fifo", sdkacctest.RandomWithPrefix(acctest.ResourcePrefix))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFIFOQueueContentBasedDeduplicationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "content_based_deduplication", "true"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSQSQueue_FIFOQueue_highThroughputMode(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := fmt.Sprintf("%s.fifo", sdkacctest.RandomWithPrefix(acctest.ResourcePrefix))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFIFOQueueHighThroughputModeConfig(rName, "null", "null"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "deduplication_scope", "queue"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "true"),
					resource.TestCheckResourceAttr(resourceName, "fifo_throughput_limit", "perQueue"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFIFOQueueHighThroughputModeConfig(rName, "messageGroup", "perMessageGroupId"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "deduplication_scope", "messageGroup"),
					resource.TestCheckResourceAttr(resourceName, "fifo_queue", "true"),
					resource.TestCheckResourceAttr(resourceName, "fifo_throughput_limit", "perMessageGroupId"),
				),
			},
		},
	})
}

func TestAccSQSQueue_StandardQueue_expectContentBasedDeduplicationError(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccStandardQueueExpectContentBasedDeduplicationErrorConfig(rName),
				ExpectError: regexp.MustCompile(`content-based deduplication can only be set for FIFO queue`),
			},
		},
	})
}

func TestAccSQSQueue_encryption(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEncryptionConfig(rName, "null"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "kms_data_key_reuse_period_seconds", "300"),
					resource.TestCheckResourceAttr(resourceName, "kms_master_key_id", "alias/aws/sqs"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEncryptionConfig(rName, "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "kms_data_key_reuse_period_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "kms_master_key_id", "alias/aws/sqs"),
				),
			},
			{
				Config: testAccManagedEncryptionConfig(rName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "sqs_managed_sse_enabled", "true"),
				),
			},
		},
	})
}

func TestAccSQSQueue_zeroVisibilityTimeoutSeconds(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZeroVisibilityTimeoutSecondsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "visibility_timeout_seconds", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// https://github.com/hashicorp/terraform-provider-aws/issues/19786.
func TestAccSQSQueue_defaultKMSDataKeyReusePeriodSeconds(t *testing.T) {
	var queueAttributes map[string]string
	resourceName := "aws_sqs_queue.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDefaultKMSDataKeyReusePeriodSecondsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQueueExists(resourceName, &queueAttributes),
					resource.TestCheckResourceAttr(resourceName, "kms_data_key_reuse_period_seconds", strconv.Itoa(tfsqs.DefaultQueueKMSDataKeyReusePeriodSeconds)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckQueuePolicyAttribute(queueAttributes *map[string]string, rName, policyTemplate string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expectedPolicy := fmt.Sprintf(policyTemplate, acctest.Partition(), acctest.Region(), acctest.AccountID(), rName)

		var actualPolicyText string
		for key, value := range *queueAttributes {
			if key == sqs.QueueAttributeNamePolicy {
				actualPolicyText = value
				break
			}
		}

		equivalent, err := awspolicy.PoliciesAreEquivalent(actualPolicyText, expectedPolicy)
		if err != nil {
			return fmt.Errorf("Error testing policy equivalence: %s", err)
		}
		if !equivalent {
			return fmt.Errorf("Non-equivalent policy error:\n\nexpected: %s\n\n     got: %s\n", expectedPolicy, actualPolicyText)
		}

		return nil
	}
}

func testAccCheckQueueExists(resourceName string, v *map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SQS Queue URL is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SQSConn

		output, err := tfsqs.FindQueueAttributesByURL(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = output

		return nil
	}
}

func testAccCheckQueueDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SQSConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sqs_queue" {
			continue
		}

		_, err := tfsqs.FindQueueAttributesByURL(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("SQS Queue %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccQueueNameGeneratedConfig = `
resource "aws_sqs_queue" "test" {}
`

const testAccQueueNameGeneratedFIFOQueueConfig = `
resource "aws_sqs_queue" "test" {
  fifo_queue = true
}
`

func testAccNameConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name = %[1]q
}
`, rName)
}

func testAccQueueNamePrefixConfig(prefix string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name_prefix = %[1]q
}
`, prefix)
}

func testAccQueueNamePrefixFIFOQueueConfig(prefix string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name_prefix = %[1]q
  fifo_queue  = true
}
`, prefix)
}

func testAccTags1Config(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccTags2Config(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccUpdatedConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                       = %[1]q
  delay_seconds              = 90
  max_message_size           = 2048
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10
  visibility_timeout_seconds = 60
}
`, rName)
}

func testAccPolicyConfig(rName string) string {
	return fmt.Sprintf(`
locals {
  queue_name = %[1]q
}

resource "aws_sns_topic" "test" {
  name = %[1]q
}

data "aws_partition" "current" {}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_sqs_queue" "test" {
  name                       = local.queue_name
  delay_seconds              = 90
  max_message_size           = 2048
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10
  visibility_timeout_seconds = 60

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "Stmt1451501026839",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "arn:${data.aws_partition.current.partition}:sqs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${local.queue_name}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_sns_topic.test.arn}"
        }
      }
    }
  ]
}
EOF
}

resource "aws_sns_topic_subscription" "test" {
  topic_arn = aws_sns_topic.test.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.test.arn
}
`, rName)
}

func testAccQueuePolicyEquivalentConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "test" {
  name = %[1]q
}

data "aws_partition" "current" {}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_sqs_queue" "test" {
  name                       = %[1]q
  delay_seconds              = 90
  max_message_size           = 2048
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10
  visibility_timeout_seconds = 60

  policy = jsonencode({
    Version = "2012-10-17"
    Id      = "sqspolicy"
    Statement = [{
      Sid       = "SID1993561419"
      Effect    = "Allow"
      Principal = "*"
      Action = [
        "sqs:SendMessage",
        "sqs:DeleteMessage",
        "sqs:ListQueues",
      ]
      Resource = "arn:${data.aws_partition.current.partition}:sqs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:%[1]s"
      Condition = {
        ArnEquals = {
          "aws:SourceArn" = aws_sns_topic.test.arn
        }
      }
    }]
  })
}

resource "aws_sns_topic_subscription" "test" {
  topic_arn = aws_sns_topic.test.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.test.arn
}
`, rName)
}

func testAccQueuePolicyNewEquivalentConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "test" {
  name = %[1]q
}

data "aws_partition" "current" {}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_sqs_queue" "test" {
  name                       = %[1]q
  delay_seconds              = 90
  max_message_size           = 2048
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10
  visibility_timeout_seconds = 60

  policy = jsonencode({
    Version = "2012-10-17"
    Id      = "sqspolicy"
    Statement = [{
      Sid       = "SID1993561419"
      Effect    = "Allow"
      Principal = ["*"]
      Action = [
        "sqs:ListQueues",
        "sqs:SendMessage",
        "sqs:DeleteMessage",
      ]
      Resource = "arn:${data.aws_partition.current.partition}:sqs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:%[1]s"
      Condition = {
        ArnEquals = {
          "aws:SourceArn" = aws_sns_topic.test.arn
        }
      }
    }]
  })
}

resource "aws_sns_topic_subscription" "test" {
  topic_arn = aws_sns_topic.test.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.test.arn
}
`, rName)
}

func testAccRedrivePolicyConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                       = "%[1]s-1"
  delay_seconds              = 0
  visibility_timeout_seconds = 300

  redrive_policy = <<EOF
{
  "maxReceiveCount": 3,
  "deadLetterTargetArn": "${aws_sqs_queue.dlq.arn}"
}
EOF
}

resource "aws_sqs_queue" "dlq" {
  name = "%[1]s-2"
}
`, rName)
}

func testAccFIFOQueueConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name       = %[1]q
  fifo_queue = true
}
`, rName)
}

func testAccFIFOQueueContentBasedDeduplicationConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                        = %[1]q
  fifo_queue                  = true
  content_based_deduplication = true
}
`, rName)
}

func testAccFIFOQueueHighThroughputModeConfig(rName, deduplicationScope, fifoThroughputLimit string) string {
	if deduplicationScope != "null" {
		deduplicationScope = strconv.Quote(deduplicationScope)
	}

	if fifoThroughputLimit != "null" {
		fifoThroughputLimit = strconv.Quote(fifoThroughputLimit)
	}

	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name       = %[1]q
  fifo_queue = true

  deduplication_scope   = %[2]s
  fifo_throughput_limit = %[3]s
}
`, rName, deduplicationScope, fifoThroughputLimit)
}

func testAccStandardQueueExpectContentBasedDeduplicationErrorConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                        = %[1]q
  content_based_deduplication = true
}
`, rName)
}

func testAccEncryptionConfig(rName, kmsDataKeyReusePeriodSeconds string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                              = %[1]q
  kms_master_key_id                 = "alias/aws/sqs"
  kms_data_key_reuse_period_seconds = %[2]s
}
`, rName, kmsDataKeyReusePeriodSeconds)
}

func testAccManagedEncryptionConfig(rName, sqsManagedSseEnabled string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                    = %[1]q
  sqs_managed_sse_enabled = %[2]s
}
`, rName, sqsManagedSseEnabled)
}

func testAccZeroVisibilityTimeoutSecondsConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                       = %[1]q
  visibility_timeout_seconds = 0
}
`, rName)
}

func testAccDefaultKMSDataKeyReusePeriodSecondsConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name                              = %[1]q
  kms_data_key_reuse_period_seconds = 300
}
`, rName)
}
