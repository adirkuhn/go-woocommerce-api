package woocommerce

import (
	"bytes"
	"encoding/json"
)

// FlexibleString unmarshals from either a JSON string or a bare JSON number,
// coercing both to a Go string. WooCommerce documents price fields
// (price/regular_price/sale_price) as always strings, but some stores return
// them unquoted instead — observed with a bundle-product plugin ("woosb"
// type) that recalculates sale_price and re-serialises it as a number. The
// raw literal bytes are used as-is for the numeric case, so no float
// precision or formatting is lost. Always marshals back out as a normal
// JSON string, since that's what WooCommerce's write endpoints expect.
type FlexibleString string

func (s *FlexibleString) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || string(trimmed) == "null" {
		*s = ""
		return nil
	}
	if trimmed[0] == '"' {
		var str string
		if err := json.Unmarshal(trimmed, &str); err != nil {
			return err
		}
		*s = FlexibleString(str)
		return nil
	}
	*s = FlexibleString(trimmed)
	return nil
}

type MetaData struct {
	ID         int    `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      any    `json:"value,omitempty"`
	DisplayKey string `json:"display_key"`
}

type Self struct {
	Href string `json:"href,omitempty"`
}

type Collection struct {
	Href string `json:"href,omitempty"`
}

type Links struct {
	Self       []Self       `json:"self,omitempty"`
	Collection []Collection `json:"collection,omitempty"`
}

type Billing struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Company   string `json:"company,omitempty"`
	Address1  string `json:"address_1,omitempty"`
	Address2  string `json:"address_2,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	Postcode  string `json:"postcode,omitempty"`
	Country   string `json:"country,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type Shipping struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Company   string `json:"company,omitempty"`
	Address1  string `json:"address_1,omitempty"`
	Address2  string `json:"address_2,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	Postcode  string `json:"postcode,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
}
