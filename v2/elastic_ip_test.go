package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testElasticIPDescription                    = "Test Elastic IP description"
	testElasticIPID                             = new(clientTestSuite).randomID()
	testElasticIPAddress                        = "1.2.3.4"
	testElasticIPHealthcheckMode                = "https"
	testElasticIPHealthcheckPort          int64 = 8080
	testElasticIPHealthcheckInterval      int64 = 10
	testElasticIPHealthcheckTimeout       int64 = 3
	testElasticIPHealthcheckStrikesFail   int64 = 1
	testElasticIPHealthcheckStrikesOK     int64 = 1
	testElasticIPHealthcheckURI                 = "/health"
	testElasticIPHealthcheckTLSSNI              = "example.net"
	testElasticIPHealthcheckTLSSkipVerify       = true
)

func (ts *clientTestSuite) TestElasticIP_ResetField() {
	var (
		testResetField     = "description"
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/elastic-ip/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(
				fmt.Sprintf("/elastic-ip/%s/%s", testInstancePoolID, testResetField),
				req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	elasticIP := &ElasticIP{
		ID:   testInstancePoolID,
		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(elasticIP.ResetField(context.Background(), &elasticIP.Description))
}

func (ts *clientTestSuite) TestClient_CreateElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("POST", "/elastic-ip",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateElasticIpJSONRequestBody{
				Description: &testElasticIPDescription,
				Healthcheck: &papi.ElasticIpHealthcheck{
					Interval:      testElasticIPHealthcheckInterval,
					Mode:          testNLServiceHealthcheckMode,
					Port:          testElasticIPHealthcheckPort,
					StrikesFail:   testElasticIPHealthcheckStrikesFail,
					StrikesOk:     testElasticIPHealthcheckStrikesOK,
					Timeout:       testElasticIPHealthcheckTimeout,
					TlsSni:        &testElasticIPHealthcheckTLSSNI,
					TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
					Uri:           &testElasticIPHealthcheckURI,
				},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testElasticIPID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testElasticIPID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/elastic-ip/%s", testElasticIPID), papi.ElasticIp{
		Description: &testElasticIPDescription,
		Healthcheck: &papi.ElasticIpHealthcheck{
			Interval:      testElasticIPHealthcheckInterval,
			Mode:          testNLServiceHealthcheckMode,
			Port:          testElasticIPHealthcheckPort,
			StrikesFail:   testElasticIPHealthcheckStrikesFail,
			StrikesOk:     testElasticIPHealthcheckStrikesOK,
			Timeout:       testElasticIPHealthcheckTimeout,
			TlsSni:        &testElasticIPHealthcheckTLSSNI,
			TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Uri:           &testElasticIPHealthcheckURI,
		},
		Id: &testElasticIPID,
		Ip: &testElasticIPAddress,
	})

	expected := &ElasticIP{
		Description: testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      time.Duration(testElasticIPHealthcheckInterval) * time.Second,
			Mode:          testElasticIPHealthcheckMode,
			Port:          uint16(testElasticIPHealthcheckPort),
			StrikesFail:   testElasticIPHealthcheckStrikesFail,
			StrikesOK:     testElasticIPHealthcheckStrikesOK,
			TLSSNI:        testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       time.Duration(testElasticIPHealthcheckTimeout) * time.Second,
			URI:           testElasticIPHealthcheckURI,
		},
		ID:        testElasticIPID,
		IPAddress: net.ParseIP(testElasticIPAddress),

		zone: testZone,
		c:    ts.client,
	}

	actual, err := ts.client.CreateElasticIP(context.Background(), testZone, &ElasticIP{
		Description: testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      time.Duration(testElasticIPHealthcheckInterval) * time.Second,
			Mode:          testElasticIPHealthcheckMode,
			Port:          uint16(testElasticIPHealthcheckPort),
			StrikesFail:   testElasticIPHealthcheckStrikesFail,
			StrikesOK:     testElasticIPHealthcheckStrikesOK,
			TLSSNI:        testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       time.Duration(testElasticIPHealthcheckTimeout) * time.Second,
			URI:           testElasticIPHealthcheckURI,
		},
		ID: testElasticIPID,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListElasticIPs() {
	ts.mockAPIRequest("GET", "/elastic-ip", struct {
		ElasticIPs *[]papi.ElasticIp `json:"elastic-ips,omitempty"`
	}{
		ElasticIPs: &[]papi.ElasticIp{{
			Description: &testElasticIPDescription,
			Healthcheck: &papi.ElasticIpHealthcheck{
				Interval:      testElasticIPHealthcheckInterval,
				Mode:          testNLServiceHealthcheckMode,
				Port:          testElasticIPHealthcheckPort,
				StrikesFail:   testElasticIPHealthcheckStrikesFail,
				StrikesOk:     testElasticIPHealthcheckStrikesOK,
				Timeout:       testElasticIPHealthcheckTimeout,
				TlsSni:        &testElasticIPHealthcheckTLSSNI,
				TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
				Uri:           &testElasticIPHealthcheckURI,
			},
			Id: &testElasticIPID,
			Ip: &testElasticIPAddress,
		}},
	})

	expected := []*ElasticIP{{
		Description: testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      time.Duration(testElasticIPHealthcheckInterval) * time.Second,
			Mode:          testElasticIPHealthcheckMode,
			Port:          uint16(testElasticIPHealthcheckPort),
			StrikesFail:   testElasticIPHealthcheckStrikesFail,
			StrikesOK:     testElasticIPHealthcheckStrikesOK,
			TLSSNI:        testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       time.Duration(testElasticIPHealthcheckTimeout) * time.Second,
			URI:           testElasticIPHealthcheckURI,
		},
		ID:        testElasticIPID,
		IPAddress: net.ParseIP(testElasticIPAddress),

		zone: testZone,
		c:    ts.client,
	}}

	actual, err := ts.client.ListElasticIPs(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetElasticIP() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/elastic-ip/%s", testElasticIPID), papi.ElasticIp{
		Description: &testElasticIPDescription,
		Healthcheck: &papi.ElasticIpHealthcheck{
			Interval:      testElasticIPHealthcheckInterval,
			Mode:          testNLServiceHealthcheckMode,
			Port:          testElasticIPHealthcheckPort,
			StrikesFail:   testElasticIPHealthcheckStrikesFail,
			StrikesOk:     testElasticIPHealthcheckStrikesOK,
			Timeout:       testElasticIPHealthcheckTimeout,
			TlsSni:        &testElasticIPHealthcheckTLSSNI,
			TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Uri:           &testElasticIPHealthcheckURI,
		},
		Id: &testElasticIPID,
		Ip: &testElasticIPAddress,
	})

	expected := &ElasticIP{
		Description: testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      time.Duration(testElasticIPHealthcheckInterval) * time.Second,
			Mode:          testElasticIPHealthcheckMode,
			Port:          uint16(testElasticIPHealthcheckPort),
			StrikesFail:   testElasticIPHealthcheckStrikesFail,
			StrikesOK:     testElasticIPHealthcheckStrikesOK,
			TLSSNI:        testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       time.Duration(testElasticIPHealthcheckTimeout) * time.Second,
			URI:           testElasticIPHealthcheckURI,
		},
		ID:        testElasticIPID,
		IPAddress: net.ParseIP(testElasticIPAddress),

		zone: testZone,
		c:    ts.client,
	}

	actual, err := ts.client.GetElasticIP(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateElasticIP() {
	var (
		testElasticIPDescriptionUpdated              = testElasticIPDescription + "-updated"
		testElasticIPHealthcheckModeUpdated          = "http"
		testElasticIPHealthcheckPortUpdated          = testElasticIPHealthcheckPort + 1
		testElasticIPHealthcheckIntervalUpdated      = testElasticIPHealthcheckInterval + 1
		testElasticIPHealthcheckTimeoutUpdated       = testElasticIPHealthcheckTimeout + 1
		testElasticIPHealthcheckStrikesFailUpdated   = testElasticIPHealthcheckStrikesFail + 1
		testElasticIPHealthcheckStrikesOKUpdated     = testElasticIPHealthcheckStrikesOK + 1
		testElasticIPHealthcheckURIUpdated           = ""
		testElasticIPHealthcheckTLSSNIUpdated        = ""
		testElasticIPHealthcheckTLSSkipVerifyUpdated = false
		testOperationID                              = ts.randomID()
		testOperationState                           = "success"
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/elastic-ip/%s", testElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpdateElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateElasticIpJSONRequestBody{
				Description: &testElasticIPDescriptionUpdated,
				Healthcheck: &papi.ElasticIpHealthcheck{
					Interval:      testElasticIPHealthcheckIntervalUpdated,
					Mode:          testElasticIPHealthcheckModeUpdated,
					Port:          testElasticIPHealthcheckPortUpdated,
					StrikesFail:   testElasticIPHealthcheckStrikesFailUpdated,
					StrikesOk:     testElasticIPHealthcheckStrikesOKUpdated,
					Timeout:       testElasticIPHealthcheckTimeoutUpdated,
					TlsSni:        &testElasticIPHealthcheckTLSSNIUpdated,
					TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerifyUpdated,
					Uri:           &testElasticIPHealthcheckURIUpdated,
				},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testElasticIPID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testElasticIPID},
	})

	ts.Require().NoError(ts.client.UpdateElasticIP(context.Background(), testZone, &ElasticIP{
		Description: testElasticIPDescriptionUpdated,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      time.Duration(testElasticIPHealthcheckIntervalUpdated) * time.Second,
			Mode:          testElasticIPHealthcheckModeUpdated,
			Port:          uint16(testElasticIPHealthcheckPortUpdated),
			StrikesFail:   testElasticIPHealthcheckStrikesFailUpdated,
			StrikesOK:     testElasticIPHealthcheckStrikesOKUpdated,
			Timeout:       time.Duration(testElasticIPHealthcheckTimeoutUpdated) * time.Second,
			TLSSNI:        testElasticIPHealthcheckTLSSNIUpdated,
			TLSSkipVerify: testElasticIPHealthcheckTLSSkipVerifyUpdated,
			URI:           testElasticIPHealthcheckURIUpdated,
		},
		ID:        testElasticIPID,
		IPAddress: net.ParseIP(testElasticIPAddress),

		zone: testZone,
		c:    ts.client,
	}))
}

func (ts *clientTestSuite) TestClient_DeleteElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/elastic-ip/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/elastic-ip/%s", testElasticIPID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testElasticIPID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testElasticIPID},
	})

	ts.Require().NoError(ts.client.DeleteElasticIP(context.Background(), testZone, testElasticIPID))
}
