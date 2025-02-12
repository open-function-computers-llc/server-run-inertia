package scriptrunner

var envScripts = map[string]envScript{
	"export": envScript{
		filename:     "export-account.sh",
		requiredEnvs: []string{"ACCOUNT"},
	},
	"import": envScript{
		filename:     "ofco-restore-site-from-backup.sh",
		requiredEnvs: []string{"BACKUP_TARBALL", "DESTINATION_ACCOUNT"},
	},
	"addAccount": envScript{
		filename:     "create-new-account.sh",
		requiredEnvs: []string{"ACCOUNT_NAME"},
		optionalEnvs: []string{"DOMAIN"},
	},
	"cloneAccount": envScript{
		filename:     "clone-wordpress-account.sh",
		requiredEnvs: []string{"SOURCE_ACCOUNT", "DESTINATION_ACCOUNT", "DESTINATION_DOMAIN"},
	},
	"terminateAccount": envScript{
		filename:     "ofco-delete-account.sh",
		requiredEnvs: []string{"ACCOUNT_TO_DELETE"},
	},
	"addDomainToAccount": envScript{
		filename:     "add-domain-to-account.sh",
		requiredEnvs: []string{"ACCOUNT", "DOMAIN"},
	},
	"unban-ip": envScript{
		filename:     "unban-ip.sh",
		requiredEnvs: []string{"IPADDRESS"},
	},
	"update-wordpress-domain": envScript{
		filename:     "change-wordpress-domain.sh",
		requiredEnvs: []string{"DOMAIN", "ACCOUNT"},
	},
}

type envScript struct {
	filename     string
	requiredEnvs []string
	optionalEnvs []string
}
