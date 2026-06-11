package cmdb

import (
	"strings"
	"testing"
)

func TestParseCSVAssets(t *testing.T) {
	input := "unique_key,hostname,management_ip\nsn:ABC123,web-01,10.0.1.10\n"
	rows, err := ParseCSVAssets(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].UniqueKey != "sn:ABC123" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}
