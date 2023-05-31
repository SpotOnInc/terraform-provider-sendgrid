# sendgrid_sso_teammate

Provide a resource to manage SSO teammates.

## Example Usage

```hcl
resource "sendgrid_sso_teammate" "teammate" {
	email       = "jane.doe@example.com"
	first_name  = "Jane"
	last_name   = "Doe"
	is_admin    = false
	persona     = "observer"
	scopes      = ["mail.send", "alerts.read"]
}
```

## Argument Reference

The following arguments are supported:

* `scopes` - (Optional) The scopes for the SendGrid TeamMate access.


## Import

An SSO teammate can be imported, e.g.
```sh
$ terraform import sendgrid_sso_teammate.teammate <email>
```
