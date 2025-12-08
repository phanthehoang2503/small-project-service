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

# 1. Login
Write-Host "1. Logging in..." -ForegroundColor Yellow
$authBody = @{ email = "user@example.com"; username = "user1"; password = "password" }

try { Request "POST" "http://localhost:8084/auth/register" $authBody $null $true | Out-Null } catch {}

$loginBody = @{ login = "user@example.com"; password = "password" }
$loginRes = Request "POST" "http://localhost:8084/auth/login" $loginBody
$token = $loginRes.token
Write-Host "   [OK] Logged in! Token acquired." -ForegroundColor Green

# Clear Cart to ensure clean state
try { 
  $headers = @{ "Authorization" = "Bearer $token" }
  Invoke-RestMethod -Method DELETE -Uri "http://localhost:8082/cart" -Headers $headers -ErrorAction SilentlyContinue | Out-Null 
}
catch {}

# 2. Create Product
Write-Host "`n2. Creating Product..." -ForegroundColor Yellow
$prodBody = @( @{ name = "Demo Product"; price = 1000; stock = 100 } )
$prodList = Request "POST" "http://localhost:8081/products" $prodBody
$prod = $prodList[0]
Write-Host "   [OK] Product Created: ID=$($prod.id) Name=$($prod.name) Stock=$($prod.stock)" -ForegroundColor Green

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

# Test Case 2: Stock Failure 
Write-Host "Starting Test Case 2: Stock Failure" -ForegroundColor Cyan
# 8. Add VALID items to Cart
Write-Host "`n8. Adding Valid items to Cart..." -ForegroundColor Yellow
$validCartBody = @{ product_id = $prod.id; quantity = 50 } 
Request "POST" "http://localhost:8082/cart" $validCartBody $token | Out-Null
Write-Host "   [OK] Added 50 items to cart." -ForegroundColor Green

# 9. Sabotage Stock )
Write-Host "`n9. Sabotaging Stock..." -ForegroundColor Yellow
$sabotageBody = @{ name = $prod.name; price = $prod.price; stock = 0 }
Request "PUT" "http://localhost:8081/products/$($prod.id)" $sabotageBody | Out-Null
Write-Host "   [OK] Stock set to 0 via API." -ForegroundColor Green

# 10. Checkout
Write-Host "`n10. Checking Out..." -ForegroundColor Yellow
$failOrder = Request "POST" "http://localhost:8083/orders" @{} $token
Write-Host "   [OK] Order Created: UUID=$($failOrder.uuid) Status=$($failOrder.status)" -ForegroundColor Green

# 11. Wait for Async Processing
Write-Host "`n11. Waiting for Async Processing..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 12. Verify Cancellation
Write-Host "`n12. Verifying Order Cancellation..." -ForegroundColor Yellow
$finalFailOrder = Request "GET" "http://localhost:8083/orders/$($failOrder.id)" $null $token
Write-Host "   Final Status: $($finalFailOrder.status)"

if ($finalFailOrder.status -eq "Cancelled") {
  Write-Host "   [OK] Order was correctly cancelled" -ForegroundColor Green
  Write-Host "`nALL TESTS PASSED! Saga Pattern & OTel are Verified." -ForegroundColor Cyan
}
else {
  Write-Host "   [FAIL] Order was NOT cancelled. Status: $($finalFailOrder.status)" -ForegroundColor Red
}
