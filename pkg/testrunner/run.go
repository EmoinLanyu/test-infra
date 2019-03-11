// Copyright 2019 Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testrunner

import (
	"fmt"
	"time"

	tmv1beta1 "github.com/gardener/test-infra/pkg/apis/testmachinery/v1beta1"
	tmclientset "github.com/gardener/test-infra/pkg/client/testmachinery/clientset/versioned"
	"github.com/gardener/test-infra/pkg/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func runTestrun(tmClient *tmclientset.Clientset, tr *tmv1beta1.Testrun, namespace, name string) (*tmv1beta1.Testrun, error) {
	// TODO: Remove legacy name attribute. Instead enforce usage of generateName.
	tr.Name = ""
	tr.GenerateName = name
	tr.Namespace = namespace
	tr, err := tmClient.Testmachinery().Testruns(namespace).Create(tr)
	if err != nil {
		return nil, fmt.Errorf("Cannot create testrun: %s", err.Error())
	}
	log.Infof("Testrun %s deployed", tr.Name)

	testrunPhase := tmv1beta1.PhaseStatusInit
	interval := time.Duration(pollIntervalSeconds) * time.Second
	timeout := time.Duration(maxWaitTimeSeconds) * time.Second
	err = wait.PollImmediate(interval, timeout, func() (bool, error) {
		tr, err = tmClient.Testmachinery().Testruns(namespace).Get(tr.Name, metav1.GetOptions{})
		if err != nil {
			log.Errorf("Cannot get testrun: %s", err.Error())
		}

		if tr.Status.State != "" {
			testrunPhase = tr.Status.Phase
			log.Infof("Testrun %s is in %s phase. State: %s", tr.Name, testrunPhase, tr.Status.State)
		} else {
			log.Infof("Testrun %s is in %s phase. Waiting ...", tr.Name, testrunPhase)
		}
		return util.Completed(testrunPhase), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Maximum wait time of %d is exceeded by Testrun %s", maxWaitTimeSeconds, name)
	}

	return tr, nil
}