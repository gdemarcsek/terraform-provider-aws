// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apprunner

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/apprunner/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
)

const (
	ServiceCreateTimeout = 20 * time.Minute
	ServiceDeleteTimeout = 20 * time.Minute
	ServiceUpdateTimeout = 20 * time.Minute

	ObservabilityConfigurationCreateTimeout = 2 * time.Minute
	ObservabilityConfigurationDeleteTimeout = 2 * time.Minute

	VPCIngressConnectionCreateTimeout = 2 * time.Minute
	VPCIngressConnectionDeleteTimeout = 2 * time.Minute
)

func WaitObservabilityConfigurationActive(ctx context.Context, conn *apprunner.Client, observabilityConfigurationArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{},
		Target:  []string{ObservabilityConfigurationStatusActive},
		Refresh: StatusObservabilityConfiguration(ctx, conn, observabilityConfigurationArn),
		Timeout: ObservabilityConfigurationCreateTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitObservabilityConfigurationInactive(ctx context.Context, conn *apprunner.Client, observabilityConfigurationArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{ObservabilityConfigurationStatusActive},
		Target:  []string{ObservabilityConfigurationStatusInactive},
		Refresh: StatusObservabilityConfiguration(ctx, conn, observabilityConfigurationArn),
		Timeout: ObservabilityConfigurationDeleteTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitVPCIngressConnectionActive(ctx context.Context, conn *apprunner.Client, vpcIngressConnectionArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{},
		Target:  []string{VPCIngressConnectionstatusActive},
		Refresh: StatusVPCIngressConnection(ctx, conn, vpcIngressConnectionArn),
		Timeout: VPCIngressConnectionCreateTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitVPCIngressConnectionDeleted(ctx context.Context, conn *apprunner.Client, vpcIngressConnectionArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{VPCIngressConnectionstatusActive, VPCIngressConnectionstatusPendingDeletion},
		Target:  []string{VPCIngressConnectionstatusDeleted},
		Refresh: StatusVPCIngressConnection(ctx, conn, vpcIngressConnectionArn),
		Timeout: VPCIngressConnectionDeleteTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitServiceCreated(ctx context.Context, conn *apprunner.Client, serviceArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.ServiceStatusOperationInProgress),
		Target:  enum.Slice(types.ServiceStatusRunning),
		Refresh: StatusService(ctx, conn, serviceArn),
		Timeout: ServiceCreateTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitServiceUpdated(ctx context.Context, conn *apprunner.Client, serviceArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.ServiceStatusOperationInProgress),
		Target:  enum.Slice(types.ServiceStatusRunning),
		Refresh: StatusService(ctx, conn, serviceArn),
		Timeout: ServiceUpdateTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func WaitServiceDeleted(ctx context.Context, conn *apprunner.Client, serviceArn string) error {
	stateConf := &retry.StateChangeConf{
		Pending: enum.Slice(types.ServiceStatusRunning, types.ServiceStatusOperationInProgress),
		Target:  enum.Slice(types.ServiceStatusDeleted),
		Refresh: StatusService(ctx, conn, serviceArn),
		Timeout: ServiceDeleteTimeout,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
