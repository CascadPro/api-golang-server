package core_http_request

import "testing"

func TestReplaceExtension(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		newExtension string
		want         string
	}{
		{name: "replace existing extension", filename: "avatar.jpeg", newExtension: ".webp", want: "avatar.webp"},
		{name: "no extension", filename: "avatar", newExtension: ".webp", want: "avatar.webp"},
		{name: "multiple dots", filename: "profile.photo.jpeg", newExtension: ".webp", want: "profile.photo.webp"},
		{name: "dotfile", filename: ".avatar", newExtension: ".webp", want: ".avatar.webp"},
		{name: "trailing dot", filename: "avatar.", newExtension: ".webp", want: "avatar.webp"},
		{name: "empty filename", filename: "", newExtension: ".webp", want: "file.webp"},
		{name: "extension without leading dot", filename: "avatar.jpeg", newExtension: "webp", want: "avatar.webp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceExtension(tt.filename, tt.newExtension); got != tt.want {
				t.Fatalf("replaceExtension(%q, %q) = %q, want %q", tt.filename, tt.newExtension, got, tt.want)
			}
		})
	}
}
