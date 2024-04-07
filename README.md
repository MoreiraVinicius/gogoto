# redirect-url-shortener
Azure Function + Golang 

## Passo a Passo

Autenticação
```shell
az login
```
Defina a Assinatura Azure para a qual criar o Principal de Serviço
```shell
az account set -s <subscription-id>
```
This will yield something like:

```json
{
  "appId": "<servicePrincipalId>",
  "displayName": "<name>",
  "name": "<name>",
  "password": "<password>",
  "tenant": "<tenantId>"
}
```
Set environment variables with values from above service principal

Bash
```bash
export AZURE_SUBSCRIPTION_ID='<subscriptionId>'
export AZURE_TENANT_ID='<tenantId>'
export AZURE_CLIENT_ID='<servicePrincipalId>'
export AZURE_CLIENT_SECRET='<password>'
```

Inicialize um projeto de Azure Functions para manipuladores personalizados.
```shell
func init --worker-runtime custom
```
Publique na Azure
```shell
func azure functionapp publish FUNCTION_APP_NAME
```
#### Usando o Azure Service Principal para RBAC como Credencial de Implantação
Usando o Azure Service Principal para RBAC como Credencial de Implantação 

NOTA: Se você deseja implantar no plano de consumo Linux e seu aplicativo __contém um arquivo executável__ (como no caso do uso do Golang), você precisa usar este método para manter a permissão de execução.

Siga estas etapas para configurar seu fluxo de trabalho para usar um Azure Service Principal para RBAC e adicioná-los como um Segredo do GitHub em seu repositório.

Baixe o Azure CLI daqui, execute o comando az login para entrar com suas credenciais do Azure. Execute o comando Azure CLI
```shell
az ad sp create-for-rbac --name "myServicePrincipalApp" --role contributor \
                        --scopes /subscriptions/{subscription-id}/resourceGroups/{resource-group}/providers/Microsoft.Web/sites/{app-name} \
                        --sdk-auth

  # Substitua {subscription-id}, {resource-group} e {app-name} pelos nomes  de sua assinatura, grupo de recursos e aplicativo de função Azure. O comando deve gerar um objeto JSON semelhante a este:

  {
    "clientId": "<GUID>",
    "clientSecret": "<GUID>",
    "subscriptionId": "<GUID>",
    "tenantId": "<GUID>",
    (...)
  }
```
2 - Copie e cole a resposta JSON acima do Azure CLI para o seu Repositório GitHub > Configurações > Segredos > Adicionar um novo segredo > CREDENCIAIS_AZURE_RBAC

3 - Use o modelo RBAC de Aplicativo de Função DotNet do Windows como referência para construir seu fluxo de trabalho no diretório .github/workflows/. Certifique-se de usar a ação azure/login e de não usar o parâmetro publish-profile

4 - Altere os valores das variáveis na seção env: de acordo com o seu aplicativo de função.

5 - Envie e faça o push do seu projeto para o repositório do GitHub, você deverá ver um novo fluxo de trabalho do GitHub iniciado na guia Ações.


# Motivação

Este arquivo de configuração parece ser usado para definir comportamentos específicos em um ambiente Azure Functions, incluindo como os logs são tratados, quais extensões são usadas e como os manipuladores personalizados (custom handlers) são configurados. Vamos detalhar cada parâmetro:

- `"version": "2.0"`: Especifica a versão do schema de configuração usado pelo arquivo. Isso ajuda a garantir que o arquivo seja interpretado corretamente pela plataforma Azure Functions.

- `"logging"`: Configurações relacionadas ao registro de logs.
  - `"applicationInsights"`: Define configurações específicas para o Application Insights, uma ferramenta de monitoramento e análise.
    - `"samplingSettings"`: Configurações de amostragem para reduzir o volume de dados enviados ao Application Insights.
      - `"isEnabled": true`: Ativa a amostragem de dados.
      - `"excludedTypes": "Request"`: Especifica que os dados de solicitações HTTP não serão amostrados, ou seja, todos os dados de solicitação serão enviados ao Application Insights.

- `"extensionBundle"`: Define um conjunto de extensões que devem ser usadas pela função.
  - `"id": "Microsoft.Azure.Functions.ExtensionBundle"`: Identifica o pacote de extensões a ser usado.
  - `"version": "[3.*, 4.0.0)"`: Especifica a versão do pacote de extensões, neste caso, qualquer versão maior ou igual a 3.0 e menor que 4.0.0.

- `"customHandler"`: Configurações para um manipulador personalizado, que permite executar funções em linguagens ou runtimes não diretamente suportados pelo Azure Functions.
  - `"description"`: Descreve como o manipulador personalizado deve ser executado.
    - `"defaultExecutablePath": "server"`: Caminho para o executável padrão que inicia o manipulador personalizado.
    - `"workingDirectory": ""`: Diretório de trabalho para o manipulador personalizado. Um valor vazio significa que será usado o diretório padrão.
    - `"arguments": []`: Argumentos adicionais que devem ser passados para o executável.
  - `"enableForwardingHttpRequest": true`: Quando verdadeiro, permite que as solicitações HTTP sejam encaminhadas diretamente para o manipulador personalizado, facilitando a criação de APIs HTTP.

Cada um desses parâmetros permite uma personalização detalhada de como as funções Azure são executadas, monitoradas e estendidas, proporcionando flexibilidade para atender a diferentes necessidades de desenvolvimento.

### Configurando acesso ao CosmosDB localmente
 Caso queira rodar a função localmente, o Azure Identity não pode usar a identidade gerenciada pelo sistema, porque ela não está disponível. Nesse caso, o Azure Identity tentará usar outras estratégias de autenticação, como o Azure CLI ou uma conta de serviço.

Se você estiver autenticado no Azure CLI, o Azure Identity pode usar as credenciais do Azure CLI para autenticar no Cosmos DB. Você pode se autenticar no Azure CLI com o seguinte comando:
```bash
az login
```
Se você não estiver autenticado no Azure CLI ou se preferir usar uma conta de serviço, você pode criar uma conta de serviço no portal do Azure e baixar o arquivo de chave JSON. Em seguida, você pode definir a variável de ambiente AZURE_CLIENT_SECRET para o caminho do arquivo de chave JSON. O Azure Identity usará a conta de serviço para autenticar no Cosmos DB.

Por favor, note que a conta de serviço precisa ter as permissões necessárias para acessar o Cosmos DB. Você pode conceder as permissões no portal do Azure, na seção "Access control (IAM)" do Cosmos DB.

#### E na nuvem, como fica a gestão de identidade?
O sistema é composto de multiplas funções Azure com acesso aos mesmos recursos,  resolvi optar em compartilhar a mesma identidade de usuario gerenciada entre elas (User-Assigned Managed Identity). Isso proporcionaria maior flexibilidade e reutilização da identidade gerenciada em múltiplos recursos.

Ao usar uma identidade gerenciada do usuário, você pode criar uma única identidade gerenciada e associá-la à todas as funções Azure. Isso simplifica a gestão de identidades e centraliza a administração das permissões concedidas a essa identidade. Além disso, se houver a necessidade de adicionar mais funções no futuro que também precisem da mesma identidade, Eu posso facilmente associá-las à identidade gerenciada existente.

Por outro lado, se Eu optasse por usar identidades gerenciadas do sistema separadas para cada função, ia acabar gerando mais complexidade na gestão das identidades e das permissões, uma vez que cada função terá sua própria identidade gerenciada. Isso pode tornar a administração e a manutenção das permissões mais trabalhosas, especialmente porque funções compartilham o mesmo conjunto de permissões.
