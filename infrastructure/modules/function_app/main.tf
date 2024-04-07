
resource "azurerm_linux_function_app" "create-gogoto_fa" {
  name                       = "createShortenedUrl"
  location                   = var.location
  resource_group_name        = var.resource_group_name
  service_plan_id            = var.app_service_plan_id_create_gogoto
  storage_account_name       = var.storage_account_name
  storage_account_access_key = var.storage_account_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [var.user_assigned_identity_id]
  }

  app_settings = {
    FUNCTIONS_WORKER_RUNTIME = "custom"
    WEBSITE_RUN_FROM_PACKAGE = "1"
   }

  site_config {

  }

}

resource "azurerm_linux_function_app" "delete-gogoto_fa" {
  name                       = "deleteShortenedUrl"
  location                   = var.location
  resource_group_name        = var.resource_group_name
  service_plan_id            = var.app_service_plan_id_delete_gogoto
  storage_account_name       = var.storage_account_name
  storage_account_access_key = var.storage_account_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [var.user_assigned_identity_id]
  }

  app_settings = {
    FUNCTIONS_WORKER_RUNTIME = "custom"
    WEBSITE_RUN_FROM_PACKAGE = "1"
  }

  site_config {
  }
}

resource "azurerm_linux_function_app" "redirect-gogoto_fa" {
  name                       = "redirectToDestinationUrl"
  location                   = var.location
  resource_group_name        = var.resource_group_name
  service_plan_id            = var.app_service_plan_id_redirect_gogoto
  storage_account_name       = var.storage_account_name
  storage_account_access_key = var.storage_account_access_key

  identity {
    type         = "UserAssigned"
    identity_ids = [var.user_assigned_identity_id]
  }

  app_settings = {
    FUNCTIONS_WORKER_RUNTIME = "custom"
    WEBSITE_RUN_FROM_PACKAGE = "1"
  }

  site_config {
  }
}
