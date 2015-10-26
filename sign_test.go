package getstream_test

import (
	. "github.com/onedotover/go-getstream"
	a "github.com/stretchr/testify/assert"
	"testing"
)

func TestSign(t *testing.T) {
	result := Sign(MockAPISecret, TestFeedSlug.Slug+TestFeedSlug.ID)
	a.Equal(t, TestToken, result)
}

func TestSignSlug(t *testing.T) {
	expected := TestFeedSlug
	expected.Token = TestToken

	actual := TestFeedSlug
	actual.Sign(MockAPISecret)

	a.Equal(t, expected, actual)
	a.Equal(t, TestFeedSignature, actual.Signature())
}

func TestSignActivity(t *testing.T) {
	act := TestActivityTarget
	act.Sign(MockAPISecret)
	a.Equal(t, TestTargetFeedSignature, act.To[0].Signature())
	a.Equal(t, TestTargetFeedSignature2, act.To[1].Signature())
}
