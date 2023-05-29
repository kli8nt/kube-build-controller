package pkg

type Deploy struct {
	AppName string `json:"app_name"`
	Port    int    `json:"app_port"`
}

type Data struct {
	Technology           string `json:"technology"`
	Version              string `json:"version"`
	RepositoryURL        string `json:"repository_url"`
	GithubToken          string `json:"github_token"`
	ApplicationName      string `json:"application_name"`
	RunCommand           string `json:"run_command"`
	BuildCommand         string `json:"build_command"`
	InstallCommand       string `json:"install_command"`
	DependenciesFiles    string `json:"dependencies_files"`
	IsStatic             string `json:"is_static"`
	OutputDirectory      string `json:"output_directory"`
	EnvironmentVariables string `json:"environment_variables"`
	Port                 string `json:"application_port"`
}

type Status struct {
	ApplicationName string `json:"application_name"`
	Status          string `json:"status"`
	Port            string `json:"port"`
}
