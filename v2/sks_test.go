package v2

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testSKSClusterAddons                     = []string{"exoscale-cloud-controller"}
	testSKSClusterCNI                        = "calico"
	testSKSClusterCreatedAt, _               = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
	testSKSClusterDescription                = "Test Cluster description"
	testSKSClusterEndpoint                   = fmt.Sprintf("%s.sks-%s.exo.io", testSKSClusterID, testZone)
	testSKSClusterID                         = new(clientTestSuite).randomID()
	testSKSClusterName                       = "test-cluster"
	testSKSClusterServiceLevel               = "pro"
	testSKSClusterState                      = "running"
	testSKSClusterVersion                    = "1.18.6"
	testSKSNodepoolAntiAffinityGroupID       = new(clientTestSuite).randomID()
	testSKSNodepoolCreatedAt, _              = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
	testSKSNodepoolDescription               = "Test Nodepool description"
	testSKSNodepoolDiskSize            int64 = 15
	testSKSNodepoolID                        = new(clientTestSuite).randomID()
	testSKSNodepoolInstancePoolID            = new(clientTestSuite).randomID()
	testSKSNodepoolInstanceTypeID            = new(clientTestSuite).randomID()
	testSKSNodepoolName                      = "test-nodepool"
	testSKSNodepoolSecurityGroupID           = new(clientTestSuite).randomID()
	testSKSNodepoolSize                int64 = 3
	testSKSNodepoolState                     = "running"
	testSKSNodepoolTemplateID                = new(clientTestSuite).randomID()
	testSKSNodepoolVersion                   = "1.18.6"
)

func (ts *clientTestSuite) TestSKSCluster_RotateCCMCredentials() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,
	}

	ts.mockAPIRequest("PUT", fmt.Sprintf("/sks-cluster/%s/rotate-ccm-credentials", cluster.ID),
		papi.Operation{
			Id:        &testOperationID,
			State:     &testOperationState,
			Reference: &papi.Reference{Id: &testSKSNodepoolID},
		})

	ts.Require().NoError(cluster.RotateCCMCredentials(context.Background()))
}

func (ts *clientTestSuite) TestSKSCluster_AuthorityCert() {
	var (
		testAuthority   = "aggregation"
		testCertificate = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,
	}

	ts.mockAPIRequest("GET",
		fmt.Sprintf("/sks-cluster/%s/authority/%s/cert", cluster.ID, testAuthority),
		struct {
			Cacert string `json:"cacert,omitempty"`
		}{
			Cacert: testCertificate,
		})

	actual, err := cluster.AuthorityCert(context.Background(), testAuthority)
	ts.Require().NoError(err)
	ts.Require().Equal(testCertificate, actual)
}

func (ts *clientTestSuite) TestSKSCluster_RequestKubeconfig() {
	var (
		testRequestUser   = "test-user"
		testRequestGroups = []string{"system:masters"}
		testKubeconfig    = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,
	}

	ts.mockAPIRequest("POST", fmt.Sprintf("/sks-cluster-kubeconfig/%s", cluster.ID), struct {
		Kubeconfig string `json:"kubeconfig,omitempty"`
	}{
		Kubeconfig: testKubeconfig,
	})

	actual, err := cluster.RequestKubeconfig(context.Background(), testRequestUser, testRequestGroups, time.Hour)
	ts.Require().NoError(err)
	ts.Require().Equal(testKubeconfig, actual)
}

func (ts *clientTestSuite) TestSKSCluster_AddNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	ts.mockAPIRequest("POST", fmt.Sprintf("/sks-cluster/%s/nodepool", testSKSClusterID),
		papi.Operation{
			Id:        &testOperationID,
			State:     &testOperationState,
			Reference: &papi.Reference{Id: &testSKSNodepoolID},
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		testSKSClusterID, testSKSNodepoolID),
		papi.SksNodepool{
			AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
			CreatedAt:          &testSKSNodepoolCreatedAt,
			Description:        &testSKSNodepoolDescription,
			DiskSize:           &testSKSNodepoolDiskSize,
			Id:                 &testSKSNodepoolID,
			InstancePool:       &papi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
			InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
			Name:               &testSKSNodepoolName,
			SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
			Size:               &testSKSNodepoolSize,
			State:              &testSKSNodepoolState,
			Template:           &papi.Template{Id: &testSKSNodepoolTemplateID},
			Version:            &testSKSNodepoolVersion,
		})

	cluster := &SKSCluster{
		ID: testSKSClusterID,

		c:    ts.client,
		zone: testZone,
	}

	expected := &SKSNodepool{
		AntiAffinityGroupIDs: []string{testSKSNodepoolAntiAffinityGroupID},
		CreatedAt:            testSKSNodepoolCreatedAt,
		Description:          testSKSNodepoolDescription,
		DiskSize:             testSKSNodepoolDiskSize,
		ID:                   testSKSNodepoolID,
		InstancePoolID:       testSKSNodepoolInstancePoolID,
		InstanceTypeID:       testSKSNodepoolInstanceTypeID,
		Name:                 testSKSNodepoolName,
		SecurityGroupIDs:     []string{testSKSNodepoolSecurityGroupID},
		Size:                 testSKSNodepoolSize,
		State:                testSKSNodepoolState,
		TemplateID:           testSKSNodepoolTemplateID,
		Version:              testSKSNodepoolVersion,
	}

	actual, err := cluster.AddNodepool(context.Background(), expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestSKSCluster_UpdateNodepool() {
	var (
		testSKSNodepoolNameUpdated                = testSKSNodepoolName + "-updated"
		testSKSNodepoolDescriptionUpdated         = testSKSNodepoolDescription + "-updated"
		testSKSNodepoolInstanceTypeIDUpdated      = testSKSNodepoolInstanceTypeID + "-updated"
		testSKSNodepoolDiskSizeUpdated            = testSKSNodepoolDiskSize + 1
		testSKSNodepoolAntiAffinityGroupIDUpdated = ts.randomID()
		testSKSNodepoolSecurityGroupIDUpdated     = ts.randomID()
		testOperationID                           = ts.randomID()
		testOperationState                        = "success"
	)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{
			{
				Description:    testSKSNodepoolDescription,
				DiskSize:       testSKSNodepoolDiskSize,
				ID:             testSKSNodepoolID,
				InstanceTypeID: testSKSNodepoolInstanceTypeID,
				Name:           testSKSNodepoolName,
			},
		},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		cluster.ID,
		cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpdateSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateSksNodepoolJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupIDUpdated}},
				Description:        &testSKSNodepoolDescriptionUpdated,
				DiskSize:           &testSKSNodepoolDiskSizeUpdated,
				InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeIDUpdated},
				Name:               &testSKSNodepoolNameUpdated,
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupIDUpdated}},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	nodepoolUpdated := SKSNodepool{
		AntiAffinityGroupIDs: []string{testSKSNodepoolAntiAffinityGroupIDUpdated},
		Description:          testSKSNodepoolDescriptionUpdated,
		DiskSize:             testSKSNodepoolDiskSizeUpdated,
		ID:                   cluster.Nodepools[0].ID,
		InstanceTypeID:       testSKSNodepoolInstanceTypeIDUpdated,
		Name:                 testSKSNodepoolNameUpdated,
		SecurityGroupIDs:     []string{testSKSNodepoolSecurityGroupIDUpdated},
	}
	ts.Require().NoError(cluster.UpdateNodepool(context.Background(), &nodepoolUpdated))
}

func (ts *clientTestSuite) TestSKSCluster_ScaleNodepool() {
	var (
		testOperationID          = ts.randomID()
		testOperationState       = "success"
		testScaleSize      int64 = 3
	)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: testSKSNodepoolID}},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s:scale",
		cluster.ID,
		cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.ScaleSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)
			expected := papi.ScaleSksNodepoolJSONRequestBody{Size: testScaleSize}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(cluster.ScaleNodepool(context.Background(), cluster.Nodepools[0], testScaleSize))
}

func (ts *clientTestSuite) TestSKSCluster_EvictNodepoolMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = "success"
		testEvictedMemberID = ts.randomID()
	)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: testSKSNodepoolID}},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s:evict",
		cluster.ID,
		cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.EvictSksNodepoolMembersJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.EvictSksNodepoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(cluster.EvictNodepoolMembers(
		context.Background(),
		cluster.Nodepools[0],
		[]string{testEvictedMemberID}))
}

func (ts *clientTestSuite) TestSKSCluster_DeleteNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("=~^/sks-cluster/%s/nodepool/.*", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
				testSKSClusterID, testSKSNodepoolID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: testSKSNodepoolID}},
	}

	ts.Require().NoError(cluster.DeleteNodepool(context.Background(), cluster.Nodepools[0]))
}

func (ts *clientTestSuite) TestSKSCluster_ResetField() {
	var (
		testResetField     = "description"
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/sks-cluster/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(
				fmt.Sprintf("/sks-cluster/%s/%s", testSKSClusterID, testResetField),
				req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: testSKSNodepoolID}},
	}

	ts.Require().NoError(cluster.ResetField(context.Background(), &cluster.Description))
}

func (ts *clientTestSuite) TestSKSCluster_ResetNodepoolField() {
	var (
		testResetField     = "description"
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/sks-cluster/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(
				fmt.Sprintf("/sks-cluster/%s/nodepool/%s/%s",
					testSKSClusterID, testSKSNodepoolID, testResetField),
				req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: testSKSNodepoolID}},
	}

	ts.Require().NoError(cluster.ResetNodepoolField(
		context.Background(),
		cluster.Nodepools[0],
		&cluster.Nodepools[0].Description))
}

func (ts *clientTestSuite) TestClient_CreateSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("POST", "/sks-cluster",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateSksClusterJSONRequestBody{
				Addons:      &testSKSClusterAddons,
				Cni:         &testSKSClusterCNI,
				Description: &testSKSClusterDescription,
				Level:       testSKSClusterServiceLevel,
				Name:        testSKSClusterName,
				Version:     testSKSClusterVersion,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), papi.SksCluster{
		Addons:      &testSKSClusterAddons,
		Cni:         &testSKSClusterCNI,
		CreatedAt:   &testSKSClusterCreatedAt,
		Description: &testSKSClusterDescription,
		Id:          &testSKSClusterID,
		Level:       &testSKSClusterServiceLevel,
		Name:        &testSKSClusterName,
		State:       &testSKSClusterState,
		Version:     &testSKSClusterVersion,
	})

	expected := &SKSCluster{
		AddOns:       testSKSClusterAddons,
		CNI:          testSKSClusterCNI,
		CreatedAt:    testSKSClusterCreatedAt,
		Description:  testSKSClusterDescription,
		ID:           testSKSClusterID,
		Name:         testSKSClusterName,
		Nodepools:    []*SKSNodepool{},
		ServiceLevel: testSKSClusterServiceLevel,
		State:        testSKSClusterState,
		Version:      testSKSClusterVersion,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.CreateSKSCluster(context.Background(), testZone, &SKSCluster{
		AddOns:       testSKSClusterAddons,
		CNI:          testSKSClusterCNI,
		Description:  testSKSClusterDescription,
		Name:         testSKSClusterName,
		ServiceLevel: testSKSClusterServiceLevel,
		Version:      testSKSClusterVersion,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListSKSClusters() {
	ts.mockAPIRequest("GET", "/sks-cluster", struct {
		SksClusters *[]papi.SksCluster `json:"sks-clusters,omitempty"`
	}{
		SksClusters: &[]papi.SksCluster{{
			Addons:      &testSKSClusterAddons,
			Cni:         &testSKSClusterCNI,
			CreatedAt:   &testSKSClusterCreatedAt,
			Description: &testSKSClusterDescription,
			Endpoint:    &testSKSClusterEndpoint,
			Id:          &testSKSClusterID,
			Level:       &testSKSClusterServiceLevel,
			Name:        &testSKSClusterName,
			Nodepools: &[]papi.SksNodepool{{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
				CreatedAt:          &testSKSNodepoolCreatedAt,
				Description:        &testSKSNodepoolDescription,
				DiskSize:           &testSKSNodepoolDiskSize,
				Id:                 &testSKSNodepoolID,
				InstancePool:       &papi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
				InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
				Name:               &testSKSNodepoolName,
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
				Size:               &testSKSNodepoolSize,
				State:              &testSKSNodepoolState,
				Template:           &papi.Template{Id: &testSKSNodepoolTemplateID},
				Version:            &testSKSNodepoolVersion,
			}},
			State:   &testSKSClusterState,
			Version: &testSKSClusterVersion,
		}},
	})

	expected := []*SKSCluster{{
		AddOns:      testSKSClusterAddons,
		CNI:         testSKSClusterCNI,
		CreatedAt:   testSKSClusterCreatedAt,
		Description: testSKSClusterDescription,
		Endpoint:    testSKSClusterEndpoint,
		ID:          testSKSClusterID,
		Name:        testSKSClusterName,
		Nodepools: []*SKSNodepool{{
			AntiAffinityGroupIDs: []string{testSKSNodepoolAntiAffinityGroupID},
			CreatedAt:            testSKSNodepoolCreatedAt,
			Description:          testSKSNodepoolDescription,
			DiskSize:             testSKSNodepoolDiskSize,
			ID:                   testSKSNodepoolID,
			InstancePoolID:       testSKSNodepoolInstancePoolID,
			InstanceTypeID:       testSKSNodepoolInstanceTypeID,
			Name:                 testSKSNodepoolName,
			SecurityGroupIDs:     []string{testSKSNodepoolSecurityGroupID},
			Size:                 testSKSNodepoolSize,
			State:                testSKSClusterState,
			TemplateID:           testSKSNodepoolTemplateID,
			Version:              testSKSNodepoolVersion,
		}},
		ServiceLevel: testSKSClusterServiceLevel,
		State:        testSKSClusterState,
		Version:      testSKSClusterVersion,

		c:    ts.client,
		zone: testZone,
	}}

	actual, err := ts.client.ListSKSClusters(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListSKSClusterVersions() {
	var (
		testSKSClusterVersions = []string{
			"1.20.0",
			"1.18.6",
		}
		err error
	)

	ts.mockAPIRequest("GET", "/sks-cluster-version", struct {
		SksClusterVersions *[]string `json:"sks-cluster-versions,omitempty"`
	}{
		SksClusterVersions: &testSKSClusterVersions,
	})

	expected := testSKSClusterVersions
	actual, err := ts.client.ListSKSClusterVersions(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetSKSCluster() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), papi.SksCluster{
		Addons:      &testSKSClusterAddons,
		Cni:         &testSKSClusterCNI,
		CreatedAt:   &testSKSClusterCreatedAt,
		Description: &testSKSClusterDescription,
		Endpoint:    &testSKSClusterEndpoint,
		Id:          &testSKSClusterID,
		Level:       &testSKSClusterServiceLevel,
		Name:        &testSKSClusterName,
		Nodepools: &[]papi.SksNodepool{{
			AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
			CreatedAt:          &testSKSNodepoolCreatedAt,
			Description:        &testSKSNodepoolDescription,
			DiskSize:           &testSKSNodepoolDiskSize,
			Id:                 &testSKSNodepoolID,
			InstancePool:       &papi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
			InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
			Name:               &testSKSNodepoolName,
			SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
			Size:               &testSKSNodepoolSize,
			State:              &testSKSNodepoolState,
			Template:           &papi.Template{Id: &testSKSNodepoolTemplateID},
			Version:            &testSKSNodepoolVersion,
		}},
		State:   &testSKSClusterState,
		Version: &testSKSClusterVersion,
	})

	expected := &SKSCluster{
		AddOns:      testSKSClusterAddons,
		CNI:         testSKSClusterCNI,
		CreatedAt:   testSKSClusterCreatedAt,
		Description: testSKSClusterDescription,
		Endpoint:    testSKSClusterEndpoint,
		ID:          testSKSClusterID,
		Name:        testSKSClusterName,
		Nodepools: []*SKSNodepool{{
			AntiAffinityGroupIDs: []string{testSKSNodepoolAntiAffinityGroupID},
			CreatedAt:            testSKSNodepoolCreatedAt,
			Description:          testSKSNodepoolDescription,
			DiskSize:             testSKSNodepoolDiskSize,
			ID:                   testSKSNodepoolID,
			InstancePoolID:       testSKSNodepoolInstancePoolID,
			InstanceTypeID:       testSKSNodepoolInstanceTypeID,
			Name:                 testSKSNodepoolName,
			SecurityGroupIDs:     []string{testSKSNodepoolSecurityGroupID},
			Size:                 testSKSNodepoolSize,
			State:                testSKSClusterState,
			TemplateID:           testSKSNodepoolTemplateID,
			Version:              testSKSNodepoolVersion,
		}},
		ServiceLevel: testSKSClusterServiceLevel,
		State:        testSKSClusterState,
		Version:      testSKSClusterVersion,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetSKSCluster(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateSKSCluster() {
	var (
		testSKSClusterNameUpdated        = testSKSClusterName + "-updated"
		testSKSClusterDescriptionUpdated = testSKSClusterDescription + "-updated"
		testOperationID                  = ts.randomID()
		testOperationState               = "success"
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpdateSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateSksClusterJSONRequestBody{
				Name:        &testSKSClusterNameUpdated,
				Description: &testSKSClusterDescriptionUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	clusterUpdated := SKSCluster{
		ID:          testSKSClusterID,
		Name:        testSKSClusterNameUpdated,
		Description: testSKSClusterDescriptionUpdated,
	}
	ts.Require().NoError(ts.client.UpdateSKSCluster(context.Background(), testZone, &clusterUpdated))
}

func (ts *clientTestSuite) TestClient_UgradeSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/sks-cluster/%s/upgrade", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpgradeSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpgradeSksClusterJSONRequestBody{Version: testSKSClusterVersion}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(ts.client.UpgradeSKSCluster(
		context.Background(),
		testZone,
		testSKSClusterID,
		testSKSClusterVersion))
}

func (ts *clientTestSuite) TestClient_DeleteSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/sks-cluster/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	ts.Require().NoError(ts.client.DeleteSKSCluster(context.Background(), testZone, testSKSClusterID))
}
