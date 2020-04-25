// +build integration

package facest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFaces(t *testing.T) {
	c := NewClient(os.Getenv("FACEST_INTEGRATION_API_KEY"))

	res, err := c.GetFaces(nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")

	assert.Equal(t, 1, res.Count, "expecting 1 face found")
	assert.Equal(t, 1, res.PagesCount, "expecting 1 page found")

	assert.Equal(t, "integration_face_id", res.Faces[0].FaceID, "expecting correct face_id")
	assert.NotEmpty(t, res.Faces[0].FaceToken, "expecting non-empty face_token")
	assert.Greater(t, len(res.Faces[0].FaceImages), 0, "expecting non-empty face_images")
}
