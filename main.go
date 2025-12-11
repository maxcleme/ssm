package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

const ssmPrefix = "ssm://"

var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Run command with SSM parameter substitution",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		client := ssm.NewFromConfig(cfg)
		env, err := resolveSSMVars(ctx, client, os.Environ())
		if err != nil {
			return err
		}
		command := exec.Command(args[0], args[1:]...)
		command.Env = env
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		return command.Run()
	},
}

func resolveSSMVars(ctx context.Context, client *ssm.Client, environ []string) ([]string, error) {
	var result []string
	var errs error
	for _, env := range environ {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			result = append(result, env)
			continue
		}
		value := parts[1]
		if strings.HasPrefix(value, ssmPrefix) {
			name := strings.TrimPrefix(value, ssmPrefix)
			param, err := client.GetParameter(ctx, &ssm.GetParameterInput{
				Name:           aws.String(name),
				WithDecryption: aws.Bool(true),
			})
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("resolving %s: %w", value, err))
				continue
			}
			value = *param.Parameter.Value
		}
		result = append(result, parts[0]+"="+value)
	}
	return result, errs
}

func main() {
	rootCmd := &cobra.Command{Use: "ssm"}
	rootCmd.AddCommand(runCmd)
	rootCmd.Execute()
}
