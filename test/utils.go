package test

import (
	"context"
	"github.com/lgunko/beauty-organisation-service/cmd/start"
	"github.com/lgunko/beauty-reuse/env"
	"github.com/lgunko/beauty-reuse/headers"
	"github.com/lgunko/beauty-reuse/server"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type HookTestSuite struct {
	suite.Suite
	hook     *test.Hook
	database mongo.Database
}

type LogLevelMessage struct {
	Level   logrus.Level
	Message string
}

func (suite *HookTestSuite) SetupSuite() {
	suite.hook = test.NewGlobal()
	logrus.AddHook(suite.hook)

	os.Setenv("EXECUTION_ENVIRONMENT", env.Test.String())
	env.SetUpEnv()
	suite.database = server.GetDatabase()

	go func() {
		start.Start()
	}()

	time.Sleep(1 * time.Second)
}

const testingOrgID = "testingOrgID"
const testingEmail = "leonidgunko1@yandex.ru"

func (suite *HookTestSuite) testWithoutHeader(result interface{}, gqlErrors []*gqlerror.Error, err error) {
	assert.Nil(suite.T(), result)
	assert.Nil(suite.T(), gqlErrors)
	assert.Error(suite.T(), err)
	assert.Errorf(suite.T(), err, "EOF")

	assert.True(suite.T(), len(suite.hook.AllEntries()) > 0)

	logLevelMessageList := []LogLevelMessage{}
	for _, entry := range suite.hook.AllEntries() {
		logLevelMessageList = append(logLevelMessageList, LogLevelMessage{entry.Level, entry.Message})
	}
	for _, header := range headers.AllRequiredHeaders {
		assert.Contains(suite.T(), logLevelMessageList, LogLevelMessage{logrus.ErrorLevel, "No Value " + header.String() + " in context"})
	}
	for _, header := range headers.AllTracingHeaders {
		assert.Contains(suite.T(), logLevelMessageList, LogLevelMessage{logrus.ErrorLevel, "No Value " + header.String() + " in context"})
	}
}

func setHeaders(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, headers.Authorization, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjRMOUx0ZmlWTEFsS3Z4QjFDaFhJbyJ9.eyJnaXZlbl9uYW1lIjoiTGVvbmlkIiwiZmFtaWx5X25hbWUiOiJHdW5rbyIsIm5pY2tuYW1lIjoibGVvbmlkZ3Vua28xIiwibmFtZSI6Ikxlb25pZCBHdW5rbyIsInBpY3R1cmUiOiJodHRwczovL3Njb250ZW50LWxocjgtMS54eC5mYmNkbi5uZXQvdi90MS4zMDQ5Ny0xL2NwMC9jMTUuMC41MC41MGEvcDUweDUwLzg0NjI4MjczXzE3NjE1OTgzMDI3Nzg1Nl85NzI2OTMzNjM5MjI4MjkzMTJfbi5qcGc_X25jX2NhdD0xJmNjYj0xLTMmX25jX3NpZD0xMmIzYmUmX25jX29oYz1kcHNGQjluTWZPZ0FYLWtjdGtDJl9uY19odD1zY29udGVudC1saHI4LTEueHgmdHA9MjcmcmVzbG9nPWQmb2g9NGM2NTNkOGU4Y2JhNzY1N2RiY2Y0NjMzODkwYmE1MzEmb2U9NjA4NjQxMzgiLCJ1cGRhdGVkX2F0IjoiMjAyMS0wMy0yN1QxNTo0MDo0Mi40NTZaIiwiZW1haWwiOiJsZW9uaWRndW5rbzFAeWFuZGV4LnJ1IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vcHJpbWUuZXUuYXV0aDAuY29tLyIsInN1YiI6ImZhY2Vib29rfDEwMzAwOTQ0MzczMzM1ODgiLCJhdWQiOiJHQ2x0dW5EbjNid3VXSXRjZFlpRVNaMm1ZZkdWWHZBUyIsImlhdCI6MTYxNjg3NDU1OSwiZXhwIjoxNjE2OTEwNTU5LCJub25jZSI6ImNGZDFXbVZsWlRKeGFsVkxPV1JZWTNsdFUxOHhkelk0UWtrdVpYVmhkMWhYVWkxblowTjRPV0pYVWc9PSJ9.bGE6t5Fvfb0nq12DuJ9aETlzl9-b-X0RSb-STOwNXolqQwgg7q1-iHXTu_Zp7WhL1wMEmqcnaj45fmLcE6BSN6QGsnCoHIeAGR4tavS-uc2iLG6-rgq3CNfUmgEvG-MZ4IlYGwb9rls0E3SXkbsCZHnqV6dIpVBF472pjpjbkTf8qDUaaow5OedOpn6PX-of0wpGdudTcs4aVv73c1AGU-nPEBryTGS7-7IXyoYg_IXDSQDUaAciMz4k_OKDX8a7a3FM2jr8QkIuawmo7hI6RNA2cEPu10UTe3D-zHyQbbtTKchK8m0vJzYJtlLwZVRe-5fKD7Ttnf9dxl7lbsdrxg")
	ctx = context.WithValue(ctx, headers.Email, testingEmail)
	ctx = context.WithValue(ctx, headers.OrgID, testingOrgID)
	ctx = context.WithValue(ctx, headers.Role, "none")
	ctx = context.WithValue(ctx, headers.GqlOperationType, "query")
	ctx = context.WithValue(ctx, headers.GqlOperationName, "testOperation")
	ctx = context.WithValue(ctx, headers.GqlOperationSelectedFields, "[]")
	for _, header := range headers.AllTracingHeaders {
		ctx = context.WithValue(ctx, header, "test")
	}
	return ctx
}
