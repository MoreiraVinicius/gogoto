

# Conta do CosmosDB com consistência eventual para armazenar os dados da url encurtada
resource "azurerm_cosmosdb_account" "gogoto_cosmos_account_eventual" {
  name                      = "gogoto-cdb-acc-eventual"
  location                  = var.location
  resource_group_name       = var.resource_group_name
  offer_type                = "Standard"
  kind                      = "GlobalDocumentDB"
  enable_automatic_failover = true

  consistency_policy {
    consistency_level = "Eventual"
  }

  geo_location {
    location          = var.location
    failover_priority = 0
  }

  identity {
    type = "UserAssigned"
    identity_ids = [var.user_assigned_identity_id]
  }
}

# Conta do CosmosDB com consistência forte para armazenar os dados de hash sem colisão
resource "azurerm_cosmosdb_account" "gogoto_cosmos_account_session" {
  name                      = "gogoto-cdb-acc-session"
  location                  = var.location
  resource_group_name       = var.resource_group_name
  offer_type                = "Standard"
  kind                      = "GlobalDocumentDB"
  enable_automatic_failover = true

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = var.location
    failover_priority = 0
  }
}