/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package vm

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"

	"kubevirt.io/kubevirt/pkg/kubecli"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

const (
	COMMAND_START   = "start"
	COMMAND_STOP    = "stop"
	COMMAND_RESTART = "restart"
)

type VMCommand struct {
	clientConfig clientcmd.ClientConfig
	out          io.Writer
}

func GetCommands(clientConfig clientcmd.ClientConfig, out io.Writer) *VMCommand {
	return &VMCommand{out: out, clientConfig: clientConfig}
}

func (vmc *VMCommand) Start() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start (VM)",
		Short:   "Start a virtual machine.",
		Example: usage(COMMAND_START),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vmc.Run(COMMAND_START, args)
		},
	}
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (vmc *VMCommand) Stop() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop (VM)",
		Short:   "Stop a virtual machine.",
		Example: usage(COMMAND_STOP),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vmc.Run(COMMAND_STOP, args)
		},
	}
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (vmc *VMCommand) Restart() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restart (VM)",
		Short:   "Restart a virtual machine.",
		Example: usage(COMMAND_RESTART),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vmc.Run(COMMAND_RESTART, args)
		},
	}
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func usage(cmd string) string {
	usage := fmt.Sprintf("  # %s a virtual machine called 'myvm':\n", strings.Title(cmd))
	usage += fmt.Sprintf("  virtctl %s myvm", cmd)
	return usage
}

func (vmc *VMCommand) Run(command string, args []string) error {

	vmiName := args[0]

	namespace, _, err := vmc.clientConfig.Namespace()
	if err != nil {
		return err
	}

	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(vmc.clientConfig)
	if err != nil {
		return fmt.Errorf("Cannot obtain KubeVirt client: %v", err)
	}

	options := k8smetav1.GetOptions{}
	vm, err := virtClient.VirtualMachine(namespace).Get(vmiName, &options)
	if err != nil {
		return fmt.Errorf("Error fetching VirtualMachine: %v", err)
	}

	switch command {
	case COMMAND_STOP, COMMAND_START:
		running := false
		if command == COMMAND_START {
			running = true
		}

		if vm.Spec.Running != running {
			bodyStr := fmt.Sprintf("{\"spec\":{\"running\":%t}}", running)

			_, err := virtClient.VirtualMachine(namespace).Patch(vm.Name, types.MergePatchType,
				[]byte(bodyStr))

			if err != nil {
				return fmt.Errorf("Error updating VirtualMachine: %v", err)
			}

		} else {
			stateMsg := "stopped"
			if running {
				stateMsg = "running"
			}
			return fmt.Errorf("Error: VirtualMachineInstance '%s' is already %s", vmiName, stateMsg)
		}
	case COMMAND_RESTART:
		err = virtClient.VirtualMachine(namespace).Restart(vmiName)
		if err != nil {
			return fmt.Errorf("Error restarting VirtualMachine %v", err)
		}
	}

	fmt.Fprintf(vmc.out, "VM %s was scheduled to %s\n", vmiName, command)
	return nil
}
