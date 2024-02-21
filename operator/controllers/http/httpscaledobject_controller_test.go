package http

import (
	"github.com/kedacore/http-add-on/operator/controllers/http/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHttpScaledObjectController(t *testing.T) {
	r := require.New(t)

	testInfra := newCommonTestInfra("testns", "testapp")

	reconciller := &HTTPScaledObjectReconciler{
		Client:               testInfra.cl,
		Scheme:               testInfra.cl.Scheme(),
		ExternalScalerConfig: config.ExternalScaler{},
		BaseConfig:           config.Base{},
	}

	// Create required app objects for the application defined by the CRD
	err := reconciller.createOrUpdateApplicationResources(
		testInfra.ctx,
		testInfra.logger,
		testInfra.cl,
		config.Base{},
		config.ExternalScaler{},
		&testInfra.httpso,
	)
	r.NoError(err)

	// check for scaledobject, expect no error as scaledobject should get created
	_, err = getSO(
		testInfra.ctx,
		testInfra.cl,
		testInfra.httpso,
	)
	r.NoError(err)

}

func TestHttpScaledObjectControllerWhenSkipAnnotationSet(t *testing.T) {
	r := require.New(t)

	testInfra := newCommonTestInfraWithSkipScaledObjectCreation("testns", "testapp")

	reconciller := &HTTPScaledObjectReconciler{
		Client:               testInfra.cl,
		Scheme:               testInfra.cl.Scheme(),
		ExternalScalerConfig: config.ExternalScaler{},
		BaseConfig:           config.Base{},
	}

	// Create required app objects for the application defined by the CRD
	err := reconciller.createOrUpdateApplicationResources(
		testInfra.ctx,
		testInfra.logger,
		testInfra.cl,
		config.Base{},
		config.ExternalScaler{},
		&testInfra.httpso,
	)
	r.NoError(err)

	// check for scaledobject, expect error as scaledobject should not exist when skipScaledObjectCreation annotation is set
	_, err = getSO(
		testInfra.ctx,
		testInfra.cl,
		testInfra.httpso,
	)
	r.Error(err)

}
