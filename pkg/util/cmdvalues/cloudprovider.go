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

package cmdvalues

import (
	"fmt"
	"github.com/gardener/test-infra/pkg/common"
	"github.com/spf13/pflag"
)

type CloudProviderValue struct {
	allowedProvider map[common.CloudProvider]bool
	cloudprovider   *common.CloudProvider
}

func NewCloudProviderValue(value *common.CloudProvider, defaultValue common.CloudProvider, allowed ...common.CloudProvider) pflag.Value {
	*value = defaultValue
	cpvalue := &CloudProviderValue{
		allowedProvider: make(map[common.CloudProvider]bool),
		cloudprovider:   value,
	}
	for _, cp := range allowed {
		cpvalue.allowedProvider[cp] = true
	}
	return cpvalue
}

var _ pflag.Value = &CloudProviderValue{}

func (v *CloudProviderValue) String() string {
	return string(*v.cloudprovider)
}

func (v *CloudProviderValue) Type() string {
	return "CloudProvider"
}

func (v *CloudProviderValue) Set(value string) error {
	provider := common.CloudProvider(value)
	if _, ok := v.allowedProvider[provider]; len(v.allowedProvider) != 0 && !ok {
		return fmt.Errorf("unsupported cloudprovider %s", provider)
	}
	*v.cloudprovider = provider
	return nil
}
