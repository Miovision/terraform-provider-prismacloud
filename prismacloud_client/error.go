package prismacloud_client

type PrismaCloudError struct {
	Reason   string `json:"i18nKey"`
	Severity string `json:"severity"`
	Subject  string `json:"subject"`
}
