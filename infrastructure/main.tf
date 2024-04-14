resource "azurerm_user_assigned_identity" "gogoto_identity" {
  resource_group_name = var.resource_group_name
  location            =var.location
  name                = "gogoto-identity"
}

# Conta de armazenamento do Azure 
resource "azurerm_storage_account" "gogoto_sa" {
  name                     = "urlgogotosacc1"
  resource_group_name      = var.resource_group_name
  location                 =var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Delete - Plano de serviço de aplicativo
resource "azurerm_service_plan" "delete-gogoto_asp" {
  name                = "redirect-url-gogoto-asp"
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type             = "Linux"
  sku_name            = "B3"
}

#Create - Plano de serviço de aplicativo
resource "azurerm_service_plan" "create-gogoto_asp" {
  name                = "create-gogoto-asp"
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type             = "Linux"
  sku_name            = "P1v2"
}

# Redirect - Plano de serviço de aplicativo
resource "azurerm_service_plan" "redirect-gogoto_asp" {
  name                = "redirect-gogoto-asp"
  resource_group_name = var.resource_group_name
  location            = var.location
  os_type             = "Linux"
  sku_name            = "P1v2"
}


module "cosmosdb" {
  source = "./modules/cosmosdb"
  location = var.location
  resource_group_name = var.resource_group_name  
  user_assigned_identity_id = azurerm_user_assigned_identity.gogoto_identity.id
  // Passar as variáveis necessárias
}

module "function_app" {
  source = "./modules/function_app"
  location = var.location
  resource_group_name = var.resource_group_name 
  storage_account_name = azurerm_storage_account.gogoto_sa.name
  storage_account_access_key = azurerm_storage_account.gogoto_sa.primary_access_key
  app_service_plan_id_create_gogoto = azurerm_service_plan.create-gogoto_asp.id
  app_service_plan_id_redirect_gogoto = azurerm_service_plan.redirect-gogoto_asp.id
  app_service_plan_id_delete_gogoto = azurerm_service_plan.delete-gogoto_asp.id
  user_assigned_identity_id = azurerm_user_assigned_identity.gogoto_identity.id
}