package scriptrunner

var envScripts = map[string]string{
	"export":             "export-account.sh",
	"import":             "ofco-restore-site-from-backup.sh",
	"addAccount":         "create-new-account.sh",
	"cloneAccount":       "clone-wordpress-account.sh",
	"terminateAccount":   "ofco-delete-account.sh",
	"addDomainToAccount": "add-domain-to-account.sh",
	"unban-ip":           "unban-ip.sh",
}
