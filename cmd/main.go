package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	modelName               string
	modelServerName         string
	modelServerVersion      string // Now used only for manifests
	acceleratorType         string
	targetNtpotMilliseconds int
)

func main() {
	// Root command
	var rootCmd = &cobra.Command{
		Use:   "mock-gcloud",
		Short: "A mock gcloud CLI for interacting with the AI Recommender API",
	}

	// Container command
	var containerCmd = &cobra.Command{
		Use:   "container",
		Short: "Manage container resources",
	}

	// AI command
	var aiCmd = &cobra.Command{
		Use:   "ai",
		Short: "Manage AI resources",
	}

	// Recommender command
	var recommenderCmd = &cobra.Command{
		Use:   "recommender",
		Short: "Manage AI recommender resources",
	}

	// Models command
	var modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "Manage models",
	}
	modelsCmd.AddCommand(listModelsCmd())

	// Model Servers command
	var modelServersCmd = &cobra.Command{
		Use:   "model-servers",
		Short: "Manage model servers",
	}
	modelServersCmd.PersistentFlags().StringVar(&modelName, "model", "", "Model name")
	modelServersCmd.MarkPersistentFlagRequired("model")
	modelServersCmd.AddCommand(listModelServersCmd())

	// Model Server Versions command
	var modelServerVersionsCmd = &cobra.Command{
		Use:   "model-server-versions",
		Short: "Manage model server versions",
	}
	modelServerVersionsCmd.PersistentFlags().StringVar(&modelName, "model", "", "Model name")
	modelServerVersionsCmd.PersistentFlags().StringVar(&modelServerName, "model-server", "", "Model server name")
	modelServerVersionsCmd.MarkPersistentFlagRequired("model")
	modelServerVersionsCmd.MarkPersistentFlagRequired("model-server")
	modelServerVersionsCmd.AddCommand(listModelServerVersionsCmd())

	// Accelerators command
	var acceleratorsCmd = &cobra.Command{
		Use:   "accelerators",
		Short: "Manage accelerators",
	}
	acceleratorsCmd.PersistentFlags().StringVar(&modelName, "model", "", "Model name")
	acceleratorsCmd.PersistentFlags().StringVar(&modelServerName, "model-server", "", "Model server name")
	acceleratorsCmd.MarkPersistentFlagRequired("model")
	acceleratorsCmd.MarkPersistentFlagRequired("model-server")
	acceleratorsCmd.AddCommand(listAcceleratorsCmd())

	// Manifests command
	var manifestsCmd = &cobra.Command{
		Use:   "manifests",
		Short: "Manage manifests",
	}
	manifestsCmd.PersistentFlags().StringVar(&modelName, "model", "", "Model name")
	manifestsCmd.PersistentFlags().StringVar(&modelServerName, "model-server", "", "Model server name")
	manifestsCmd.PersistentFlags().StringVar(&modelServerVersion, "model-server-version", "", "Model server version")
	manifestsCmd.PersistentFlags().StringVar(&acceleratorType, "accelerator-type", "", "Accelerator type")
	manifestsCmd.PersistentFlags().IntVar(&targetNtpotMilliseconds, "target-ntpot-milliseconds", 0, "Target NTPOT milliseconds")
	manifestsCmd.MarkPersistentFlagRequired("model")
	manifestsCmd.MarkPersistentFlagRequired("model-server")
	manifestsCmd.MarkPersistentFlagRequired("model-server-version")
	manifestsCmd.MarkPersistentFlagRequired("accelerator-type")
	manifestsCmd.AddCommand(createManifestCmd())

	// Models and Servers command
	var modelsAndServersCmd = &cobra.Command{
		Use:   "modelsAndServers",
		Short: "Manage models and servers",
	}
	modelsAndServersCmd.AddCommand(listModelsAndServersCmd())

	recommenderCmd.AddCommand(modelsCmd)
	recommenderCmd.AddCommand(modelServersCmd)
	recommenderCmd.AddCommand(modelServerVersionsCmd)
	recommenderCmd.AddCommand(acceleratorsCmd)
	recommenderCmd.AddCommand(manifestsCmd)
	recommenderCmd.AddCommand(modelsAndServersCmd)
	aiCmd.AddCommand(recommenderCmd)
	containerCmd.AddCommand(aiCmd)
	rootCmd.AddCommand(containerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func listModelsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List models",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewAPIClient()
			models, err := client.ListModels()
			if err != nil {
				log.Fatalf("Error listing models: %v", err)
			}
			for _, model := range models {
				fmt.Printf("Name: %s\n", model.Name)
			}
		},
	}
}

func listModelServersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List model servers",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewAPIClient()
			modelServers, err := client.ListModelServers(modelName)
			if err != nil {
				log.Fatalf("Error listing model servers: %v", err)
			}
			for _, modelServer := range modelServers {
				fmt.Printf("Name: %s\n", modelServer.Name)
			}
		},
	}
}

func listModelServerVersionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List model server versions",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewAPIClient()
			modelServerVersions, err := client.ListModelServerVersions(modelName, modelServerName)
			if err != nil {
				log.Fatalf("Error listing model server versions: %v", err)
			}
			for _, modelServerVersion := range modelServerVersions {
				fmt.Printf("Name: %s\n", modelServerVersion.Name)
			}
		},
	}
}

func listAcceleratorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List accelerators",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewAPIClient()
			listAcceleratorsResponse, err := client.ListAccelerators(modelName, modelServerName)
			if err != nil {
				log.Fatalf("Error listing accelerators: %v", err)
			}
			fmt.Printf("Min Tpot Milliseconds: %d\n", listAcceleratorsResponse.MinTpotMilliseconds)
			fmt.Printf("Max Tpot Milliseconds: %d\n", listAcceleratorsResponse.MaxTpotMilliseconds)
			fmt.Printf("Min Throughput Tokens Per Second: %d\n", listAcceleratorsResponse.MinThroughputTokensPerSecond)
			fmt.Printf("Max Throughput Tokens Per Second: %d\n", listAcceleratorsResponse.MaxThroughputTokensPerSecond)
			fmt.Printf("Min Ntpot Milliseconds: %d\n", listAcceleratorsResponse.MinNtpotMilliseconds)
			fmt.Printf("Max Ntpot Milliseconds: %d\n", listAcceleratorsResponse.MaxNtpotMilliseconds)
			for _, option := range listAcceleratorsResponse.AcceleratorOptions {
				fmt.Printf("  Accelerator Type: %s\n", option.AcceleratorType)
				fmt.Printf("    Model Name: %s\n", option.ModelAndModelServerInfo.ModelName)
				fmt.Printf("    Model Server Name: %s\n", option.ModelAndModelServerInfo.ModelServerName)
				fmt.Printf("    Model Server Version: %s\n", option.ModelAndModelServerInfo.ModelServerVersion)
				if option.MachineType != "" {
					fmt.Printf("    Machine Type: %s\n", option.MachineType)
				}
				if option.TpuTopology != "" {
					fmt.Printf("    Tpu Topology: %s\n", option.TpuTopology)
				}
				fmt.Printf("    Accelerator Count: %d\n", option.ResourcesUsed.AcceleratorCount)
				fmt.Printf("    Tpot Milliseconds: %d\n", option.PerformanceStats.TpotMilliseconds)
				fmt.Printf("    Queries Per Second: %d\n", option.PerformanceStats.QueriesPerSecond)
				fmt.Printf("    Output Tokens Per Second: %d\n", option.PerformanceStats.OutputTokensPerSecond)
				fmt.Printf("    Ntpot Milliseconds: %d\n", option.PerformanceStats.NtpotMilliseconds)
			}
		},
	}
}

func createManifestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a manifest",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewAPIClient()
			manifestResponse, err := client.CreateManifest(modelName, modelServerName, modelServerVersion, acceleratorType, targetNtpotMilliseconds)
			if err != nil {
				log.Fatalf("Error creating manifest: %v", err)
			}
			for _, k8sManifest := range manifestResponse.K8sManifests {
				fmt.Printf("K8s Manifest Kind: %s\n", k8sManifest.Kind)
				fmt.Printf("K8s Manifest API Version: %s\n", k8sManifest.ApiVersion)
				fmt.Printf("K8s Manifest Content: \n%s\n", k8sManifest.Content)
			}
			for _, comment := range manifestResponse.Comments {
				fmt.Printf("Comment: %s\n", comment)
			}
		},
	}
}

func listModelsAndServersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List models and servers",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewAPIClient()
			modelsAndServers, err := client.ListModelsAndServers()
			if err != nil {
				log.Fatalf("Error listing models and servers: %v", err)
			}
			for _, item := range modelsAndServers {
				fmt.Printf("Model Name: %s, Model Server Name: %s, Create Time: %s, Update Time: %s\n", item.ModelName, item.ModelServerName, item.CreateTime, item.UpdateTime)
			}
		},
	}
}
