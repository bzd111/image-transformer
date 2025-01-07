package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNeedChangeImage(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name            string
		image           string
		newRepo         string
		originalRepo    string
		expectedImage   string
		expectedChanged bool
	}{
		{
			name:            "simple image",
			image:           "busybox:latest",
			newRepo:         "aa.bb.cc",
			originalRepo:    "asdf",
			expectedImage:   "aa.bb.cc/busybox:latest",
			expectedChanged: true,
		},
		{
			name:            "image with single slash",
			image:           "docker.io/busybox:latest",
			newRepo:         "aa.bb.cc",
			originalRepo:    "docker.io",
			expectedImage:   "aa.bb.cc/docker.io/busybox:latest",
			expectedChanged: true,
		},
		{
			name:            "image with multiple slashes",
			image:           "docker.io/library/busybox:latest",
			newRepo:         "aa.bb.cc",
			originalRepo:    "docker.io",
			expectedImage:   "aa.bb.cc/docker.io/library/busybox:latest",
			expectedChanged: true,
		},
		{
			name:            "image without change",
			image:           "aa.bb.cc/library/busybox:latest",
			newRepo:         "aa.bb.cc",
			originalRepo:    "docker.io",
			expectedImage:   "",
			expectedChanged: false,
		},
		{
			name:            "image with different original repo",
			image:           "gcr.io/library/busybox:latest",
			newRepo:         "aa.bb.cc",
			originalRepo:    "gcr.io",
			expectedImage:   "aa.bb.cc/gcr.io/library/busybox:latest",
			expectedChanged: true,
		},
		{
			name:            "image with repo and name ",
			image:           "groundnuty/k8s-wait-for:v2.0",
			newRepo:         "aa.bb.cc",
			originalRepo:    "",
			expectedImage:   "aa.bb.cc/groundnuty/k8s-wait-for:v2.0",
			expectedChanged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(t.Name())
			newImage, changed := needChangeImage(tt.image, tt.newRepo, tt.originalRepo)
			assert.Equal(tt.expectedImage, newImage)
			assert.Equal(tt.expectedChanged, changed)
		})
	}
}
