package ruleengine

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var rule_id int64
	var app_id int64
	var phase string
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete rule-engine --id 1234 --application-id 99887766 --phase request
		$ azion delete rule-engine
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("id") {

				answer, err := askInput(msg.AskInputRulesId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRule
				}

				rule_id = num
			}

			if !cmd.Flags().Changed("application-id") {

				answer, err := askInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRule
				}

				app_id = num
			}

			if !cmd.Flags().Changed("phase") {

				answer, err := askInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, app_id, phase, rule_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.DeleteOutputSuccess, rule_id)

			return nil
		},
	}

	cmd.Flags().Int64Var(&rule_id, "id", 0, msg.FlagRuleID)
	cmd.Flags().Int64Var(&app_id, "application-id", 0, msg.FlagAppID)
	cmd.Flags().StringVar(&phase, "phase", "", msg.FlagPhase)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

func askInput(msg string) (string, error) {
	qs := []*survey.Question{
		{
			Name:     "id",
			Prompt:   &survey.Input{Message: msg},
			Validate: survey.Required,
		},
	}

	answer := ""

	err := survey.Ask(qs, &answer)
	if err != nil {
		logger.Debug("Error while parsing answer", zap.Error(err))
		return "", utils.ErrorParseResponse
	}

	return answer, nil
}
