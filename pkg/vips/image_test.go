package vips_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/davidbyttow/govips/pkg/vips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadImage_AccessMode(t *testing.T) {
	srcBytes, err := ioutil.ReadFile("testdata/test.png")
	require.NoError(t, err)

	// defaults to random access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src)
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// random access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src, vips.WithAccessMode(vips.AccessRandom))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// sequential access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src, vips.WithAccessMode(vips.AccessSequential))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check sequential access by encoding twice where the second fails
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of order")
		}
	}
}

func TestNewImageFromFile_AccessMode(t *testing.T) {
	// defaults to random access
	{
		img, err := vips.NewImageFromFile("testdata/test.png")
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// random access
	{
		img, err := vips.NewImageFromFile("testdata/test.png", vips.WithAccessMode(vips.AccessRandom))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// sequential access
	{
		img, err := vips.NewImageFromFile("testdata/test.png", vips.WithAccessMode(vips.AccessSequential))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check sequential access by encoding twice where the second fails
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of order")
		}
	}
}

func TestNewImageFromBuffer_AccessMode(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/test.png")
	require.NoError(t, err)

	// defaults to random access
	{
		img, err := vips.NewImageFromBuffer(src)
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// random access
	{
		img, err := vips.NewImageFromBuffer(src, vips.WithAccessMode(vips.AccessRandom))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// sequential access
	{
		img, err := vips.NewImageFromBuffer(src, vips.WithAccessMode(vips.AccessSequential))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check sequential access by encoding twice where the second fails
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of order")
		}
	}
}

func TestImageRef_HasProfile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			"image with profile",
			"testdata/with_icc_profile.jpg",
			true,
		},
		{
			"image without profile",
			"testdata/without_icc_profile.jpg",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := vips.NewImageFromFile(tt.path)
			require.NoError(t, err)
			got := ref.HasProfile()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestImageRef_AutorotAngle(t *testing.T) {
	tests := []struct {
		name string
		path string
		want vips.Angle
	}{
		{
			"image without EXIF orientation",
			"testdata/without_exif.jpg",
			vips.Angle0,
		},
		{
			"image with EXIF orientation top-left",
			"testdata/with_exif_orientation_top_left.jpg",
			vips.Angle0,
		},
		{
			"image with EXIF orientation right-top",
			"testdata/with_exif_orientation_right_top.jpg",
			vips.Angle90,
		},
		{
			"image with EXIF orientation bottom-right",
			"testdata/with_exif_orientation_bottom_right.jpg",
			vips.Angle180,
		},
		{
			"image with EXIF orientation left-bottom",
			"testdata/with_exif_orientation_left_bottom.jpg",
			vips.Angle270,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := vips.NewImageFromFile(tt.path)
			require.NoError(t, err)
			got := ref.AutorotAngle()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestImageRef_Export(t *testing.T) {
	img, err := vips.NewImageFromFile("testdata/test.png")
	require.NoError(t, err)
	tests := map[string]struct {
		params   vips.ExportParams
		wantFile string
		wantType vips.ImageType
	}{
		"WebP lossy with Quality=25": {
			params:   vips.ExportParams{Format: vips.ImageTypeWEBP, Quality: 25},
			wantFile: "testdata/test_q25.webp",
			wantType: vips.ImageTypeWEBP,
		},
		"WebP lossy with Quality=75": {
			params:   vips.ExportParams{Format: vips.ImageTypeWEBP, Quality: 75},
			wantFile: "testdata/test_q75.webp",
			wantType: vips.ImageTypeWEBP,
		},
		"WebP lossless with ReductionEffort=0": {
			params:   vips.ExportParams{Format: vips.ImageTypeWEBP, Lossless: true, ReductionEffort: 0},
			wantFile: "testdata/test_lossless_re0.webp",
			wantType: vips.ImageTypeWEBP,
		},
		"WebP lossless with defaults (ReductionEffort=4)": {
			params:   vips.ExportParams{Format: vips.ImageTypeWEBP, Lossless: true, ReductionEffort: 4},
			wantFile: "testdata/test_lossless_re4.webp",
			wantType: vips.ImageTypeWEBP,
		},
		"WebP lossless with ReductionEffort=6": {
			params:   vips.ExportParams{Format: vips.ImageTypeWEBP, Lossless: true, ReductionEffort: 6},
			wantFile: "testdata/test_lossless_re6.webp",
			wantType: vips.ImageTypeWEBP,
		},
	}
	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			buf, typ, err := img.Export(tt.params)
			require.NoError(t, err)
			goldenBuf, err := ioutil.ReadFile(tt.wantFile)
			assert.Equal(t, goldenBuf, buf)
			assert.Equal(t, tt.wantType, typ)
		})
	}
}
