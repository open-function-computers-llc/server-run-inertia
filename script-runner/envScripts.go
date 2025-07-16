package scriptrunner

var envScripts = map[string]envScript{
	"export": {
		filename:     "export-account.sh",
		requiredEnvs: []string{"ACCOUNT"},
	},
	"import": {
		filename:     "ofco-restore-site-from-backup.sh",
		requiredEnvs: []string{"BACKUP_TARBALL", "DESTINATION_ACCOUNT"},
	},
	"addAccount": {
		filename:     "create-new-account.sh",
		requiredEnvs: []string{"ACCOUNT_NAME"},
		optionalEnvs: []string{"DOMAIN"},
	},
	"cloneAccount": {
		filename:     "clone-wordpress-account.sh",
		requiredEnvs: []string{"SOURCE_ACCOUNT", "DESTINATION_ACCOUNT", "DESTINATION_DOMAIN"},
	},
	"terminateAccount": {
		filename:     "ofco-delete-account.sh",
		requiredEnvs: []string{"ACCOUNT_TO_DELETE"},
	},
	"addDomainToAccount": {
		filename:     "add-domain-to-account.sh",
		requiredEnvs: []string{"ACCOUNT", "DOMAIN"},
	},
	"unban-ip": {
		filename:     "unban-ip.sh",
		requiredEnvs: []string{"IPADDRESS"},
	},
	"ban-ip": {
		filename:     "ban-ip.sh",
		requiredEnvs: []string{"IPADDRESS"},
	},
	"update-wordpress-domain": {
		filename:     "change-wordpress-domain.sh",
		requiredEnvs: []string{"DOMAIN", "ACCOUNT"},
	},
}

type envScript struct {
	filename     string
	requiredEnvs []string
	optionalEnvs []string
}
