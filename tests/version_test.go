package tests

import (
	"testing"

	"github.com/lordofscripts/goapp/app"
)

const (
	tstNAME string = "testapp"
	tstDESC string = "Just testing"
)

// Alpha versions
func TestPackageVersion_Alpha(t *testing.T) {
	const BASE_VERSION string = "1.0.5"
	const REVISION int = 3
	pver := app.NewAlphaVersion(tstNAME, tstDESC, BASE_VERSION, REVISION)

	const EXPECT string = "v1.0.5-Alpha.3"
	got := pver.Short()
	if got != EXPECT {
		t.Errorf("alpha version got '%s' exp '%s'", got, EXPECT)
	}

	pver = app.NewPackageVersion(tstNAME, tstDESC, BASE_VERSION, app.DevStatusAlpha)
	pver.WithRevision(REVISION)
	got = pver.Short()
	if got != EXPECT {
		t.Errorf("alpha version got '%s' exp '%s'", got, EXPECT)
	}
}

// Beta versions
func TestPackageVersion_Beta(t *testing.T) {
	const BASE_VERSION string = "1.0.5"
	const REVISION int = 2
	pver := app.NewBetaVersion(tstNAME, tstDESC, BASE_VERSION, REVISION)

	const EXPECT string = "v1.0.5-Beta.2"
	got := pver.Short()
	if got != EXPECT {
		t.Errorf("alpha version got '%s' exp '%s'", got, EXPECT)
	}

	pver = app.NewPackageVersion(tstNAME, tstDESC, BASE_VERSION, app.DevStatusBeta)
	pver.WithRevision(REVISION)
	got = pver.Short()
	if got != EXPECT {
		t.Errorf("beta version got '%s' exp '%s'", got, EXPECT)
	}
}

// Release Candidate
func TestPackageVersion_RC(t *testing.T) {
	const BASE_VERSION string = "1.0.3"
	const REVISION int = 8
	pver := app.NewReleaseCandidateVersion(tstNAME, tstDESC, BASE_VERSION, REVISION)

	const EXPECT string = "v1.0.3-RC.8"
	got := pver.Short()
	if got != EXPECT {
		t.Errorf("alpha version got '%s' exp '%s'", got, EXPECT)
	}

	pver = app.NewPackageVersion(tstNAME, tstDESC, BASE_VERSION, app.DevStatusRC)
	pver.WithRevision(REVISION)
	got = pver.Short()
	if got != EXPECT {
		t.Errorf("release candidate version got '%s' exp '%s'", got, EXPECT)
	}
}

// Final Release
func TestPackageVersion_Final(t *testing.T) {
	const BASE_VERSION string = "1.0.3"
	pver := app.NewReleaseVersion(tstNAME, tstDESC, BASE_VERSION)

	const EXPECT string = "v1.0.3"
	got := pver.Short()
	if got != EXPECT {
		t.Errorf("final release version got '%s' exp '%s'", got, EXPECT)
	}

	pver = app.NewPackageVersion(tstNAME, tstDESC, BASE_VERSION, app.DevStatusReleased)
	got = pver.Short()
	if got != EXPECT {
		t.Errorf("final release version got '%s' exp '%s'", got, EXPECT)
	}
}
