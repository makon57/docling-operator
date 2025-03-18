package reconcilers

import (
	"context"

	"github.com/go-logr/logr"
	"github.io/opdev/docling-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewServiceReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *ServiceReconciler {
	return &ServiceReconciler{
		Client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *ServiceReconciler) Reconcile(ctx context.Context, doclingServ *v1alpha1.DoclingServ) (bool, error) {

	service := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: doclingServ.Name, Namespace: doclingServ.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Service", "Service", doclingServ.Name)
		svc := r.createService(doclingServ)

		err = r.Create(ctx, svc)
		if err != nil {
			r.Log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return true, err
		}
	}

	return false, nil
}

func (r *ServiceReconciler) createService(doclingServ *v1alpha1.DoclingServ) *corev1.Service {

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:		doclingServ.Name,
			Namespace:  doclingServ.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 5001,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": "docling-serve",
			},
		},
	}
	return service
}
