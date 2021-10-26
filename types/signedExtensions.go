package types

type SignedExtension struct {
	Name             string             `json:"name"`
	AdditionalSigned []AdditionalSigned `json:"additional_signed"`
}

type AdditionalSigned struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
