package protocol

import "testing"

func TestParseMOTDStringInput(t *testing.T) {
	input := "A Minecraft Server"
	got := parseMOTD(input)
	if got != input {
		t.Fatalf("parseMOTD() = %q, want %q", got, input)
	}
}

func TestParseMOTDStructuredMapWithExtra(t *testing.T) {
	input := map[string]interface{}{
		"text": "Base ",
		"extra": []interface{}{
			map[string]interface{}{"text": "One "},
			"Two ",
			map[string]interface{}{
				"text": "Nested ",
				"extra": []interface{}{
					map[string]interface{}{"text": "Depth"},
				},
			},
		},
	}

	want := "Base One Two Nested Depth"
	got := parseMOTD(input)
	if got != want {
		t.Fatalf("parseMOTD() = %q, want %q", got, want)
	}
}

func TestParseMOTDUnsupportedTypeReturnsEmpty(t *testing.T) {
	got := parseMOTD(12345)
	if got != "" {
		t.Fatalf("parseMOTD() = %q, want empty string", got)
	}
}
