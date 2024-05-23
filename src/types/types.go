package types

type LanguageDesc_t struct {
	Language string   `json:"language"`
	Compiler string   `json:"compiler"`
	Flags    []string `json:"flags"`
}

type Project_t struct {
	Name            string         `json:"name"`
	Version         string         `json:"version"`
	Description     string         `json:"descprition"`
	PackageUrl      string         `json:"package-url"`
	BuildExecutable bool           `json:"build-executable"`
	Author          string         `json:"author"`
	LanguageDesc    LanguageDesc_t `json:"lang-descprition"`
}

type Src_t struct {
	Main        string   `json:"main"`
	Includes    []string `json:"includes"`
	Libs        []string `json:"libs"`
	LinkerFlags []string `json:"linker-flags"`
}

type Dependencies_t struct {
	PackageUrl string `json:"package-url"`
	Version    string `json:"version"`
}

type Config_t struct {
	Project      Project_t        `json:"project"`
	Src          Src_t            `json:"src"`
	Dependencies []Dependencies_t `json:"deps"`
}
