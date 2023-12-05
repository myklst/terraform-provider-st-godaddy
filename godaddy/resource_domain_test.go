package godaddy_provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeAndCalculateMode(t *testing.T) {

	now := time.Now().UTC()
	const layout string = "2006-01-02T15:04:05.000Z"
	fmt.Println(now.Format(layout))
	var output string
	var err diag.Diagnostic

	const oneDayRemaining int64 = 1
	const thirtyDaysRemaining int64 = 30
	const negativeDaysRemaining int64 = -1
	const moreThanOneYearRemaining int64 = 367

	domainExpiresOneYearFromNow := now.AddDate(1, 0, 0)
	domainExpiry := domainExpiresOneYearFromNow.Format(layout)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(oneDayRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(thirtyDaysRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(negativeDaysRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(moreThanOneYearRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "renew", output)

	domainExpiresThirtyOneDaysFromNow := now.AddDate(0, 0, 31)
	domainExpiry = domainExpiresThirtyOneDaysFromNow.Format(layout)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(oneDayRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(thirtyDaysRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(negativeDaysRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(moreThanOneYearRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "renew", output)

	domainExpiresOneDayFromNow := now.AddDate(0, 0, 1)
	domainExpiry = domainExpiresOneDayFromNow.Format(layout)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(oneDayRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "renew", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(thirtyDaysRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "renew", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(negativeDaysRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "skip", output)

	output, err = ParseTimeAndCalculateMode(domainExpiry, int64(moreThanOneYearRemaining))
	assert.Equal(t, err, nil)
	assert.Equal(t, "renew", output)
}
