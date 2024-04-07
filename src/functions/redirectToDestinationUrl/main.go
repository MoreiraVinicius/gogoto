package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

type CosmosDBURLItem struct {
	ID             string                `json:"id"`
	DestinationUrl string                `json:"destination_url"`
	HashID         string                `json:"hash_id"`
	PK             azcosmos.PartitionKey `json:"PK"`
}

type CosmosDBHashItem struct {
	ID          string `json:"id"`
	IsAvailable bool   `json:"is_available"`
}

type RequestBody struct {
	DestinationURL string `json:"destination_url"`
}

type ResponseBody struct {
	Hash string `json:"hash"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é Post
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Consultar o Cosmos DB para obter a URL de destino
	hash, err := GetHashFromUrl(&w, r)
	if err != nil {
		log.Printf("Erro ao obter a URL de destino: %s", err.Error())
		return
	}

	log.Printf("Hash criado: %s", hash)

	response := ResponseBody{
		Hash: hash,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao serializar a resposta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func isValidURL(url string) bool {
	// Expressão regular para validar a URL
	regex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`)
	return regex.MatchString(url)
}

func GetHashFromUrl(w *http.ResponseWriter, r *http.Request) (string, error) {
	// Decodificar o corpo da requisição
	var requestBody RequestBody

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(*w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return "", fmt.Errorf("erro ao decodificar o corpo da requisição: %s", err.Error())
	}

	if !isValidURL(requestBody.DestinationURL) {
		httpMsgErro := "URL inválida"
		http.Error(*w, httpMsgErro, http.StatusBadRequest)

		log.Printf("%s: %s", httpMsgErro, err) // Fixed: Removed err.Error()
		return "", err
	}

	hash, err := getAvailableHash()

	if err != nil {
		httpMsgErro := "Erro ao obter um hash disponível"
		http.Error(*w, httpMsgErro, http.StatusInternalServerError)

		log.Printf("%s", err.Error())
		return "", err
	}

	item, err := saveShortenedURL(hash, requestBody.DestinationURL)

	if err != nil {
		httpMsgErro := "Erro ao obter um hash disponível"
		http.Error(*w, httpMsgErro, http.StatusInternalServerError)

		log.Printf("%s", err.Error())
		return "", err
	}

	return item.HashID, nil
}

func getAvailableHash() (string, error) {
	var cosmosItem CosmosDBHashItem
	var cosmosDBHashesContainerID = "hashes"
	var cosmosAccountDBEndpoint = os.Getenv("COSMOS_DB_ACC_ENDPOINT")
	var cosmosDatabaseID = os.Getenv("COSMOS_DB_DATABASE_ID")
	pk := azcosmos.NewPartitionKeyString("br01")

	// Crie um DefaultAzureCredential sem opções, o que fará com que use a identidade gerenciada pelo sistema
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return "", fmt.Errorf("erro ao obter credenciais via 'Azure Identity': %s", err.Error())
	}

	client, err := azcosmos.NewClient(cosmosAccountDBEndpoint, cred, nil)

	if err != nil {
		return "", fmt.Errorf("erro ao obter nova instância do cliente Cosmos via token de acesso do Azure AD: %s", err.Error())
	}

	container, _ := client.NewContainer(cosmosDatabaseID, cosmosDBHashesContainerID)

	query := "SELECT TOP 1 * FROM c WHERE c.is_available = true"

	queryPager := container.NewQueryItemsPager(query, pk, nil)

	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.Background())

		if err != nil {
			return "", fmt.Errorf("erro ao paginar banco de dados: %s", err.Error())
		}

		for _, item := range queryResponse.Items {

			err := json.Unmarshal(item, &cosmosItem)
			if err != nil {
				return "", fmt.Errorf("erro ao fazer Unmarshal: %s", err.Error())
			}

			// Deixa o item indisponível logo para evitar que outro processo o pegue
			cosmosItem.IsAvailable = false
			// Transforma o item em bytes para salvar no banco novamente
			cosmosItemBytes, _ := json.Marshal(cosmosItem)
			_, err = container.ReplaceItem(context.Background(), pk, cosmosItem.ID, cosmosItemBytes, nil)

			if err != nil {
				return "", fmt.Errorf("erro ao da o 'lock' no uso do hash '%s': %s", cosmosItem.ID, err.Error())
			}

			return cosmosItem.ID, nil
		}
	}
	return "", fmt.Errorf("no available hashes")
}

// Salva a URL encurtada no banco de dados e retorna o item salvo ou a mensagem de erro
func saveShortenedURL(hash string, url string) (*CosmosDBURLItem, error) {
	var cosmosItem CosmosDBURLItem
	var cosmosDBEndpoint = os.Getenv("COSMOS_DB_ACC_ENDPOINT")
	var cosmosDBDatabaseID = os.Getenv("COSMOS_DB_DATABASE_ID")
	cosmosDBURLContainerID := "urls"
	partitionKey := "br01"

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, fmt.Errorf("erro ao obter credenciais via 'Azure Identity': %s", err.Error())
	}

	pk := azcosmos.NewPartitionKeyString(partitionKey)
	client, err := azcosmos.NewClient(cosmosDBEndpoint, cred, nil)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar uma nova instância do cliente Cosmos via token de acesso do Azure AD: %s", err.Error())
	}

	container, err := client.NewContainer(cosmosDBDatabaseID, cosmosDBURLContainerID)

	if err != nil {
		return nil, fmt.Errorf("erro ao se conectar ao Container do CosmosDB: %s", err.Error())
	}

	// Monta o item a ser inserido no Cosmos DB
	cosmosItem.ID = uuid.New().String()
	cosmosItem.DestinationUrl = url
	cosmosItem.HashID = hash
	cosmosItem.PK = pk

	itemBytes, err := json.Marshal(cosmosItem)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar o Item para ser salvo no banco de dados: %s", err.Error())
	}

	_, err = container.CreateItem(context.Background(), pk, itemBytes, nil)

	if err != nil {
		return nil, fmt.Errorf("erro ao salvar item no banco de dados: %s", err.Error())
	}

	return &cosmosItem, nil
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/shortener", PostHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Println(http.ListenAndServe(listenAddr, nil))
}
