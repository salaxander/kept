package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/salaxander/kept/pkg/check"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"strconv"
)

const (
	yamlStage     = "stage"
	yamlMilestone = "milestone"
	yamlStatus    = "status" // `implementable`, if graduating needs to be updated to `implemented`

	alpha      = "alpha"
	beta       = "beta"
	graduating = "graduating"
	stable     = "stable"
)

// options for running `meshctl check`, which checks all clusters
type checkOpts struct {
	url     string
	release string
	stage   string
}

func (o *checkOpts) addToFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&o.url,
		"url",
		"",
		"url to KEP yaml",
	)
	flags.StringVarP(
		&o.release,
		"release",
		"r",
		"",
		"release version",
	)
	flags.StringVarP(
		&o.stage,
		"stage",
		"s",
		"",
		"stage (alpha/beta/graduating/stable)",
	)
}

func Command() *cobra.Command {
	opts := &checkOpts{}

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check a KEP complies with template for a specific version",
		Long:  `Check a KEP by providing an individual KEP yaml, the release it is targeting and the stage (alpha/beta/stable).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runYamlChecks(opts)
		},
	}

	opts.addToFlags(cmd.Flags())

	return cmd
}

func runYamlChecks(opts *checkOpts) error {
	pterm.Println("")

	err := opts.validate()
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	kepYaml, err := check.GetKepYaml(opts.url)
	if err != nil {
		return err
	}

	var kepMap map[string]interface{}
	err = yaml.Unmarshal([]byte(kepYaml), &kepMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %s", err)
	}

	// TODO: imo this should be versioned
	urlTemplateForLatest := "https://github.com/kubernetes/enhancements/blob/master/keps/NNNN-kep-template/kep.yaml"
	templateYaml, err := check.GetKepYaml(urlTemplateForLatest)
	if err != nil {
		return err
	}

	var templateMap map[string]interface{}
	err = yaml.Unmarshal([]byte(templateYaml), &templateMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %s", err)
	}

	// Check all fields in the template exist in the current KEP
	for key := range templateMap {
		if _, exists := kepMap[key]; !exists {
			fmt.Printf("Field '%s' does not exist in KEP.\n", key)
		}
	}

	// Check that correct stage is set in yaml
	if _, exists := kepMap[yamlStage]; !exists {
		fmt.Printf("stage '%s' does not exist in KEP yaml.\n", yamlStage)
	} else {
		if kepMap[yamlStage] != opts.stage {
			fmt.Printf("stage '%s' is not set to '%s' in KEP yaml.\n", yamlStage, opts.stage)
		}
	}

	// Check that correct release is set in yaml
	if _, exists := kepMap[yamlMilestone]; !exists {
		fmt.Printf("stage '%s' does not exist in KEP yaml.\n", yamlMilestone)
	} else {
		milestoneMap := convertToMapString(kepMap[yamlMilestone].(map[interface{}]interface{}))
		if err != nil {
			return fmt.Errorf("failed to unmarshal YAML: %s", err)
		}
		if _, stageExists := milestoneMap[opts.stage]; !stageExists {
			return fmt.Errorf("stage '%s' does not exist in KEP yaml for milestone %s", opts.stage, milestoneMap)
		} else {
			if milestoneMap[opts.stage] != opts.release {
				fmt.Printf("release '%s' is not set to '%s' in KEP yaml.\n", milestoneMap[opts.stage], opts.release)
			}
		}
	}

	// Check that correct status is set in yaml
	if opts.stage == graduating || opts.stage == stable {
		if _, exists := kepMap[yamlStatus]; !exists {
			fmt.Printf("status '%s' does not exist in KEP yaml.\n", yamlStatus)
		} else {
			if kepMap[yamlStatus] != "implemented" {
				fmt.Printf("status '%s' is not set to '%s' in KEP yaml.\n", yamlStatus, "implemented")
			}
		}
	} else {
		// alpha/beta
		if _, exists := kepMap[yamlStatus]; !exists {
			fmt.Printf("status '%s' does not exist in KEP yaml.\n", yamlStatus)
		} else {
			if kepMap[yamlStatus] != "implementable" {
				fmt.Printf("status '%s' is not set to '%s' in KEP yaml.\n", yamlStatus, "implementable")
			}
		}
	}

	fmt.Println("KEP passed checks.")
	return nil
}

func (opts *checkOpts) validate() error {
	if opts.stage != alpha && opts.stage != beta && opts.stage != stable && opts.stage != graduating {
		return fmt.Errorf("Invalid stage type '%s'. Supported output types are: %s, %s, %s, %s", opts.stage, alpha, beta, stable, graduating)
	}

	if opts.release == "" {
		return fmt.Errorf("Invalid release version '%s'", opts.release)
	}

	if opts.url == "" {
		return fmt.Errorf("Invalid url '%s'", opts.url)
	}

	return nil
}

func init() {
	checkCmd := Command()

	rootCmd.AddCommand(checkCmd)
}

// Recursive function to convert a map[interface{}]interface{} to map[string]string
func convertToMapString(data map[interface{}]interface{}) map[string]string {
	strMap := make(map[string]string)
	for key, value := range data {
		strKey := convertToString(key)
		strValue := convertToString(value)
		strMap[strKey] = strValue
	}
	return strMap
}

// Recursive function to convert an interface value to a string
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	case map[interface{}]interface{}:
		strMap := make(map[string]string)
		for k, v := range v {
			strMap[convertToString(k)] = convertToString(v)
		}
		return fmt.Sprintf("%v", strMap)
	default:
		return fmt.Sprintf("%v", v)
	}
}
