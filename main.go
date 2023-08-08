package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mkideal/cli"
	"github.com/ortelius/scec-commons/model"
	toml "github.com/pelletier/go-toml"
)

const (
	LicenseFile int = 0
	SwaggerFile int = 1
	ReadmeFile  int = 2
)

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

func getCompToml(derivedAttrs map[string]string) (model.CompAttrs, map[string]string) {
	var extraAttrs map[string]string
	var attrs model.CompAttrs

	for k, v := range derivedAttrs {

		if _, found := os.LookupEnv(strings.ToUpper(k)); !found {
			os.Setenv(strings.ToUpper(k), v)
		}

		switch strings.ToUpper(k) {
		case "BLDDATE":
			attrs.BuildDate = v
		case "BUILDID":
			attrs.BuildID = v
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
		case "GITCOMMIt":
			attrs.GitCommit = v
		case "GITREPO":
			attrs.GitRepo = v
		case "GITTAG":
			attrs.GitTag = v
		case "GITURL":
			attrs.GitURL = v
		case "HIPCHATCHANNEL":
			attrs.HipchatChannel = v
		case "PAGERDUTYBUSINESSURL":
			attrs.PagerdutyBusinessURL = v
		case "PAGERDUTYURL":
			attrs.PagerdutyURL = v
		case "REPOSITORY":
			attrs.Repository = v
			//		case "SERVICEOWNER":
			//			attrs.ServiceOwner = v
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

func run_git(cmdline string) string {
    cmd := exec.Command(cmdline)
    output, _ := cmd.CombinedOutput()
    return string(output)
}

func getDerived() map[string]string {
	var mapping map[string]string

    run_git("git fetch --unshallow")

    mapping["BLDDATE"] = time.Now().UTC().String()
    mapping["SHORT_SHA"] = run_git("git log --oneline -n 1 | cut -d' '  -f1")
    mapping["GIT_COMMIT"] = run_git("git log --oneline -n 1 | cut -d' '  -f1")
    mapping["GIT_VERIFY_COMMIT"] = run_git("git verify-commit " + mapping.get("GIT_COMMIT","") + " 2>&1 | grep -i 'Signature made' | wc -l")
    mapping["GIT_SIGNED_OFF_BY"] = run_git("git log -1 " + mapping.get("GIT_COMMIT","") + " | grep 'Signed-off-by:' | cut -d: -f2 | sed 's/^[ \t]*//;s/[ \t]*$//' | sed 's/&/\\&amp;/g; s/</\\&lt;/g; s/>/\\&gt;/g;'")
    mapping["BUILDNUM"] = run_git('git log --oneline | wc -l | tr -d " "')
    mapping["GIT_REPO"] = run_git("git config --get remote.origin.url | awk -F/ '{print $(NF-1)\"/\"$NF}'| sed 's/.git$//'")
    mapping["GIT_REPO_PROJECT"] = run_git("git config --get remote.origin.url | awk -F/ '{print $(NF-1)}'")
    mapping["GIT_ORG"] = run_git("git config --get remote.origin.url | awk -F/ '{print $NF}'| sed 's/.git$//'")
    mapping["GIT_URL"] = run_git("git config --get remote.origin.url")
    mapping["GIT_BRANCH"] = run_git("git rev-parse --abbrev-ref HEAD")
    mapping["GIT_COMMIT_TIMESTAMP"] = run_git("git log --pretty='format:%%cd' " + mapping.get("SHORT_SHA", "") + " | head -1")
    mapping["GIT_BRANCH_PARENT"] = run_git('git show-branch -a 2>/dev/null | sed "s/].*//" | grep "\\*" | grep -v "$(git rev-parse --abbrev-ref HEAD)" | head -n1 | sed "s/^.*\\[//"')
    mapping["GIT_BRANCH_CREATE_COMMIT"] = run_git("git log --oneline --reverse " + mapping.get("GIT_BRANCH_PARENT", "main") + ".." + mapping.get("GIT_BRANCH") + " | head -1 | awk -F' ' '{print $1}'")
    mapping["GIT_BRANCH_CREATE_TIMESTAMP"] = run_git("git log --pretty='format:%cd' " + mapping.get("GIT_BRANCH_CREATE_COMMIT", "HEAD") + " | head -1")
    mapping["GIT_COMMIT_AUTHORS"] = run_git(
        "git rev-list --remotes --pretty --since='"
        + mapping.get("GIT_BRANCH_CREATE_TIMESTAMP", "")
        + "' --until='"
        + mapping.get("GIT_COMMIT_TIMESTAMP", "")
        + "' | grep -i 'Author:' | grep -v dependabot | awk -F'[:<>]' '{print $3}' | sed 's/^ //' | sed 's/ $//' | sort -u | tr '\n' ',' | sed 's/,$//'"
    )
    mapping["GIT_LINES_TOTAL"] = run_git("wc -l $(git ls-files) | grep total | awk -F' ' '{print $1}'")
    mapping["BASENAME"] = os.path.basename(os.getcwd())
	return mapping
}

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

func gatherEvidence(Userid string, Password string) {

	create_time := time.Now().UTC()
	user := makeUser(Userid)
	license := model.License{Content: gatherFile(LicenseFile)}
	swagger := model.Swagger{Content: gatherFile(SwaggerFile)}
	readme := model.Readme{Content: gatherFile(ReadmeFile)}

	derivedAttrs := getDerived()
	attrs, _ := getCompToml(derivedAttrs)

	var compver model.ComponentVersionDetails

	compver.Attrs = attrs
	compver.CompType = "docker"
	compver.Created = create_time
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
