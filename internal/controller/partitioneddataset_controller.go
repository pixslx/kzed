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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kzedv1alpha1 "kzed/m/v2/api/v1alpha1"
	"kzed/m/v2/zowe"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const pdsFinalizer = "partitioneddatasets.kzed.io/finalizer"

// PartitionedDataSetReconciler reconciles a PartitionedDataSet object
type PartitionedDataSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Zowe   zowe.ZOWE
	SYSUID string
}

//+kubebuilder:rbac:groups=kzed.io,resources=partitioneddatasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kzed.io,resources=partitioneddatasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kzed.io,resources=partitioneddatasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PartitionedDataSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *PartitionedDataSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	pds := &kzedv1alpha1.PartitionedDataSet{}
	err := r.Get(ctx, req.NamespacedName, pds)
	logger.Info("starting reconciliation for PartitionedDataSet " + pds.Name)
	if err != nil {
		logger.Info("deleting PartitionedDataSet " + pds.Name)
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(pds, pdsFinalizer) {
		controllerutil.AddFinalizer(pds, pdsFinalizer)
		err = r.Update(ctx, pds)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	isPDSToBeDeleted := pds.GetDeletionTimestamp() != nil
	if isPDSToBeDeleted {
		if controllerutil.ContainsFinalizer(pds, pdsFinalizer) {
			if err := r.finalize(logger, pds); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(pds, pdsFinalizer)
			err := r.Update(ctx, pds)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	pdsPath := strings.ToUpper(r.SYSUID + "." + pds.Name)

	if pds.Status.Status == "" {
		logger.Info("creating PartitionedDataSet " + pds.Name)
		zoweInput := zowe.ZOWEDataSetCreateInput{}
		zoweInput.Name = pdsPath
		zoweInput.AllocationSpaceUnit = pds.Spec.AllocationSpaceUnit
		zoweInput.BlockSize = pds.Spec.BlockSize
		zoweInput.DataClass = pds.Spec.DataClass
		zoweInput.DataSetType = pds.Spec.DataSetType
		zoweInput.DeviceType = pds.Spec.DeviceType
		zoweInput.DirectoryBlocks = pds.Spec.DirectoryBlocks
		zoweInput.ManagementClass = pds.Spec.ManagementClass
		zoweInput.PrimarySpace = pds.Spec.PrimarySpace
		zoweInput.RecordFormat = pds.Spec.RecordFormat
		zoweInput.RecordLength = pds.Spec.RecordLength
		zoweInput.SecondarySpace = pds.Spec.SecondarySpace
		zoweInput.Size = pds.Spec.Size
		zoweInput.StorageClass = pds.Spec.StorageClass
		zoweInput.VolumeSerial = pds.Spec.VolumeSerial
		zoweInput.Like = pds.Spec.Like
		zoweResponse, err := r.Zowe.FilesCreatePDS(zoweInput)

		if err != nil {
			logger.Error(err, "error creating PartitionedDataSet through ZOWE CLI")
			return ctrl.Result{}, err
		}
		if zoweResponse.ExitCode == 0 {
			pds.Status.Status = "CREATED"
			pds.Status.CreatedAt.Time = time.Now()
			pds.Status.LastSyncAt.Time = time.Now()
			r.Status().Update(ctx, pds)
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, errors.New("error creating PartitionedDataSet through ZOWE CLI")
		}
	} else if pds.Status.Status == "CREATED" || pds.Status.Status == "SYNCED" {
		membersInCR := map[string]bool{}
		for fileName, fileContent := range pds.Data {
			membersInCR[fileName] = true
			zoweResponse, err := r.Zowe.FilesUploadSTDIN2DS(pdsPath+"("+fileName+")", fileContent)
			if err != nil {
				logger.Error(err, "error uploading content to PartitionedDataSet through ZOWE CLI")
				return ctrl.Result{}, err
			}
			if zoweResponse.ExitCode != 0 {
				return ctrl.Result{}, errors.New("error uploading content to PartitionedDataSet")
			}
		}

		zoweResponse, err := r.Zowe.FilesDSListMembers(pdsPath)
		if err == nil {
			for _, pdsItem := range zoweResponse.Data.APIResponse.Items {
				_, existsInCR := membersInCR[pdsItem.Member]
				if !existsInCR {
					err = r.Zowe.FilesDSDelete(strings.ToUpper(pdsPath + "(" + pdsItem.Member + ")"))
					if err != nil {
						logger.Error(err, "error deleting member from PartitionedDataSet through ZOWE CLI")
					}
				}
			}
		}

		pds.Status.Status = "SYNCED"
		pds.Status.LastSyncAt.Time = time.Now()
		r.Status().Update(ctx, pds)

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PartitionedDataSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kzedv1alpha1.PartitionedDataSet{}).
		WithEventFilter(
			predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.Funcs{UpdateFunc: func(ue event.UpdateEvent) bool {
					return ue.ObjectNew.(*kzedv1alpha1.PartitionedDataSet).Status.Status == "CREATED"
				}},
			),
		).
		Complete(r)
}

func (r *PartitionedDataSetReconciler) finalize(logger logr.Logger, pds *kzedv1alpha1.PartitionedDataSet) error {
	logger.Info("finalizing PartitionedDataSet " + pds.Name)
	pdsPath := strings.ToUpper(r.SYSUID + "." + pds.Name)
	if r.Zowe.FilesDSExists(pdsPath) {
		err := r.Zowe.FilesDSDelete(pdsPath)
		return err
	}
	return nil
}
