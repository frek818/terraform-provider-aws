// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backup

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	awstypes "github.com/aws/aws-sdk-go-v2/service/backup/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/sdkv2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_backup_vault_policy", name="Vault Policy")
func resourceVaultPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVaultPolicyPut,
		ReadWithoutTimeout:   resourceVaultPolicyRead,
		UpdateWithoutTimeout: resourceVaultPolicyPut,
		DeleteWithoutTimeout: resourceVaultPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"backup_vault_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_vault_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrPolicy: sdkv2.IAMPolicyDocumentSchemaRequired(),
		},
	}
}

func resourceVaultPolicyPut(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BackupClient(ctx)

	policy, err := structure.NormalizeJsonString(d.Get(names.AttrPolicy).(string))
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	name := d.Get("backup_vault_name").(string)
	input := &backup.PutBackupVaultAccessPolicyInput{
		BackupVaultName: aws.String(name),
		Policy:          aws.String(policy),
	}

	_, err = tfresource.RetryWhenAWSErrMessageContains(ctx, propagationTimeout,
		func() (any, error) {
			return conn.PutBackupVaultAccessPolicy(ctx, input)
		},
		errCodeInvalidParameterValueException, "Provided principal is not valid",
	)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Backup Vault Policy (%s): %s", name, err)
	}

	d.SetId(name)

	return append(diags, resourceVaultPolicyRead(ctx, d, meta)...)
}

func resourceVaultPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BackupClient(ctx)

	output, err := findVaultAccessPolicyByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Backup Vault Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Backup Vault Policy (%s): %s", d.Id(), err)
	}

	d.Set("backup_vault_arn", output.BackupVaultArn)
	d.Set("backup_vault_name", output.BackupVaultName)

	policyToSet, err := verify.SecondJSONUnlessEquivalent(d.Get(names.AttrPolicy).(string), aws.ToString(output.Policy))
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	policyToSet, err = structure.NormalizeJsonString(policyToSet)
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	d.Set(names.AttrPolicy, policyToSet)

	return diags
}

func resourceVaultPolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).BackupClient(ctx)

	log.Printf("[DEBUG] Deleting Backup Vault Policy (%s)", d.Id())
	input := backup.DeleteBackupVaultAccessPolicyInput{
		BackupVaultName: aws.String(d.Id()),
	}
	_, err := conn.DeleteBackupVaultAccessPolicy(ctx, &input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) || tfawserr.ErrCodeEquals(err, errCodeAccessDeniedException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Backup Vault Policy (%s): %s", d.Id(), err)
	}

	return diags
}

func findVaultAccessPolicyByName(ctx context.Context, conn *backup.Client, name string) (*backup.GetBackupVaultAccessPolicyOutput, error) {
	input := &backup.GetBackupVaultAccessPolicyInput{
		BackupVaultName: aws.String(name),
	}

	output, err := conn.GetBackupVaultAccessPolicy(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) || tfawserr.ErrCodeEquals(err, errCodeAccessDeniedException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
