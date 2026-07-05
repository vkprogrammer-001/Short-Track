package base62

import "testing"

func TestEncodeDecode(t *testing.T) {
	tests := []int64{
		0,
		1,
		62,
		63,
		999,
		123456789,
		9223372036854775807, // Max int64
	}

	for _, original := range tests {
		encoded := Encode(original)
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf("Decode failed for encoded value of %d (%s): %v", original, encoded, err)
			continue
		}
		if decoded != original {
			t.Errorf("Mismatch for %d: got %d (encoded: %s)", original, decoded, encoded)
		}
	}
}

func TestDecodeInvalid(t *testing.T) {
	_, err := Decode("invalid-character!")
	if err == nil {
		t.Error("Expected error for invalid base62 character, but got nil")
	}
}
