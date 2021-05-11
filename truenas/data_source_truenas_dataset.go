package truenas

import (
	"context"
	"github.com/dariusbakunas/terraform-provider-truenas/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
)

func dataSourceTrueNASDataset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrueNASDatasetRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"acl_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Determine how chmod behaves when adjusting file ACLs. See the zfs(8) aclmode property.",
				Computed:    true,
			},
			"acl_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"atime": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Choose 'on' to update the access time for files when they are read. Choose 'off' to prevent producing log traffic when reading files",
				Computed:    true,
			},
			"case_sensitivity": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"comments": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"compression": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"copies": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"deduplication": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encryption_root": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"inherit_encryption": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encryption_algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_loaded": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
			},
			"pbkdf2iters": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			//"passphrase": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"encryption_key": &schema.Schema{
			//	Type:     schema.TypeString,
			//  Sensitive: true,
			//	Computed: true,
			//},
			//"generate_key": &schema.Schema{
			//	Type:     schema.TypeBool,
			//	Computed: true,
			//},
			"exec": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_format": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"managed_by": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_point": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_bytes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_critical": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_warning": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ref_quota_bytes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ref_quota_critical": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ref_quota_warning": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ref_reservation": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"readonly": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_size": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_size_bytes": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
			},
			"reservation": &schema.Schema{
				Type: schema.TypeInt,
				Computed: true,
			},
			"sync": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"snap_dir": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"xattr": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"origin": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTrueNASDatasetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*api.Client)
	id := d.Get("id").(string)

	resp, err := c.DatasetAPI.Get(ctx, id)

	if err != nil {
		return diag.Errorf("error getting dataset: %s", err)
	}

	dpath := newDatasetPath(id)

	d.Set("pool", dpath.Pool)

	if dpath.Parent != "" {
		d.Set("parent", dpath.Parent)
	}

	d.Set("name", dpath.Name)
	d.Set("mount_point", resp.MountPoint)
	d.Set("encryption_root", resp.EncryptionRoot)
	d.Set("key_loaded", resp.KeyLoaded)

	if resp.ACLMode != nil {
		if err := d.Set("acl_mode", strings.ToLower(*resp.ACLMode.Value)); err != nil {
			return diag.Errorf("error setting acl_mode: %s", err)
		}
	}

	if resp.ACLType != nil {
		if err := d.Set("acl_type", strings.ToLower(*resp.ACLType.Value)); err != nil {
			return diag.Errorf("error setting acl_type: %s", err)
		}
	}

	if resp.ATime != nil {
		if err := d.Set("atime", strings.ToLower(*resp.ATime.Value)); err != nil {
			return diag.Errorf("error setting atime: %s", err)
		}
	}

	if resp.CaseSensitivity != nil {
		if err := d.Set("case_sensitivity", strings.ToLower(*resp.CaseSensitivity.Value)); err != nil {
			return diag.Errorf("error setting case_sensitivity: %s", err)
		}
	}

	if resp.Comments != nil {
		// TrueNAS does not seem to change comments case in any way
		if err := d.Set("comments", resp.Comments.Value); err != nil {
			return diag.Errorf("error setting comments: %s", err)
		}
	}

	if resp.Compression != nil {
		if err := d.Set("compression", strings.ToLower(*resp.Compression.Value)); err != nil {
			return diag.Errorf("error setting compression: %s", err)
		}
	}

	if resp.Deduplication != nil {
		if err := d.Set("deduplication", strings.ToLower(*resp.Deduplication.Value)); err != nil {
			return diag.Errorf("error setting deduplication: %s", err)
		}
	}

	if resp.Exec != nil {
		if err := d.Set("exec", strings.ToLower(*resp.Exec.Value)); err != nil {
			return diag.Errorf("error setting exec: %s", err)
		}
	}

	if resp.KeyFormat != nil && resp.KeyFormat.Value != nil {
		if err := d.Set("key_format", strings.ToLower(*resp.KeyFormat.Value)); err != nil {
			return diag.Errorf("error setting key_format: %s", err)
		}
	}

	if resp.ManagedBy != nil {
		if err := d.Set("managed_by", resp.ManagedBy.Value); err != nil {
			return diag.Errorf("error setting managed_by: %s", err)
		}
	}

	if resp.Copies != nil {
		copies, err := strconv.Atoi(*resp.Copies.Value)

		if err != nil {
			return diag.Errorf("error parsing copies: %s", err)
		}

		if err := d.Set("copies", copies); err != nil {
			return diag.Errorf("error setting copies: %s", err)
		}
	}

	if resp.Quota != nil {
		quota, err := strconv.Atoi(resp.Quota.RawValue)

		if err != nil {
			return diag.Errorf("error parsing quota: %s", err)
		}

		if err := d.Set("quota_bytes", quota); err != nil {
			return diag.Errorf("error setting quota_bytes: %s", err)
		}
	}

	if resp.QuotaCritical != nil {
		quota, err := strconv.Atoi(*resp.QuotaCritical.Value)

		if err != nil {
			return diag.Errorf("error parsing quota_critical: %s", err)
		}

		if err := d.Set("quota_critical", quota); err != nil {
			return diag.Errorf("error setting quota_critical: %s", err)
		}
	}

	if resp.QuotaWarning != nil {
		quota, err := strconv.Atoi(*resp.QuotaWarning.Value)

		if err != nil {
			return diag.Errorf("error parsing quota_warning: %s", err)
		}

		if err := d.Set("quota_warning", quota); err != nil {
			return diag.Errorf("error setting quota_warning: %s", err)
		}
	}

	if resp.Reservation != nil {
		resrv, err := strconv.Atoi(resp.Reservation.RawValue)

		if err != nil {
			return diag.Errorf("error parsing reservation: %s", err)
		}

		if err := d.Set("reservation", resrv); err != nil {
			return diag.Errorf("error setting reservation: %s", err)
		}

	}

	if resp.RefReservation != nil {
		resrv, err := strconv.Atoi(resp.RefReservation.RawValue)

		if err != nil {
			return diag.Errorf("error parsing refreservation: %s", err)
		}

		if err := d.Set("ref_reservation", resrv); err != nil {
			return diag.Errorf("error setting ref_reservation: %s", err)
		}

	}

	if resp.RefQuota != nil {
		quota, err := strconv.Atoi(resp.RefQuota.RawValue)

		if err != nil {
			return diag.Errorf("error parsing refquota: %s", err)
		}

		if err := d.Set("ref_quota_bytes", quota); err != nil {
			return diag.Errorf("error setting ref_quota_bytes: %s", err)
		}
	}

	if resp.RefQuotaCritical != nil {
		quota, err := strconv.Atoi(*resp.RefQuotaCritical.Value)

		if err != nil {
			return diag.Errorf("error parsing refquota_critical: %s", err)
		}

		if err := d.Set("ref_quota_critical", quota); err != nil {
			return diag.Errorf("error setting ref_quota_critical: %s", err)
		}
	}

	if resp.RefQuotaWarning != nil {
		quota, err := strconv.Atoi(*resp.RefQuotaWarning.Value)

		if err != nil {
			return diag.Errorf("error parsing refquota_warning: %s", err)
		}

		if err := d.Set("ref_quota_warning", quota); err != nil {
			return diag.Errorf("error setting ref_quota_warning: %s", err)
		}
	}

	if resp.Readonly != nil {
		if err := d.Set("readonly", strings.ToLower(*resp.Readonly.Value)); err != nil {
			return diag.Errorf("error setting readonly: %s", err)
		}
	}

	if resp.Recordsize != nil {
		if err := d.Set("record_size", *resp.Recordsize.Value); err != nil {
			return diag.Errorf("error setting record_size: %s", err)
		}

		sz, err := strconv.Atoi(resp.Recordsize.RawValue)

		if err != nil {
			return diag.Errorf("error parsing recordsize rawvalue: %s", err)
		}

		if err := d.Set("record_size_bytes", sz); err != nil {
			return diag.Errorf("error setting record_size_bytes: %s", err)
		}
	}

	if resp.Sync != nil {
		if err := d.Set("sync", strings.ToLower(*resp.Sync.Value)); err != nil {
			return diag.Errorf("error setting sync: %s", err)
		}
	}

	if resp.SnapDir != nil {
		if err := d.Set("snap_dir", strings.ToLower(*resp.SnapDir.Value)); err != nil {
			return diag.Errorf("error setting snap_dir: %s", err)
		}
	}

	if resp.EncryptionAlgorithm != nil && resp.EncryptionAlgorithm.Value != nil {
		if err := d.Set("encryption_algorithm", *resp.EncryptionAlgorithm.Value); err != nil {
			return diag.Errorf("error setting encryption_algorithm: %s", err)
		}
	}

	if resp.PBKDF2Iters != nil {
		iters, err := strconv.Atoi(*resp.PBKDF2Iters.Value)

		if err != nil {
			return diag.Errorf("error parsing PBKDF2Iters: %s", err)
		}

		if iters >= 0 {
			if err := d.Set("pbkdf2iters", iters); err != nil {
				return diag.Errorf("error setting PBKDF2Iters: %s", err)
			}
		}
	}

	if resp.Origin != nil {
		if err := d.Set("origin", strings.ToLower(*resp.Origin.Value)); err != nil {
			return diag.Errorf("error setting origin: %s", err)
		}
	}

	if resp.XATTR != nil {
		if err := d.Set("xattr", strings.ToLower(*resp.XATTR.Value)); err != nil {
			return diag.Errorf("error setting xattr: %s", err)
		}
	}

	if err := d.Set("encrypted", resp.Encrypted); err != nil {
		return diag.Errorf("error setting encrypted: %s", err)
	}

	d.SetId(resp.ID)

	return diags
}
