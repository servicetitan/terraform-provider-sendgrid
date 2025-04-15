/*
Provide a resource to manage an event webhook settings.
Example Usage
```hcl

	resource "sendgrid_event_webhook" "default" {
		enabled = true
	    url = "https://foo.bar/sendgrid/inbound"
	    group_resubscribe = true
	    delivered = true
	    group_unsubscribe = true
	    spam_report = true
	    bounce = true
	    deferred = true
	    unsubscribe = true
	    processed = true
	    open = true
	    click = true
	    dropped = true
	    oauth_client_id = "a-client-id"
	    oauth_client_secret = "a-client-secret"
	    oauth_token_url = "https://oauth.example.com/token"
	}

```
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/sapronov-st/terraform-provider-sendgrid/sdk"
)

func resourceSendgridEventWebhook() *schema.Resource { //nolint:funlen
	return &schema.Resource{
		CreateContext: resourceSendgridEventWebhookCreate,
		ReadContext:   resourceSendgridEventWebhookRead,
		UpdateContext: resourceSendgridEventWebhookUpdate,
		DeleteContext: resourceSendgridEventWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indicates if the event webhook is enabled.",
				Required:    true,
			},
			"url": {
				Type: schema.TypeString,
				Description: "The public URL where you would like SendGrid to POST the data events from your email. " +
					"Any emails sent with the given hostname provided (whose MX records have been updated to point to SendGrid) " +
					"will be eventd and POSTed to this URL.",
				Required: true,
			},
			"group_resubscribe": {
				Type: schema.TypeBool,
				Description: "Recipient resubscribes to specific group by updating preferences. " +
					"You need to enable Subscription Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"delivered": {
				Type:        schema.TypeBool,
				Description: "Message has been successfully delivered to the receiving server.",
				Optional:    true,
				Default:     true,
			},
			"group_unsubscribe": {
				Type: schema.TypeBool,
				Description: "Recipient unsubscribe from specific group, by either direct link or updating preferences. " +
					"You need to enable Subscription Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"spam_report": {
				Type:        schema.TypeBool,
				Description: "Recipient marked a message as spam.",
				Optional:    true,
				Default:     true,
			},
			"bounce": {
				Type:        schema.TypeBool,
				Description: "Receiving server could not or would not accept message.",
				Optional:    true,
				Default:     true,
			},
			"deferred": {
				Type:        schema.TypeBool,
				Description: "Recipient's email server temporarily rejected message.",
				Optional:    true,
				Default:     true,
			},
			"unsubscribe": {
				Type: schema.TypeBool,
				Description: "Recipient clicked on message's subscription management link. " +
					"You need to enable Subscription Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"processed": {
				Type:        schema.TypeBool,
				Description: "Message has been received and is ready to be delivered.",
				Optional:    true,
				Default:     true,
			},
			"open": {
				Type: schema.TypeBool,
				Description: "Recipient has opened the HTML message. " +
					"You need to enable Open Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"click": {
				Type: schema.TypeBool,
				Description: "Recipient clicked on a link within the message. " +
					"You need to enable Click Tracking for getting this type of event.",
				Optional: true,
				Default:  true,
			},
			"dropped": {
				Type: schema.TypeBool,
				Description: "You may see the following drop reasons: " +
					"Invalid SMTPAPI header, Spam Content (if spam checker app enabled), " +
					"Unsubscribed Address, Bounced Address, Spam Reporting Address, Invalid, Recipient List over Package Quota.",
				Optional: true,
				Default:  true,
			},
			"oauth_client_id": {
				Type: schema.TypeString,
				Description: "The client ID Twilio SendGrid sends to your OAuth server or " +
					"service provider to generate an OAuth access token.",
				Optional: true,
			},
			"oauth_client_secret": {
				Type: schema.TypeString,
				Description: "This secret is needed only once to create an access token. SendGrid will store this secret, " +
					"allowing you to update your Client ID and Token URL without passing the secret to SendGrid again. " +
					"When passing data in this field, you must also include the oauth_client_id and oauth_token_url fields.",
				Optional:  true,
				Sensitive: true,
			},
			"oauth_token_url": {
				Type: schema.TypeString,
				Description: "The URL where Twilio SendGrid sends the Client ID and Client Secret to generate an access token. " +
					"This should be your OAuth server or service provider. " +
					"When passing data in this field, you must also include the oauth_client_id field.",
				Optional: true,
			},
			"signed": {
				Type:        schema.TypeBool,
				Description: "Should the event webhook use signing?",
				Optional:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public key used to sign the event webhook. Only present if 'signed' is true",
				Computed:    true,
			},
		},
	}
}

func resourceSendgridEventWebhookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// no op as there is no way to delete it
	return nil
}

func resourceSendgridEventWebhookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	webhookResource := readEventWebhookFromResource(d)

	result, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateEventWebhook(webhookResource)
	})
	webhook := result.(*sendgrid.EventWebhook)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(webhook.Id)

	signed := d.Get("signed").(bool)
	if signed {
		if _, err := c.ConfigureEventWebhookSigning(webhook.Id, signed); err.Err != nil {
			return diag.FromErr(err.Err)
		}
	}

	return resourceSendgridEventWebhookRead(ctx, d, c)
}

func resourceSendgridEventWebhookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	webhookResouce := readEventWebhookFromResource(d)

	if d.HasChangesExcept("signed", "public_key") {
		_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
			return c.PatchEventWebhook(webhookResouce)
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("signed") {
		signed := d.Get("signed").(bool)
		if _, err := c.ConfigureEventWebhookSigning(webhookResouce.Id, signed); err.Err != nil {
			return diag.FromErr(err.Err)
		}
	}

	return resourceSendgridEventWebhookRead(ctx, d, c)
}

func resourceSendgridEventWebhookRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	webhook, err := c.ReadEventWebhook(d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	webhookSigning, err := c.ReadEventWebhookSigning(d.Id())
	if err.Err != nil {
		return diag.FromErr(err.Err)
	}

	setEventWebhookToResource(d, webhook, webhookSigning)

	return nil
}

func readEventWebhookFromResource(d *schema.ResourceData) *sendgrid.EventWebhook {
	return &sendgrid.EventWebhook{
		Id:                d.Id(),
		Enabled:           d.Get("enabled").(bool),
		URL:               d.Get("url").(string),
		GroupResubscribe:  d.Get("group_resubscribe").(bool),
		Delivered:         d.Get("delivered").(bool),
		GroupUnsubscribe:  d.Get("group_unsubscribe").(bool),
		SpamReport:        d.Get("spam_report").(bool),
		Bounce:            d.Get("bounce").(bool),
		Deferred:          d.Get("deferred").(bool),
		Unsubscribe:       d.Get("unsubscribe").(bool),
		Processed:         d.Get("processed").(bool),
		Open:              d.Get("open").(bool),
		Click:             d.Get("click").(bool),
		Dropped:           d.Get("dropped").(bool),
		OAuthClientID:     d.Get("oauth_client_id").(string),
		OAuthClientSecret: d.Get("oauth_client_secret").(string),
		OAuthTokenURL:     d.Get("oauth_token_url").(string),
	}
}

func setEventWebhookToResource(d *schema.ResourceData, webhook *sendgrid.EventWebhook, webhookSigning *sendgrid.EventWebhookSigning) {
	d.SetId(webhook.Id)
	//nolint:errcheck
	d.Set("enabled", webhook.Enabled)
	//nolint:errcheck
	d.Set("url", webhook.URL)
	//nolint:errcheck
	d.Set("group_resubscribe", webhook.GroupResubscribe)
	//nolint:errcheck
	d.Set("delivered", webhook.Delivered)
	//nolint:errcheck
	d.Set("group_unsubscribe", webhook.GroupUnsubscribe)
	//nolint:errcheck
	d.Set("spam_report", webhook.SpamReport)
	//nolint:errcheck
	d.Set("bounce", webhook.Bounce)
	//nolint:errcheck
	d.Set("deferred", webhook.Delivered)
	//nolint:errcheck
	d.Set("unsubscribe", webhook.Unsubscribe)
	//nolint:errcheck
	d.Set("processed", webhook.Processed)
	//nolint:errcheck
	d.Set("open", webhook.Open)
	//nolint:errcheck
	d.Set("click", webhook.Click)
	//nolint:errcheck
	d.Set("dropped", webhook.Dropped)
	//nolint:errcheck
	d.Set("oauth_client_id", webhook.OAuthClientID)
	//nolint:errcheck
	d.Set("oauth_token_url", webhook.OAuthTokenURL)

	//nolint:errcheck
	d.Set("public_key", webhookSigning.PublicKey)
	//nolint:errcheck
	d.Set("signed", webhookSigning.PublicKey != "")
}
