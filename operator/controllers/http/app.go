package http

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	"github.com/kedacore/http-add-on/operator/controllers/http/config"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
)

func (r *HTTPScaledObjectReconciler) createOrUpdateApplicationResources(
	ctx context.Context,
	logger logr.Logger,
	cl client.Client,
	baseConfig config.Base,
	externalScalerConfig config.ExternalScaler,
	httpso *v1alpha1.HTTPScaledObject,
) error {
	defer SaveStatus(context.Background(), logger, cl, httpso)
	logger = logger.WithValues(
		"reconciler.appObjects",
		"addObjects",
		"HTTPScaledObject.name",
		httpso.Name,
		"HTTPScaledObject.namespace",
		httpso.Namespace,
	)

	// set initial statuses
	AddOrUpdateCondition(
		httpso,
		*SetMessage(
			CreateCondition(
				v1alpha1.Pending,
				v1.ConditionUnknown,
				v1alpha1.PendingCreation,
			),
			"Identified HTTPScaledObject creation signal"),
	)

	// We want to integrate http scaler with other
	// scalers. when SkipScaledObjectCreation is set to true,
<<<<<<< HEAD
	// reconciler will skip the KEDA core ScaledObjects creation or delete the ScaledOBject if it already exists.
	// you can then create your own ScaledObject, and add http scaler as one of your triggers.
=======
	// reconciler will skip the KEDA core ScaledObjects creation.
	// you can create your own so, and add http scaler as one of your triggers.
>>>>>>> 501c6c1 (feat: provide support to allow HTTP scaler to work alongside other core KEDA scalers)
	if httpso.Annotations["skipScaledObjectCreation"] == "true" {
		logger.Info(
			"Skip scaled objects creation with flag SkipScaledObjectCreation=true",
			"HTTPScaledObject", httpso.Name)
<<<<<<< HEAD
		err := r.deleteScaledObject(ctx, cl, logger, httpso)
		if err != nil {
			logger.Info("Failed to delete ScaledObject",
				"HTTPScaledObject", httpso.Name)
		}
=======
>>>>>>> 501c6c1 (feat: provide support to allow HTTP scaler to work alongside other core KEDA scalers)
		return nil
	}

	// create the KEDA core ScaledObjects (not the HTTP one) for
	// the app deployment and the interceptor deployment.
	// this needs to be submitted so that KEDA will scale both the app and
	// interceptor
	return r.createOrUpdateScaledObject(
		ctx,
		cl,
		logger,
		externalScalerConfig.HostName(baseConfig.CurrentNamespace),
		httpso,
	)
}

func (r *HTTPScaledObjectReconciler) deleteScaledObject(
	ctx context.Context,
	cl client.Client,
	logger logr.Logger,
	httpso *v1alpha1.HTTPScaledObject,
) error {

	var fetchedSO kedav1alpha1.ScaledObject

	objectKey := types.NamespacedName{
		Namespace: httpso.Namespace,
		Name:      httpso.Name,
	}

	if err := cl.Get(ctx, objectKey, &fetchedSO); err != nil {
		logger.Info("Failed to retreive ScaledObject",
			"ScaledObject", &fetchedSO.Name)
		return nil
	}

	if isOwnerReferenceMatch(&fetchedSO, httpso) {
		if err := cl.Delete(ctx, &fetchedSO); err != nil {
			logger.Info("Failed to delete ScaledObject",
				"ScaledObject", &fetchedSO.Name)
			return nil
		}
		logger.Info("Deleted ScaledObject",
			"ScaledObject", &fetchedSO.Name)
	}

	return nil
}

// function to check if the owner reference of ScaledObject matches the HTTPScaledObject
func isOwnerReferenceMatch(scaledObject *kedav1alpha1.ScaledObject, httpso *v1alpha1.HTTPScaledObject) bool {

	for _, ownerRef := range scaledObject.OwnerReferences {
		if strings.ToLower(ownerRef.Kind) == "httpscaledobject" &&
			ownerRef.Name == httpso.Name {
			return true
		}
	}
	return false
}
