package main

import (
	"flag"
	"fmt"
	"gitlab-approvers-check/gitlab"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type idsUserApprovals []string

func (i *idsUserApprovals) String() string {
	return "aff"
}

func (i *idsUserApprovals) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {

	gitlabHostUrl := flag.String("host", "", "Gitlab host")

	personalAccessToken := flag.String("personal-access-token", "", "Personal access token generated from Gitlab")

	projectId := flag.Int("project-id", 0, "Project ID from Gitlab")

	mergeRequestId := flag.Int("merge-request-id", 0, "Merge Request ID from Gitlab")

	fileCheckRegex := flag.String("file-check", "", "Regex to check if target file exists in merge request changes")

	var idsUserApprovals idsUserApprovals
	flag.Var(&idsUserApprovals, "user-approval-id", "Approval required user id")
	flag.Parse()

	if len(*gitlabHostUrl) == 0 {
		fmt.Println("The parameter host is required!")
		os.Exit(1)
	}

	if len(*personalAccessToken) == 0 {
		fmt.Println("The parameter personal-access-token is required!")
		os.Exit(1)
	}

	if *projectId == 0 {
		fmt.Println("The parameter project-id is required!")
		os.Exit(1)
	}

	if *mergeRequestId == 0 {
		fmt.Println("The parameter merge-request-id is required!")
		os.Exit(1)
	}

	if len(*fileCheckRegex) == 0 {
		fmt.Println("The parameter file-check is required!")
		os.Exit(1)
	}

	_, error := regexp.Compile(*fileCheckRegex)

	if error != nil {
		fmt.Println("The parameter file-check needs to be a valid regex!")
		os.Exit(1)
	}

	if len(*&idsUserApprovals) == 0 {
		fmt.Println("The parameter user-approval-id is required!")
		os.Exit(1)
	}

	var gitlabAccessConfig = gitlab.GitlabAccessConfig{
		Host:           *gitlabHostUrl,
		Token:          *personalAccessToken,
		ProjectId:      *projectId,
		MergeRequestId: *mergeRequestId}

	mergeRequestChanges, err := gitlab.GetMergeRequestChanges(gitlabAccessConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(Green("##########################################################################################"))

	fmt.Println(Green(fmt.Sprintf("#### Title: %s", mergeRequestChanges.Title)))

	fmt.Println(Green(fmt.Sprintf("#### Description: %s", mergeRequestChanges.Description)))

	fmt.Println(Green(fmt.Sprintf("#### Author: %s", mergeRequestChanges.Author.Name)))

	fmt.Println(Green("##########################################################################################"))

	fmt.Println()

	mergeRequestApprovals, err := gitlab.GetMergeRequestApprovals(gitlabAccessConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(Yellow("Os arquivos alterados ou criados no merge request foram os seguintes:"))

	for _, change := range mergeRequestChanges.Changes {

		fmt.Println(Yellow(fmt.Sprintf("* %s", change.NewPath)))
	}

	fmt.Println()

	for _, change := range mergeRequestChanges.Changes {

		match, _ := regexp.MatchString(*fileCheckRegex, strings.ToLower(change.NewPath))

		var mandatoryUserApproved = false

		if match {

			for _, approval := range mergeRequestApprovals.ApprovedBy {

				for _, idApprovalRequired := range idsUserApprovals {

					mandatoryUserId, _ := strconv.Atoi(idApprovalRequired)

					if approval.User.Id == mandatoryUserId && mergeRequestChanges.Author.Id != approval.User.Id {

						mandatoryUserApproved = true
					}
				}
			}

			if !mandatoryUserApproved {

				fmt.Println(Red("É necessário a aprovação de um dos seguintes usuários: "))

				for _, idApprovalRequired := range idsUserApprovals {

					mandatoryUserId, _ := strconv.Atoi(idApprovalRequired)

					if mandatoryUserId != mergeRequestChanges.Author.Id {

						user, _ := gitlab.GetUser(gitlabAccessConfig, mandatoryUserId)

						fmt.Println(Red(fmt.Sprintf("* %s", user.Name)))
					}
				}

				os.Exit(1)
			}
		}
	}

	fmt.Println(Green("Tudo certo por aqui! ;)"))
}
