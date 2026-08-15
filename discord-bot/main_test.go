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

func TestStartupProbeServeMux(t *testing.T) {
	testCases := []struct {
		desc         string
		isStarted    bool
		target       string
		expectedCode int
	}{
		{
			desc:         "not started",
			isStarted:    false,
			target:       "/",
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			desc:         "started",
			isStarted:    true,
			target:       "/",
			expectedCode: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a := App{
				isStarted: &atomic.Bool{},
			}
			a.isStarted.Store(tC.isStarted)

			mux := a.newStartupProbeServeMux(":8081")
			req := httptest.NewRequest(http.MethodGet, tC.target, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tC.expectedCode, rec.Code)
		})
	}
}
