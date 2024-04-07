package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type CosmosDBURLItem struct {
	ID             string `json:"id"`
	DestinationUrl string `json:"destination_url"`
	HashID         string `json:"hash_id"`
	PK             string `json:"PK"`
}

type CosmosDBHashItem struct {
	ID          string `json:"id"`
	IsAvailable bool   `json:"is_available"`
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obter o hash da URL
	parts := strings.Split(r.URL.Path, "/")
	hash := parts[len(parts)-1]

	// Consultar o Cosmos DB para obter a URL de destino
	ok, err := DeleteDestinationURL(hash)
	if err != nil {
		log.Printf("Erro ao obter a URL de destino: %s", err.Error())
	}

	if !ok {
		http.Error(w, "URL não encontrada para o hash fornecido", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteDestinationURL(hash string) (bool, error) {
	var cosmosItem CosmosDBURLItem
	//URL do endpoint da conta do Azure Cosmos DB
	var cosmosDBEndpoint = os.Getenv("COSMOS_DB_ENDPOINT")
	var cosmosDBDatabaseID = os.Getenv("COSMOS_DB_DATABASE_ID")
	var cosmosDBContainerID = "urls"

	query := "SELECT TOP 1 c.destination_url FROM c WHERE c.hash_id = @hash_id"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Println(err)
		return false, err
	}

	opt := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@hash_id", Value: hash},
		},
	}

	pk := azcosmos.NewPartitionKeyString("br01")
	client, _ := azcosmos.NewClient(cosmosDBEndpoint, cred, nil)
	container, _ := client.NewContainer(cosmosDBDatabaseID, cosmosDBContainerID)

	queryPager := container.NewQueryItemsPager(query, pk, &opt)

	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.Background())
		if err != nil {
			return false, errors.New("URL not found")
		}

		for _, item := range queryResponse.Items {

			err := json.Unmarshal(item, &cosmosItem)
			if err != nil {
				log.Println(err)
				return false, err
			}

			_, err = container.DeleteItem(context.Background(), pk, cosmosItem.ID, nil)
			if err != nil {
				log.Println(err)
				return false, err
			}
		}
	}
	// makeHashAvailable(hash, client, &pk, &cosmosDBDatabaseID)
	return true, nil
}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/url/{hash}", DeleteHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Println(http.ListenAndServe(listenAddr, nil))
}

// func makeHashAvailable(hash string, client *azcosmos.Client, pk *azcosmos.PartitionKey, databaseId *string) (bool, error) {
// 	var cosmosDBContainerID = "hashes"
// 	containerHash, _ := client.NewContainer(*databaseId, cosmosDBContainerID)
// 	item, err := containerHash.ReadItem(context.Background(), *pk, hash, nil)

// 	if err != nil {
// 		log.Println("Erro ao disponibilizar hash para uso: ", err)
// 		return false, err
// 	}

// 	var cosmosHashItem CosmosDBHashItem
// 	err = json.Unmarshal(item.Value, &cosmosHashItem)
// 	if err != nil {
// 		log.Println(err)
// 		return false, err
// 	}

// 	// Atualizar o atributo is_available para true para disponibilizar o hash para uso
// 	cosmosHashItem.IsAvailable = true

// 	// Converter cosmosHashItem para um azcosmos.Item
// 	itemBytes, err := json.Marshal(cosmosHashItem)

// 	if err != nil {
// 		log.Println("Erro ao converter cosmosHashItem para azcosmos.Item: ", err)
// 		// Retorna o erro mas não interrompe a execução porque a exclusão da URL encurtada já foi realizada
// 		return true, err
// 	}

// 	// Substituir o item no Cosmos DB
// 	_, err = containerHash.ReplaceItem(context.Background(), *pk, hash, itemBytes, nil)
// 	if err != nil {
// 		// Retorna o erro mas não interrompe a execução porque a exclusão da URL encurtada já foi realizada
// 		log.Println("Erro ao disponibilizar o hash para nova utilização: ", err)
// 	}

// 	return true, nil
// }
