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
	"add-account": {
		filename:     "create-new-account.sh",
		requiredEnvs: []string{"ACCOUNT_NAME"},
		optionalEnvs: []string{"DOMAIN", "INSTALL_WP", "INSTALL_RAD", "ADMIN_PASSWORD", "SITE_TITLE", "ADMIN_USERNAME", "ADMIN_EMAIL"},
	},
	"clone-account": {
		filename:     "clone-wordpress-account.sh",
		requiredEnvs: []string{"SOURCE_ACCOUNT", "DESTINATION_ACCOUNT"},
		optionalEnvs: []string{"DESTINATION_DOMAIN"},
	},
	"terminate-account": {
		filename:     "ofco-delete-account.sh",
		requiredEnvs: []string{"ACCOUNT_TO_DELETE"},
	},
	"add-domain-to-account": {
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
	// read-only lookup of why an IP is banned, without unbanning it
	"ban-lookup": {
		filename:     "ban-lookup.sh",
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
