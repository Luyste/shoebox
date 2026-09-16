package media

import (
	"testing"
)

func TestKindFromMimeType(t *testing.T) {
	tests :=
		[]struct {
			name     string
			mimeType string
			want     string
			wantErr  bool
		}{
			{name: "jpeg maps to photo", mimeType: "image/jpeg", want: "photo", wantErr: false},
			{name: "mp4 maps to video", mimeType: "video/mp4", want: "video", wantErr: false},
			{name: "json maps to error", mimeType: "application/json", want: "", wantErr: true},
		}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := kindFromMimeType(tc.mimeType)

			if got != tc.want {
				t.Errorf("want=%v, got=%v", tc.want, got)
			}

			if tc.wantErr != (err != nil) {
				t.Errorf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
		})
	}
}

func TestInspect(t *testing.T) {
	tests := []struct {
		name         string
		content      []byte
		want         string
		wantMimeType string
		wantErr      bool
	}{
		{name: "jpeg signature returns photo kind and checksum matches", content: []byte{0xFF,
			0xD8, 0xFF, 0xFD, 0xD3}, want: "photo", wantMimeType: "image/jpeg", wantErr: false},
		{name: "mp4 signature returns video kind and checksum matches", content: []byte{
			0x00, 0x00, 0x00, 0x10, 'f', 't', 'y', 'p', 'm', 'p', '4', '2', 0x00, 0x00, 0x00, 0x00}, want: "video", wantMimeType: "video/mp4", wantErr: false},
		{name: "random signature returns error and random checksum", content: []byte("some random string"), want: "", wantMimeType: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			checksum := buildChecksum(tc.content)
			size := len(tc.content)

			got, err := Inspect(tc.content)

			if err == nil {
				if got.Checksum != checksum {
					t.Errorf("want checksum: %v, got checksum: %v", checksum, got.Checksum)
				}

				if got.Size != size {
					t.Errorf("want size=%v, got size=%v", size, got.Size)
				}
			}

			if got.Kind != tc.want {
				t.Errorf("want kind=%v, got kind=%v", tc.want, got.Kind)
			}

			if got.MimeType != tc.wantMimeType {
				t.Errorf("want kind=%v, got kind=%v", tc.wantMimeType, got.MimeType)
			}

			if tc.wantErr != (err != nil) {
				t.Errorf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
		})
	}
}
