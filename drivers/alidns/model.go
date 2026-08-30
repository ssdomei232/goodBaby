package alidns

type AliDNSAccount struct {
	AK string `json:"ak"`
	SK string `json:"sk"`
}

type DeleteRecordConfig struct {
	RecordID string `json:"record_id"`
}
