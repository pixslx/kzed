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
	"encoding/json"
	"os/exec"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kzedv1alpha1 "kzed/m/v2/api/v1alpha1"
)

// JCLJobReconciler reconciles a JCLJob object
type JCLJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type ZOWEJobSubmitOutput struct {
	Success  bool   `json:"success"`
	ExitCode int    `json:"exitCode"`
	Message  string `json:"message"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Data     struct {
		Owner         string `json:"owner"`
		Phase         int    `json:"phase"`
		Subsystem     string `json:"subsystem"`
		PhaseName     string `json:"phase-name"`
		JobCorrelator string `json:"job-correlator"`
		Type          string `json:"type"`
		URL           string `json:"url"`
		Jobid         string `json:"jobid"`
		Class         string `json:"class"`
		FilesURL      string `json:"files-url"`
		Jobname       string `json:"jobname"`
		Status        string `json:"status"`
		Retcode       string `json:"retcode"`
	} `json:"data"`
}

type ZOWEJobSpoolsOutput struct {
	Success  bool   `json:"success"`
	ExitCode int    `json:"exitCode"`
	Message  string `json:"message"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Data     []struct {
		ID       int    `json:"id"`
		DdName   string `json:"ddName"`
		StepName string `json:"stepName"`
		ProcName string `json:"procName"`
		Data     string `json:"data"`
	} `json:"data"`
}

//+kubebuilder:rbac:groups=kzed.io,resources=jcljobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kzed.io,resources=jcljobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kzed.io,resources=jcljobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JCLJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *JCLJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	jclJob := &kzedv1alpha1.JCLJob{}
	err := r.Get(ctx, req.NamespacedName, jclJob)
	logger.Info("Starting reconciliation for JCLJob " + jclJob.Name)
	if err != nil {
		return ctrl.Result{}, nil
	}

	if jclJob.Status.Status == "" {
		logger.Info("Submitting JCLJob with path " + jclJob.Spec.Path)
		zoweOut, err := exec.Command("sh", "-c", "zowe jobs submit data-set --rfj '"+jclJob.Spec.Path+"'").Output()
		logger.Info(string(zoweOut))
		if err != nil {
			logger.Error(err, "error submitting job through ZOWE CLI")
			return ctrl.Result{}, err
		}
		zoweResponse := ZOWEJobSubmitOutput{}
		err = json.Unmarshal([]byte(zoweOut), &zoweResponse)
		if err != nil {
			logger.Error(err, "error unmarshalling ZOWE CLI output")
			return ctrl.Result{}, err
		}

		jclJob.Status.Status = zoweResponse.Data.Status
		jclJob.Status.JobID = zoweResponse.Data.Jobid
		jclJob.Status.JobName = zoweResponse.Data.Jobname
		jclJob.Status.ReturnCode = zoweResponse.Data.Retcode

		r.Status().Update(ctx, jclJob)
		return ctrl.Result{}, nil
	} else {
		logger.Info("Querying JCLJob with JobID " + jclJob.Status.JobID)
		zoweOut, err := exec.Command("sh", "-c", "zowe jobs vw jsbj --rfj "+jclJob.Status.JobID).Output()
		if err != nil {
			logger.Error(err, "error querying job through ZOWE CLI")
			return ctrl.Result{}, err
		}
		zoweResponse := ZOWEJobSubmitOutput{}
		err = json.Unmarshal([]byte(zoweOut), &zoweResponse)
		if err != nil {
			logger.Error(err, "error unmarshalling ZOWE CLI output")
			return ctrl.Result{}, err
		}

		jclJob.Status.Status = zoweResponse.Data.Status
		jclJob.Status.JobID = zoweResponse.Data.Jobid
		jclJob.Status.JobName = zoweResponse.Data.Jobname
		jclJob.Status.ReturnCode = zoweResponse.Data.Retcode

		if jclJob.Status.Status == "OUTPUT" {
			logger.Info("Querying JCLJob spools with JobID " + jclJob.Status.JobID)
			zoweOut, err = exec.Command("sh", "-c", "zowe jobs vw asc --rfj "+jclJob.Status.JobID).Output()
			if err != nil {
				logger.Error(err, "error getting job spools through ZOWE CLI")
				return ctrl.Result{}, err
			}
			zoweSpoolsResponse := ZOWEJobSpoolsOutput{}
			err = json.Unmarshal([]byte(zoweOut), &zoweSpoolsResponse)
			if err != nil {
				logger.Error(err, "error unmarshalling ZOWE CLI output")
				return ctrl.Result{}, err
			}

			jclJob.Status.SpoolFiles = []kzedv1alpha1.JCLJobSpoolFiles{}
			for _, spoolFile := range zoweSpoolsResponse.Data {
				jclJobSpoolFile := &kzedv1alpha1.JCLJobSpoolFiles{}
				jclJobSpoolFile.SpoolID = string(spoolFile.ID)
				jclJobSpoolFile.StepName = spoolFile.StepName
				jclJobSpoolFile.DDName = spoolFile.DdName
				jclJobSpoolFile.Data = spoolFile.Data
				jclJob.Status.SpoolFiles = append(jclJob.Status.SpoolFiles, *jclJobSpoolFile)
			}
		}

		r.Status().Update(ctx, jclJob)
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *JCLJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kzedv1alpha1.JCLJob{}).
		Complete(r)
}
