# Azure Event Hub Setup Guide for Gommunity Service

This guide will walk you through the process of setting up Azure Event Hub for the Gommunity microservice.

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Create Event Hubs Namespace](#create-event-hubs-namespace)
3. [Create Event Hubs (Topics)](#create-event-hubs-topics)
4. [Get Connection String](#get-connection-string)
5. [Configure the Application](#configure-the-application)
6. [Verify Configuration](#verify-configuration)
7. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- An active Azure subscription
- Azure CLI installed (optional but recommended)
- Access to Azure Portal

---

## Create Event Hubs Namespace

### Using Azure Portal:

1. Navigate to [Azure Portal](https://portal.azure.com)
2. Click **"Create a resource"**
3. Search for **"Event Hubs"**
4. Click **"Create"**

5. Fill in the details:
   - **Subscription**: Select your subscription
   - **Resource Group**: Create new or select existing
   - **Namespace name**: `your-namespace` (e.g., `levelup-journey`)
   - **Location**: Choose closest to your services (e.g., `West US 3`)
   - **Pricing tier**: Standard (supports Kafka protocol)
   - **Throughput Units**: 1 (or more based on your needs)

6. Click **"Review + create"** → **"Create"**

### Using Azure CLI:

```bash
# Create resource group (if needed)
az group create --name your-resource-group --location westus3

# Create Event Hubs namespace
az eventhubs namespace create \
  --name your-namespace \
  --resource-group your-resource-group \
  --location westus3 \
  --sku Standard
```

---

## Create Event Hubs (Topics)

The Gommunity service requires these Event Hubs:

- `community.registration` - For user registration events
- `community.profile.updated` - For profile update events

### Using Azure Portal:

1. Go to your **Event Hubs Namespace**
2. In the left menu, click **"Event Hubs"**
3. Click **"+ Event Hub"**

4. Create the first Event Hub:
   - **Name**: `community.registration`
   - **Partition Count**: 2 (default)
   - **Message Retention**: 1 day (default)
   - Click **"Create"**

5. Repeat for the second Event Hub:
   - **Name**: `community.profile.updated`
   - Click **"Create"**

### Using Azure CLI:

```bash
# Create first Event Hub
az eventhubs eventhub create \
  --name community.registration \
  --namespace-name your-namespace \
  --resource-group your-resource-group \
  --partition-count 2 \
  --message-retention 1

# Create second Event Hub
az eventhubs eventhub create \
  --name community.profile.updated \
  --namespace-name your-namespace \
  --resource-group your-resource-group \
  --partition-count 2 \
  --message-retention 1
```

---

## Get Connection String

### Using Azure Portal:

1. Go to your **Event Hubs Namespace**
2. In the left menu, click **"Shared access policies"**
3. Click on **"RootManageSharedAccessKey"** (or create a new policy)
4. Copy the **"Primary Connection String"**

   It will look like:
   ```
   Endpoint=sb://your-namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=YOUR_KEY_HERE=
   ```

### Using Azure CLI:

```bash
az eventhubs namespace authorization-rule keys list \
  --namespace-name your-namespace \
  --resource-group your-resource-group \
  --name RootManageSharedAccessKey \
  --query primaryConnectionString \
  --output tsv
```

### ⚠️ Security Best Practices:

- **Never commit** the connection string to version control
- Use **Azure Key Vault** in production
- Create **separate policies** with minimum required permissions:
  - Producers need: **Send** permission
  - Consumers need: **Listen** permission
- Rotate keys regularly

---

## Configure the Application

### Step 1: Copy `.env.example` to `.env`

```bash
cp .env.example .env
```

### Step 2: Edit `.env` file

Update the following variables with your Azure Event Hub details:

```bash
# ===================================================
# Kafka Configuration (Azure Event Hub)
# ===================================================

# Your Event Hub namespace endpoint
KAFKA_BOOTSTRAP_SERVERS=your-namespace.servicebus.windows.net:9093

# Security protocol (required for Azure Event Hub)
KAFKA_SECURITY_PROTOCOL=SASL_SSL

# SASL mechanism (Azure Event Hub uses PLAIN)
KAFKA_SASL_MECHANISM=PLAIN

# Username (MUST be exactly this)
KAFKA_SASL_USERNAME=$ConnectionString

# Your complete connection string from Azure Portal
KAFKA_SASL_PASSWORD=Endpoint=sb://your-namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=YOUR_KEY_HERE

# Consumer Group ID
KAFKA_GROUP_ID=gommunity-consumer-group
```

### Step 3: Important Configuration Notes

#### Username
- **MUST** be exactly: `$ConnectionString` (including the `$` symbol)
- Do NOT remove or modify this value

#### Password (Connection String)
- Must be the **complete** connection string from Azure
- Should start with: `Endpoint=sb://`
- Should end with: `SharedAccessKey=YOUR_KEY=`
- **No quotes needed** unless your `.env` parser requires them
- **No spaces** at the beginning or end

#### Bootstrap Servers
- Format: `namespace-name.servicebus.windows.net:9093`
- Port **must be** `9093` (Kafka protocol port for Azure Event Hub)

---

## Verify Configuration

### Step 1: Run the Application

```bash
go run cmd/api/main.go
```

### Step 2: Check the Logs

You should see output similar to:

```
========================================
Configuring Kafka consumer for Azure Event Hub...
Bootstrap Servers: your-namespace.servicebus.windows.net:9093
SASL Mechanism: PLAIN
SASL Username: $ConnectionString
SASL Password: [SET - length: XXX characters]
✓ Password appears to be a valid Azure Event Hub connection string
✓ All Azure Event Hub configuration validations passed
Required Event Hubs (topics) in Azure:
   1. community.registration
   2. community.profile.updated
⚠️  Make sure these Event Hubs exist in your Azure Event Hubs Namespace!
========================================
Kafka consumer created for topics: [community.registration community.profile.updated]
Starting Kafka message consumption...
```

### Step 3: Test with Sample Messages

You can test sending messages using Azure Portal:

1. Go to your Event Hub (e.g., `community.registration`)
2. Click **"Generate data"** or **"Event Hubs Explorer"**
3. Send a test message
4. Check your application logs for the received message

---

## Troubleshooting

### ❌ Error: `SASL Authentication Failed`

**Causes:**
- Incorrect connection string
- Wrong username (must be `$ConnectionString`)
- Expired or invalid access key

**Solutions:**
1. Verify the connection string is complete and correct
2. Ensure username is exactly `$ConnectionString`
3. Regenerate keys in Azure Portal if needed
4. Check for extra spaces or quotes in `.env` file

---

### ❌ Error: `EOF` or Connection Closed

**Causes:**
- Event Hub namespace doesn't exist
- Incorrect bootstrap servers URL
- Firewall or network issues

**Solutions:**
1. Verify namespace name in `KAFKA_BOOTSTRAP_SERVERS`
2. Ensure port is `9093`
3. Check Azure Event Hub is running (not paused)
4. Verify network connectivity to Azure

---

### ❌ Error: `Unknown Topic or Partition`

**Causes:**
- Event Hubs (topics) don't exist in Azure
- Wrong Event Hub names

**Solutions:**
1. Go to Azure Portal → Event Hubs Namespace → Event Hubs
2. Create missing Event Hubs:
   - `community.registration`
   - `community.profile.updated`
3. Ensure names match exactly (case-sensitive)

---

### ❌ Configuration Validation Errors

If you see validation warnings in the logs:

```
⚠️  CONFIGURATION VALIDATION ERRORS:
   1. Password does NOT appear to be a valid Azure Event Hub connection string
```

**Check:**
1. Password starts with `Endpoint=sb://`
2. No extra quotes or spaces
3. Complete connection string (not truncated)
4. Username is `$ConnectionString` (not `ConnectionString`)

---

### 🔍 Enable Debug Logging

For more detailed logs, the application already includes comprehensive error messages:

```
💡 Troubleshooting tips for Azure Event Hub:
   1. Verify the Event Hub (topic) exists in Azure Portal
   2. Check that your Connection String is correct and not expired
   3. Ensure your Shared Access Policy has 'Listen' permission
   4. Confirm the username is exactly: $ConnectionString
   5. Verify your Consumer Group exists (default: $Default)
```

---

## Consumer Groups

By default, the application uses a custom consumer group: `gommunity-consumer-group`

### View Consumer Groups in Azure Portal:

1. Go to Event Hub (e.g., `community.registration`)
2. Click **"Consumer groups"** in the left menu
3. You should see `$Default` and your custom group will be created automatically

### Create Consumer Group Manually (Optional):

```bash
az eventhubs eventhub consumer-group create \
  --eventhub-name community.registration \
  --namespace-name your-namespace \
  --resource-group your-resource-group \
  --name gommunity-consumer-group
```

---

## Security Recommendations

### Production Deployment:

1. **Use Azure Key Vault**
   ```bash
   # Store connection string in Key Vault
   az keyvault secret set \
     --vault-name your-keyvault \
     --name eventhub-connection-string \
     --value "Endpoint=sb://..."
   ```

2. **Create Dedicated Access Policies**
   - Create separate policies for producers and consumers
   - Grant minimum required permissions
   - Use different keys for different environments

3. **Enable Monitoring**
   - Set up Azure Monitor alerts
   - Track failed authentication attempts
   - Monitor throughput and throttling

4. **Rotate Keys Regularly**
   - Azure Event Hub provides primary and secondary keys
   - Rotate without downtime by using secondary key first

---

## Additional Resources

- [Azure Event Hubs Documentation](https://docs.microsoft.com/azure/event-hubs/)
- [Kafka Protocol Support](https://docs.microsoft.com/azure/event-hubs/event-hubs-for-kafka-ecosystem-overview)
- [Best Practices](https://docs.microsoft.com/azure/event-hubs/event-hubs-best-practices)
- [Monitoring and Diagnostics](https://docs.microsoft.com/azure/event-hubs/event-hubs-diagnostic-logs)

---

## Next Steps

After successful setup:

1. ✅ Verify messages are being consumed
2. ✅ Set up monitoring and alerts
3. ✅ Configure auto-scaling if needed
4. ✅ Implement error handling in your message handlers
5. ✅ Set up dead-letter queue for failed messages

---

## Support

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review application logs for detailed error messages
3. Verify Event Hub status in Azure Portal
4. Contact your Azure administrator for access issues