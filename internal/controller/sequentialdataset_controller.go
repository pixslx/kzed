/*
Copyright 2024.

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

package controller

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kzedv1alpha1 "kzed/m/v2/api/v1alpha1"
	"kzed/m/v2/zowe"
)

const sdsFinalizer = "sequentialdatasets.kzed.io/finalizer"

// SequentialDataSetReconciler reconciles a SequentialDataSet object
type SequentialDataSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Zowe   zowe.ZOWE
	SYSUID string
}

//+kubebuilder:rbac:groups=kzed.io,resources=sequentialdatasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kzed.io,resources=sequentialdatasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kzed.io,resources=sequentialdatasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SequentialDataSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *SequentialDataSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	sds := &kzedv1alpha1.SequentialDataSet{}
	err := r.Get(ctx, req.NamespacedName, sds)
	logger.Info("starting reconciliation for SequentialDataSet " + sds.Name)
	if err != nil {
		logger.Info("deleting SequentialDataSet " + sds.Name)
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(sds, sdsFinalizer) {
		controllerutil.AddFinalizer(sds, sdsFinalizer)
		err = r.Update(ctx, sds)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	isSDSToBeDeleted := sds.GetDeletionTimestamp() != nil
	if isSDSToBeDeleted {
		if controllerutil.ContainsFinalizer(sds, sdsFinalizer) {
			if err := r.finalize(logger, sds); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(sds, sdsFinalizer)
			err := r.Update(ctx, sds)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	sdsPath := strings.ToUpper(r.SYSUID + "." + sds.Name)

	if sds.Status.Status == "" {
		logger.Info("creating SequentialDataSet " + sds.Name)
		zoweInput := zowe.ZOWEDataSetCreateInput{}
		zoweInput.Name = sdsPath
		zoweInput.BlockSize = sds.Spec.BlockSize
		zoweInput.DataClass = sds.Spec.DataClass
		zoweInput.DeviceType = sds.Spec.DeviceType
		zoweInput.DirectoryBlocks = sds.Spec.DirectoryBlocks
		zoweInput.ManagementClass = sds.Spec.ManagementClass
		zoweInput.PrimarySpace = sds.Spec.PrimarySpace
		zoweInput.RecordFormat = sds.Spec.RecordFormat
		zoweInput.RecordLength = sds.Spec.RecordLength
		zoweInput.SecondarySpace = sds.Spec.SecondarySpace
		zoweInput.Size = sds.Spec.Size
		zoweInput.StorageClass = sds.Spec.StorageClass
		zoweInput.VolumeSerial = sds.Spec.VolumeSerial
		zoweInput.Like = sds.Spec.Like
		zoweResponse, err := r.Zowe.FilesCreateSDS(zoweInput)

		if err != nil {
			logger.Error(err, "error creating SequentialDataSet through ZOWE CLI")
			return ctrl.Result{}, err
		}
		if zoweResponse.ExitCode == 0 {
			sds.Status.Status = "CREATED"
			sds.Status.CreatedAt.Time = time.Now()
			sds.Status.LastSyncAt.Time = time.Now()
			r.Status().Update(ctx, sds)
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, errors.New("error creating SequentialDataSet through ZOWE CLI")
		}
	} else if sds.Status.Status == "CREATED" || sds.Status.Status == "SYNCED" {
		if sds.Data != "" {
			zoweResponse, err := r.Zowe.FilesUploadSTDIN2DS(strings.ToUpper(sdsPath), sds.Data)
			if err != nil {
				logger.Error(err, "error uploading content to SequentialDataSet through ZOWE CLI")
				return ctrl.Result{}, err
			}
			if zoweResponse.ExitCode != 0 {
				return ctrl.Result{}, errors.New("error uploading content to SequentialDataSet")
			}

			sds.Status.Status = "SYNCED"
			sds.Status.LastSyncAt.Time = time.Now()
			r.Status().Update(ctx, sds)
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SequentialDataSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kzedv1alpha1.SequentialDataSet{}).
		WithEventFilter(
			predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.Funcs{UpdateFunc: func(ue event.UpdateEvent) bool {
					return ue.ObjectNew.(*kzedv1alpha1.SequentialDataSet).Status.Status == "CREATED"
				}},
			),
		).
		Complete(r)
}

func (r *SequentialDataSetReconciler) finalize(logger logr.Logger, pds *kzedv1alpha1.SequentialDataSet) error {
	logger.Info("finalizing PartitionedDataSet " + pds.Name)
	sdsPath := strings.ToUpper(r.SYSUID + "." + pds.Name)
	if r.Zowe.FilesDSExists(sdsPath) {
		err := r.Zowe.FilesDSDelete(sdsPath)
		return err
	}
	return nil
}
