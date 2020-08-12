package uptime_test

import (
	"net/http"
	"time"

	"code.cloudfoundry.org/cf-k8s-networking/test/uptime/internal/uptime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data Plane Uptime", func() {
	var (
		results *uptime.DataPlaneResults
		// upgradeChecker *checker.Upgrade
		startTime time.Time
	)

	BeforeEach(func() {
		results = &uptime.DataPlaneResults{}
		// upgradeChecker = &checker.Upgrade{
		// 	PollInterval: 1 * time.Second,
		// }
		// upgradeChecker.Start()

		startTime = time.Now()
	})

	AfterEach(func() {
		// upgradeChecker.Stop()
	})

	FIt("measures the data plane uptime", func() {
		By("checking whether X% of requests are successful within the acceptable response time during an upgrade", func() {
			for {
				// if !upgradeChecker.HasFoundUpgrade() && time.Since(startTime) > upgradeDiscoveryTimeout {
				// 	Fail(fmt.Sprintf("failed to find cf upgrade in %s", upgradeDiscoveryTimeout.String()))
				// }

				// if the upgrade is finished (learned by checking the "finished at" in
				// kapp app-change ls), stop running the test
				if time.Since(startTime) > (time.Minute * 1) {
					break
				}

				resp, err, requestLatency := timeGetRequest(dataPlaneSLIAppRouteURL)
				if err != nil {
					results.RecordError(err)
					continue
				}

				hasStatusOK := resp.StatusCode == http.StatusOK
				hasMetRequestLatencySLO := requestLatency < dataPlaneSLOMaxRequestLatency
				hasPassedSLI := hasStatusOK && hasMetRequestLatencySLO

				results.Record(hasPassedSLI,
					hasStatusOK,
					hasMetRequestLatencySLO,
					requestLatency)

			}

			results.PrintResults()

			Expect(results.SuccessPercentage()).To(BeNumerically(">=", dataPlaneSLOPercentage))
		})
	})
})
