# Demo Script for Small Project Microservices
# Usage: ./demo.ps1

$ErrorActionPreference = "Stop"

function Request($method, $url, $body = $null, $token = $null, $quiet = $false) {
  $headers = @{ "Content-Type" = "application/json" }
  if ($token) { $headers["Authorization"] = "Bearer $token" }
    
  $params = @{
    Method  = $method
    Uri     = $url
    Headers = $headers
  }
  if ($body) { $params["Body"] = (ConvertTo-Json -InputObject $body -Depth 10) }

  try {
    $response = Invoke-RestMethod @params
    return $response
  }
  catch {
    if (-not $quiet) {
      Write-Host "Error calling $url" -ForegroundColor Red
      Write-Host $_.Exception.Message -ForegroundColor Red
        
      if ($_.ErrorDetails) {
        Write-Host $_.ErrorDetails.Message -ForegroundColor Red
      }
      elseif ($_.Exception.Response) {
        # Fallback for older PS
        try {
          $reader = New-Object System.IO.StreamReader $_.Exception.Response.GetResponseStream()
          Write-Host $reader.ReadToEnd() -ForegroundColor Red
        }
        catch {
          Write-Host "Could not read response body." -ForegroundColor DarkGray
        }
      }
      exit 1
    }
    else {
      throw $_
    }
  }
}

Write-Host "`nStarting E-Commerce Demo Flow...`n" -ForegroundColor Cyan

# 1. Login
Write-Host "1. Logging in..." -ForegroundColor Yellow
$authBody = @{ email = "user@example.com"; username = "user1"; password = "password" }
# Register first just in case (Quiet mode)
try { Request "POST" "http://localhost:8084/auth/register" $authBody $null $true | Out-Null } catch {}

$loginBody = @{ login = "user@example.com"; password = "password" }
$loginRes = Request "POST" "http://localhost:8084/auth/login" $loginBody
$token = $loginRes.token
Write-Host "   [OK] Logged in! Token acquired." -ForegroundColor Green

# 2. Create Product
Write-Host "`n2. Creating Product..." -ForegroundColor Yellow
$prodBody = @( @{ name = "Demo Product"; price = 1000; stock = 100 } )
$prodList = Request "POST" "http://localhost:8081/products" $prodBody
$prod = $prodList[0]
Write-Host "   [OK] Product Created: ID=$($prod.id) Name=$($prod.name) Stock=$($prod.stock)" -ForegroundColor Green

# 3. Add to Cart
Write-Host "`n3. Adding to Cart..." -ForegroundColor Yellow
$cartBody = @{ product_id = $prod.id; quantity = 2 }
Request "POST" "http://localhost:8082/cart" $cartBody $token | Out-Null
Write-Host "   [OK] Added 2 items to cart." -ForegroundColor Green

# 4. Checkout (Create Order)
Write-Host "`n4. Checking Out..." -ForegroundColor Yellow
$order = Request "POST" "http://localhost:8083/orders" @{} $token
Write-Host "   [OK] Order Created: UUID=$($order.uuid) Total=$($order.total) Status=$($order.status)" -ForegroundColor Green

# 5. Wait for Async Processing
Write-Host "`n5. Waiting for Async Processing (Payment & Stock)..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 6. Verify Order Status
Write-Host "`n6. Verifying Order Status..." -ForegroundColor Yellow
$finalOrder = Request "GET" "http://localhost:8083/orders/$($order.id)" $null $token
Write-Host "   [OK] Final Status: $($finalOrder.status)" -ForegroundColor Green

if ($finalOrder.status -eq "Paid") {
  Write-Host "`nSUCCESS! Order flow completed successfully." -ForegroundColor Cyan
}
else {
  Write-Host "`nWARNING: Order status is $($finalOrder.status). Check logs/traces." -ForegroundColor Red
}

# 7. Verify Stock Deduction
Write-Host "`n7. Verifying Stock Deduction..." -ForegroundColor Yellow
$finalProd = Request "GET" "http://localhost:8081/products/$($prod.id)"
Write-Host "   Initial Stock: 100"
Write-Host "   Bought: 2"
Write-Host "   Current Stock: $($finalProd.stock)"
if ($finalProd.stock -eq 98) {
  Write-Host "   [OK] Stock deducted correctly!" -ForegroundColor Green
}
else {
  Write-Host "   [FAIL] Stock deduction failed!" -ForegroundColor Red
}
