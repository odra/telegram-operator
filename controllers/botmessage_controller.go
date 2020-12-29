/*
Copyright 2020.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	telegramv1alpha1 "github.com/odra/telegram-operator/api/v1alpha1"
)

// BotMessageReconciler reconciles a BotMessage object
type BotMessageReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=telegram.my.domain,resources=botmessages,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=telegram.my.domain,resources=botmessages/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=telegram.my.domain,resources=botmessages/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BotMessage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *BotMessageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	l := r.Log.WithValues("botmessage", req.NamespacedName)

	instance := &telegramv1alpha1.BotMessage{}

	err = r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	if instance.Status.Type == telegramv1alpha1.BotMessageError {
		l.Info("Message error status, changing status to new again..")
		err = r.setStatusToNew(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if instance.Status.Type == telegramv1alpha1.BotMessageSent {
		l.Info("Message already sent, do nothing")
		return ctrl.Result{}, nil
	}

	if instance.Status.Type == telegramv1alpha1.BotMessageNew {
		err = r.setStatusToSending(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		l.Info("Provision pod and set a reference in the CR annotation")
		return ctrl.Result{Requeue: true}, nil
	}

	if instance.Status.Type == telegramv1alpha1.BotMessageSending {
		err = r.setStatusToCompleted(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	l.Info("Nothing new, finishing reconcile loop")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BotMessageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&telegramv1alpha1.BotMessage{}).
		Owns(&v12.Pod{}).
		Complete(r)
}

//controller utils
func (r *BotMessageReconciler) setStatusToNew(ctx context.Context, instance *telegramv1alpha1.BotMessage) error {
	instance.Status.Type = telegramv1alpha1.BotMessageNew
	instance.Status.Status = v1.ConditionTrue
	instance.Status.Reason = "New"
	instance.Status.Message = "New message to be sent"

	return r.Client.Status().Update(ctx, instance)
}

func (r *BotMessageReconciler) setStatusToSending(ctx context.Context, instance *telegramv1alpha1.BotMessage) error {
	instance.Status.Type = telegramv1alpha1.BotMessageSending
	instance.Status.Status = v1.ConditionTrue
	instance.Status.Reason = "Sending"
	instance.Status.Message = "Sending Message"

	pod := &v12.Pod{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "telegram-sender-" + instance.Name,
			Namespace: instance.Namespace,
			Labels:    nil,
		},
		Spec: v12.PodSpec{
			RestartPolicy: v12.RestartPolicyOnFailure,
			Containers: []v12.Container{
				{
					Name:  "sender",
					Image: instance.Spec.Image,
					Args:  []string{instance.Spec.Text},
					EnvFrom: []v12.EnvFromSource{
						{
							SecretRef: &v12.SecretEnvSource{
								LocalObjectReference: v12.LocalObjectReference{
									Name: instance.Spec.Secret.Name,
								},
							},
						},
					},
				},
			},
		},
	}

	err := r.Client.Create(ctx, pod)
	if err != nil {
		return err
	}

	return r.Client.Status().Update(ctx, instance)
}

func (r *BotMessageReconciler) setStatusToCompleted(ctx context.Context, instance *telegramv1alpha1.BotMessage) error {
	changed := false
	pod := &v12.Pod{}

	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      "telegram-sender-" + instance.Name,
	}, pod)
	if err != nil {
		return err
	}

	if pod.Status.Phase == v12.PodSucceeded {
		instance.Status.Type = telegramv1alpha1.BotMessageSent
		instance.Status.Status = v1.ConditionTrue
		instance.Status.Reason = "Sent"
		instance.Status.Message = "Message successfully sent"
		changed = true
	}

	if pod.Status.Phase == v12.PodFailed {
		instance.Status.Type = telegramv1alpha1.BotMessageError
		instance.Status.Status = v1.ConditionTrue
		instance.Status.Reason = "Faailed"
		instance.Status.Message = "Sender pod failed, please check pod logs"
		changed = true
	}

	if changed {
		return r.Client.Status().Update(ctx, instance)
	}

	return nil
}
