## Instruções

### Para criar a imagem docker

```sh
docker build -t cezaraf/gitlab-approvers-check .
```

### Para executar as imagens com os respectivos parâmetros

```sh
docker run --rm \
	cezaraf/gitlab-approvers-check \
	--host <GITLAB_HOST> \
	--personal-access-token <PERSONAL_ACCESS_TOKEN> \
	--project-id <PROJECT_ID> \
	--merge-request-id <MERGE_REQUEST_ID> \
	--file-check <FILE_CHECK_REGEX> \
	--user-approval-id <USER_APPROVAL_ID> \
	--user-approval-id <ANOTHER_USER_APPROVAL_ID>
```

Sendo:

* GITLAB_HOST = endereço do gitlab com o protocolo http (https://git.tcm.go)
* PERSONAL_ACCESS_TOKEN = token de acesso pessoal gerado no gitlab
* PROJECT_ID = id do projeto no gitlab
* MERGE_REQUEST_ID = id do merge request no gitlab
* FILE_CHECK_REGEX = regex que vai validar se, dado os arquivos do merge request, vai precisar de aprovação dos usuários parametrizados
* USER_APPROVAL_ID = id do usuário no gitlab que precisar aprovar o merge request

### Configuração na pipeline

```yaml
'Checar se é necessário a aprovação por parte da equipe de banco de dados':
  image: docker:stable
  variables:
    DOCKER_HOST: tcp://10.10.110.158:2375
  script: 
    - >
      docker run --rm cezaraf/gitlab-approvers-check
      --host $GITLAB_HOST
      --personal-access-token $PERSONAL_ACCESS_TOKEN
      --project-id $CI_PROJECT_ID
      --merge-request-id $CI_MERGE_REQUEST_IID
      --file-check $SQL_FILE_CHECKER
      --user-approval-id $AISLAN_USER_ID
      --user-approval-id $ROBSON_USER_ID
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event" && $CI_MERGE_REQUEST_TARGET_BRANCH_NAME =~ /^(develop|master|release.+)$/'
```