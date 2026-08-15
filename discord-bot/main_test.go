package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*func Benchmark(b *testing.B) {
	envVars, err := loadEnvVars()
	if err != nil {
		b.Errorf("failed to load env vars: %v", err)
	}
	tok := envVars.discordToken
	client, err := disgo.New(tok,
		bot.WithDefaultGateway(),
	)
	if err != nil {
		b.Errorf("error while building disgo: %v", err)
	}
	defer client.Close(context.TODO())

	for b.Loop() {
	}
}*/

type testCase struct {
	desc         string
	isStarted    bool
	isReady      bool
	isLiving     bool
	target       string
	expectedCode int
}

func newTestCase(desc, target string, expectedCode int) testCase {
	return testCase{
		desc:         desc,
		isStarted:    true,
		isReady:      true,
		isLiving:     true,
		target:       target,
		expectedCode: expectedCode,
	}
}

func (tc testCase) WithNotStarted() testCase {
	tc.isStarted = false
	return tc
}

func (tc testCase) WithNotReady() testCase {
	tc.isReady = false
	return tc
}

func (tc testCase) WithNotLiving() testCase {
	tc.isLiving = false
	return tc
}

func TestACAProbeServeMux(t *testing.T) {
	testCases := []testCase{
		newTestCase("not started", "/startup", http.StatusServiceUnavailable).WithNotStarted(),
		newTestCase("started", "/startup", http.StatusOK),
		newTestCase("not ready", "/readiness", http.StatusServiceUnavailable).WithNotReady(),
		newTestCase("ready", "/readiness", http.StatusOK),
		newTestCase("not living", "/liveness", http.StatusServiceUnavailable).WithNotLiving(),
		newTestCase("living", "/liveness", http.StatusOK),
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a := App{
				isStarted: &atomic.Bool{},
				isReady:   &atomic.Bool{},
				isLiving:  &atomic.Bool{},
			}
			a.isStarted.Store(tC.isStarted)
			a.isReady.Store(tC.isReady)
			a.isLiving.Store(tC.isLiving)

			mux := a.newACAProbeServeMux()
			req := httptest.NewRequest(http.MethodGet, tC.target, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tC.expectedCode, rec.Code)
		})
	}
}
