// Package main - Ortelius CLI for adding Component Versions to the DB from the CI/CD pipeline
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/mkideal/cli"
	"github.com/ortelius/scec-commons/model"
	toml "github.com/pelletier/go-toml"
)

const (
	LicenseFile int = 0 // LicenseFile is used to read the License file
	SwaggerFile int = 1 // SwaggerFile is used to read the Swagger/OpenApi file
	ReadmeFile  int = 2 // ReadmeFile is used to read the Readme file
)

// resolveVars will resolve the ${var} with a value from the component.toml or environment variables
func resolveVars(val string, data map[interface{}]interface{}) string {

	for k, v := range data {
		switch t := v.(type) {
		case map[string]interface{}:
			for a, b := range t {
				val = strings.ReplaceAll(val, "${"+a+"}", b.(string))
			}
		case string:
			val = strings.ReplaceAll(val, "${"+k.(string)+"}", v.(string))
		}
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		val = strings.ReplaceAll(val, "${"+pair[0]+"}", pair[1])
	}
	return val
}

// getCompToml reads the component.toml file and assignes the key/values to the fields in the CompAttrs struct
func getCompToml(derivedAttrs map[string]string) (model.CompAttrs, map[string]string) {
	var attrs model.CompAttrs
	extraAttrs := make(map[string]string, 0)

	for k, v := range derivedAttrs {

		if _, found := os.LookupEnv(strings.ToUpper(k)); !found {
			os.Setenv(strings.ToUpper(k), v)
		}

		switch strings.ToUpper(k) {
		case "BASENAME":
			attrs.Basename = v
		case "BLDDATE":
			attrs.BuildDate = v
		case "BUILDID":
			attrs.BuildID = v
		case "BUILDNUM":
			attrs.BuildNum = v
		case "BUILDURL":
			attrs.BuildURL = v
		case "CHART":
			attrs.Chart = v
		case "CHARTNAMESPACE":
			attrs.ChartNamespace = v
		case "CHARTREPO":
			attrs.ChartRepo = v
		case "CHARTREPOURL":
			attrs.ChartRepoURL = v
		case "CHARTVERSION":
			attrs.ChartVersion = v
		case "DISCORDCHANNEL":
			attrs.DiscordChannel = v
		case "DOCKERREPO":
			attrs.DockerRepo = v
		case "DOCKERSHA":
			attrs.DockerSha = v
		case "DOCKERTAG":
			attrs.DockerTag = v
		case "SHORT_SHA":
			attrs.GitCommit = v
		case "GIT_BRANCH":
			attrs.GitBranch = v
		case "GIT_BRANCH_PARENT":
			attrs.GitBranchParent = v
		case "GIT_BRANCH_CREATE_COMMIT":
			attrs.GitBranchCreateCommit = v
		case "GIT_BRANCH_CREATE_TIMESTAMP":
			attrs.GitBranchCreateTimestamp = v
		case "GIT_COMMIT":
			attrs.GitCommit = v
		case "GITCOMMIT":
			attrs.GitCommit = v
		case "GIT_COMMIT_AUTHORS":
			attrs.GitCommitAuthors = v
		case "GIT_COMMIT_TIMESTAMP":
			attrs.GitCommitTimestamp = v
		case "GIT_COMMITTERS_CNT":
			attrs.GitCommittersCnt = v
		case "GIT_CONTRIB_PERCENTAGE":
			attrs.GitContribPercentage = v
		case "GIT_LINES_ADDED":
			attrs.GitLinesAdded = v
		case "GIT_LINES_DELETED":
			attrs.GitLinesDeleted = v
		case "GIT_LINES_TOTAL":
			attrs.GitLinesTotal = v
		case "GIT_ORG":
			attrs.GitOrg = v
		case "GIT_PREVIOUS_COMPONENT_COMMIT":
			attrs.GitPrevCompCommit = v
		case "GIT_REPO_PROJECT":
			attrs.GitRepoProject = v
		case "GIT_REPO":
			attrs.GitRepo = v
		case "GITREPO":
			attrs.GitRepo = v
		case "GIT_TAG":
			attrs.GitTag = v
		case "GITTAG":
			attrs.GitTag = v
		case "GIT_TOTAL_COMMITTERS_CNT  ":
			attrs.GitTotalCommittersCnt = v
		case "GIT_URL":
			attrs.GitURL = v
		case "GITURL":
			attrs.GitURL = v
		case "GIT_VERIFY_COMMIT":
			attrs.GitVerifyCommit = v
		case "GIT_SIGNED_OFF_BY":
			attrs.GitSignedOffBy = v
		case "HIPCHATCHANNEL":
			attrs.HipchatChannel = v
		case "PAGERDUTYBUSINESSURL":
			attrs.PagerdutyBusinessURL = v
		case "PAGERDUTYURL":
			attrs.PagerdutyURL = v
		case "REPOSITORY":
			attrs.Repository = v
		case "SERVICEOWNER":
			attrs.ServiceOwner = makeUser(v)
		case "SLACKCHANNEL":
			attrs.SlackChannel = v

		}
	}

	f, err := os.ReadFile("component.toml")

	if err != nil {
		log.Println(err)
		return attrs, extraAttrs
	}

	var data map[interface{}]interface{}

	err = toml.Unmarshal(f, &data)
	log.Println(data)

	if err != nil {
		log.Println(err)
		return attrs, extraAttrs
	}

	for k, v := range data {
		switch t := v.(type) {
		case map[string]interface{}:
			for a, b := range t {
				extraAttrs[a] = resolveVars(b.(string), data)
			}
		case string:
			switch strings.ToUpper(k.(string)) {
			case "BLDDATE":
				attrs.BuildDate = resolveVars(v.(string), data)
			case "BUILDID":
				attrs.BuildID = resolveVars(v.(string), data)
			case "BUILDURL":
				attrs.BuildURL = resolveVars(v.(string), data)
			case "CHART":
				attrs.Chart = resolveVars(v.(string), data)
			case "CHARTNAMESPACE":
				attrs.ChartNamespace = resolveVars(v.(string), data)
			case "CHARTREPO":
				attrs.ChartRepo = resolveVars(v.(string), data)
			case "CHARTREPOURL":
				attrs.ChartRepoURL = resolveVars(v.(string), data)
			case "CHARTVERSION":
				attrs.ChartVersion = resolveVars(v.(string), data)
			case "DISCORDCHANNEL":
				attrs.DiscordChannel = resolveVars(v.(string), data)
			case "DOCKERREPO":
				attrs.DockerRepo = resolveVars(v.(string), data)
			case "DOCKERSHA":
				attrs.DockerSha = resolveVars(v.(string), data)
			case "DOCKERTAG":
				attrs.DockerTag = resolveVars(v.(string), data)
			case "GITCOMMIt":
				attrs.GitCommit = resolveVars(v.(string), data)
			case "GITREPO":
				attrs.GitRepo = resolveVars(v.(string), data)
			case "GITTAG":
				attrs.GitTag = resolveVars(v.(string), data)
			case "GITURL":
				attrs.GitURL = resolveVars(v.(string), data)
			case "HIPCHATCHANNEL":
				attrs.HipchatChannel = resolveVars(v.(string), data)
			case "PAGERDUTYBUSINESSURL":
				attrs.PagerdutyBusinessURL = resolveVars(v.(string), data)
			case "PAGERDUTYURL":
				attrs.PagerdutyURL = resolveVars(v.(string), data)
			case "REPOSITORY":
				attrs.Repository = resolveVars(v.(string), data)
				//		case "SERVICEOWNER":
				//			attrs.ServiceOwner = resolveVars(v.(string), data)
			case "SLACKCHANNEL":
				attrs.SlackChannel = resolveVars(v.(string), data)
			}
		}
	}
	return attrs, extraAttrs
}

// gatherFile finds and reads the license, swagger or readme into a string array
func gatherFile(filetype int) []string {

	lines := make([]string, 0)
	filename := ""

	if filetype == LicenseFile {
		if _, err := os.Stat("LICENSE"); err == nil {
			filename = "LICENSE"
		} else if _, err := os.Stat("LICENSE.md"); err == nil {
			filename = "LICENSE.md"
		} else if _, err := os.Stat("license"); err == nil {
			filename = "license"
		} else if _, err := os.Stat("license.md"); err == nil {
			filename = "license.md"
		}
	} else if filetype == SwaggerFile {
		if _, err := os.Stat("swagger.yaml"); err == nil {
			filename = "swagger.yaml"
		} else if _, err := os.Stat("swagger.yml"); err == nil {
			filename = "swagger.yml"
		} else if _, err := os.Stat("swagger.json"); err == nil {
			filename = "swagger.json"
		} else if _, err := os.Stat("openapi.json"); err == nil {
			filename = "openapi.json"
		} else if _, err := os.Stat("openapi.yaml"); err == nil {
			filename = "openapi.yaml"
		} else if _, err := os.Stat("openapi.yml"); err == nil {
			filename = "openapi.yml"
		}
	} else if filetype == ReadmeFile {
		if _, err := os.Stat("README"); err == nil {
			filename = "README"
		} else if _, err := os.Stat("README.md"); err == nil {
			filename = "README.md"
		} else if _, err := os.Stat("readme"); err == nil {
			filename = "readme"
		} else if _, err := os.Stat("readme.md"); err == nil {
			filename = "readme.md"
		}
	}

	if len(filename) > 0 {
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Println(err)
			return lines
		}

		lines = strings.Split(string(data), "\n")
		return lines
	}
	return lines
}

// runGit executes a shell command and returns the output as a string
func runGit(cmdline string) string {
	cmd := exec.Command("sh", "-c", cmdline)
	output, _ := cmd.CombinedOutput()

	return strings.TrimSuffix(string(output), "\n")
}

// getWithDefault is a helper function for finding a key in a map and return a default value if the key is not found
func getWithDefault(m map[string]string, key string, defaultStr string) string {
	if x, found := m[key]; found {
		return x
	}
	return defaultStr
}

// getDerived will run commands in the current working directory to derive data mainly from git
func getDerived() map[string]string {
	mapping := make(map[string]string, 0)

	runGit("git fetch --unshallow 2>/dev/null")

	mapping["BLDDATE"] = time.Now().UTC().String()
	mapping["SHORT_SHA"] = runGit("git log --oneline -n 1 | cut -d' '  -f1")
	mapping["GIT_COMMIT"] = runGit("git log --oneline -n 1 | cut -d' '  -f1")
	mapping["GIT_VERIFY_COMMIT"] = runGit("git verify-commit " + getWithDefault(mapping, "GIT_COMMIT", "") + " 2>&1 | grep -i 'Signature made' | wc -l")
	mapping["GIT_SIGNED_OFF_BY"] = runGit("git log -1 " + getWithDefault(mapping, "GIT_COMMIT", "") + " | grep 'Signed-off-by:' | cut -d: -f2 | sed 's/^[ \t]*//;s/[ \t]*$//' | sed 's/&/\\&amp;/g; s/</\\&lt;/g; s/>/\\&gt;/g;'")
	mapping["BUILDNUM"] = runGit("git log --oneline | wc -l | tr -d \" \"")
	mapping["GIT_REPO"] = runGit("git config --get remote.origin.url | awk -F/ '{print $(NF-1)\"/\"$NF}'| sed 's/.git$//'")
	mapping["GIT_REPO_PROJECT"] = runGit("git config --get remote.origin.url | awk -F/ '{print $(NF-1)}'")
	mapping["GIT_ORG"] = runGit("git config --get remote.origin.url | awk -F/ '{print $NF}'| sed 's/.git$//'")
	mapping["GIT_URL"] = runGit("git config --get remote.origin.url")
	mapping["GIT_BRANCH"] = runGit("git rev-parse --abbrev-ref HEAD")
	mapping["GIT_COMMIT_TIMESTAMP"] = runGit("git log --pretty='format:%cd' " + getWithDefault(mapping, "SHORT_SHA", "") + " | head -1")
	mapping["GIT_BRANCH_PARENT"] = runGit("git show-branch -a 2>/dev/null | sed \"s/].*//\" | grep \"\\*\" | grep -v \"$(git rev-parse --abbrev-ref HEAD)\" | head -n1 | sed \"s/^.*\\[//\"")
	mapping["GIT_BRANCH_CREATE_COMMIT"] = runGit("git log --oneline --reverse " + getWithDefault(mapping, "GIT_BRANCH_PARENT", "main") + ".." + getWithDefault(mapping, "GIT_BRANCH", "main") + " | head -1 | awk -F' ' '{print $1}'")
	mapping["GIT_BRANCH_CREATE_TIMESTAMP"] = runGit("git log --pretty='format:%cd' " + getWithDefault(mapping, "GIT_BRANCH_CREATE_COMMIT", "HEAD") + " | head -1")
	mapping["GIT_COMMIT_AUTHORS"] = runGit("git rev-list --remotes --pretty --since='" + getWithDefault(mapping, "GIT_BRANCH_CREATE_TIMESTAMP", "") + "' --until='" + getWithDefault(mapping, "GIT_COMMIT_TIMESTAMP", "") + "' | grep -i 'Author:' | grep -v dependabot | awk -F'[:<>]' '{print $3}' | sed 's/^ //' | sed 's/ $//' | sort -u | tr '\n' ',' | sed 's/,$//'")
	mapping["GIT_LINES_TOTAL"] = runGit("wc -l $(git ls-files) | grep total | awk -F' ' '{print $1}'")
	cwd, _ := os.Getwd()
	mapping["BASENAME"] = path.Base(cwd)
	return mapping
}

// makeUser takes a string and creates a User struct.  Handles setting the domain if the string contains dots.
func makeUser(name string) model.User {
	var user model.User

	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		name = parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		user.Domain = model.Domain{Name: strings.Join(parts, ".")}
	}
	user.Name = name

	return user
}

// gatherEvidence collects data from the component.toml and git repo for the component version
func gatherEvidence(Userid string, Password string) {

	fmt.Println(Password)
	createTime := time.Now().UTC()
	user := makeUser(Userid)
	license := model.License{Content: gatherFile(LicenseFile)}
	swagger := model.Swagger{Content: gatherFile(SwaggerFile)}
	readme := model.Readme{Content: gatherFile(ReadmeFile)}

	derivedAttrs := getDerived()
	attrs, _ := getCompToml(derivedAttrs)

	var compver model.ComponentVersionDetails

	compver.Attrs = attrs
	compver.CompType = "docker"
	compver.Created = createTime
	compver.Creator = user
	//	compver.Domain = domain
	compver.License = license
	//	compver.Name = compname
	//	compver.Owner = Userid
	//	compver.Packages = sbom
	//	compver.ParentKey = parent
	compver.Readme = readme
	compver.Swagger = swagger

	URL := "https://cat-fact.herokuapp.com/facts"
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}

// main is the entrypoint for the CLI.  Takes --user and --pass parameters
func main() {
	type argT struct {
		cli.Helper
		Userid   string `cli:"*user" usage:"User id (required)"`
		Password string `cli:"*pass" usage:"User password (required)"`
	}

	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		gatherEvidence(argv.Userid, argv.Password)
		return nil
	}))
}
