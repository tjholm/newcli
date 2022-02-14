// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure/sdk/v4/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/cli/pkg/provider/pulumi/common"
)

type StorageArgs struct {
	ResourceGroupName pulumi.StringInput
}

type Storage struct {
	pulumi.ResourceState

	Name       string
	Account    *storage.Account
	Queues     map[string]*storage.Queue
	Containers map[string]*storage.Container
}

func (a *azureProvider) newStorageResources(ctx *pulumi.Context, name string, args *StorageArgs, opts ...pulumi.ResourceOption) (*Storage, error) {
	res := &Storage{
		Name:       name,
		Queues:     map[string]*storage.Queue{},
		Containers: map[string]*storage.Container{},
	}
	err := ctx.RegisterComponentResource("nitric:storage:AzureStorage", name, res, opts...)
	if err != nil {
		return nil, err
	}

	accName := resourceName(ctx, name, StorageAccountRT)
	res.Account, err = storage.NewAccount(ctx, accName, &storage.AccountArgs{
		AccessTier:             pulumi.String("Hot"),
		ResourceGroupName:      args.ResourceGroupName,
		AccountKind:            pulumi.String("StorageV2"),
		AccountTier:            pulumi.String("Standard"),
		AccountReplicationType: pulumi.String("LRS"),
		Tags:                   common.Tags(ctx, accName),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, errors.WithMessage(err, "account create")
	}

	for bName := range a.s.Buckets {
		res.Containers[bName], err = storage.NewContainer(ctx, resourceName(ctx, bName, StorageContainerRT), &storage.ContainerArgs{
			StorageAccountName: res.Account.Name,
		}, pulumi.Parent(res))
		if err != nil {
			return nil, errors.WithMessage(err, "container create")
		}
	}

	for qName := range a.s.Queues {
		res.Queues[qName], err = storage.NewQueue(ctx, resourceName(ctx, qName, StorageQueueRT), &storage.QueueArgs{
			StorageAccountName: res.Account.Name,
		}, pulumi.Parent(res))
		if err != nil {
			return nil, errors.WithMessage(err, "queue create")
		}
	}
	return res, nil
}
