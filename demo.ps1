$ErrorActionPreference = "Stop"

$Global:TotalTests = 0
$Global:PassedTests = 0
$Global:FailedTests = 0

function Write-Log($type, $message) {
  switch ($type) {
    "INFO" { Write-Host " [INFO] $message" -ForegroundColor Gray }
    "STEP" { Write-Host "`n [STEP] $message" -ForegroundColor Yellow }
    "PASS" { Write-Host " [PASS] $message" -ForegroundColor Green }
    "FAIL" { Write-Host " [FAIL] $message" -ForegroundColor Red }
    "WARN" { Write-Host " [WARN] $message" -ForegroundColor Magenta }
    "HEADER" { Write-Host "`n=== $message ===" -ForegroundColor Cyan }
  }
}

function Assert-Equal($actual, $expected, $message) {
  $Global:TotalTests++
  if ($actual -eq $expected) {
    Write-Log "PASS" "$message"
    $Global:PassedTests++
  }
  else {
    Write-Log "FAIL" "$message (Expected: '$expected', Got: '$actual')"
    $Global:FailedTests++
  }
}

function Request($method, $url, $body = $null, $token = $null) {
  $headers = @{ "Content-Type" = "application/json" }
  if ($token) { $headers["Authorization"] = "Bearer $token" }
    
  $params = @{ Method = $method; Uri = $url; Headers = $headers }
  if ($body) { $params["Body"] = (ConvertTo-Json $body -Depth 10) }

  try {
    return Invoke-RestMethod @params
  }
  catch {
    throw "API Call Failed ($method $url): $($_.Exception.Message)"
  }
}

function Initialize-Environment {
  Write-Log "HEADER" "Environment Setup"
    
  Write-Log "INFO" "Authenticating..."
  try { Request "POST" "http://localhost:8084/auth/register" @{ email = "test@example.com"; username = "tester"; password = "password" } | Out-Null } catch {}
  $token = (Request "POST" "http://localhost:8084/auth/login" @{ login = "test@example.com"; password = "password" }).token
  Write-Log "INFO" "Token acquired."

  Write-Log "INFO" "Clearing cart..."
  try { Invoke-RestMethod -Method DELETE -Uri "http://localhost:8082/cart" -Headers @{ "Authorization" = "Bearer $token" } -ErrorAction SilentlyContinue | Out-Null } catch {}
    
  return $token
}

function Test-HappyPath($token) {
  Write-Log "HEADER" "Test 1: Successful Order"

  Write-Log "STEP" "Creating Product & Adding to Cart"
  $prod = (Request "POST" "http://localhost:8081/products" @( @{ name = "Valid Product"; price = 1000; stock = 100 } ))[0]
  Request "POST" "http://localhost:8082/cart" @{ product_id = $prod.id; quantity = 2 } $token | Out-Null
    
  Write-Log "STEP" "Checking Out"
  $order = Request "POST" "http://localhost:8083/orders" @{} $token
  Write-Log "INFO" "Order Created: $($order.uuid)"

  Write-Log "STEP" "Waiting for Saga..."
  Start-Sleep -Seconds 5
    
  $finalOrder = Request "GET" "http://localhost:8083/orders/$($order.id)" $null $token
  Assert-Equal $finalOrder.status "Paid" "Order status should be Paid"

  $finalProd = Request "GET" "http://localhost:8081/products/$($prod.id)"
  Assert-Equal $finalProd.stock 98 "Stock should decrease by 2"
}

function Test-StockFailure($token) {
  Write-Log "HEADER" "Test 2: Stock Failure"

  Write-Log "STEP" "Setting up Stock Failure Scenario"
  $prod = (Request "POST" "http://localhost:8081/products" @( @{ name = "Fail Product"; price = 1000; stock = 100 } ))[0]
    
  Request "POST" "http://localhost:8082/cart" @{ product_id = $prod.id; quantity = 50 } $token | Out-Null
    
  Request "PUT" "http://localhost:8081/products/$($prod.id)" @{ name = $prod.name; price = $prod.price; stock = 0 } | Out-Null
  Write-Log "INFO" "Stock sabotaged to 0."

  Write-Log "STEP" "Checking Out"
  $order = Request "POST" "http://localhost:8083/orders" @{} $token
    
  Write-Log "STEP" "Waiting for Saga..."
  Start-Sleep -Seconds 5
    
  $finalOrder = Request "GET" "http://localhost:8083/orders/$($order.id)" $null $token
  Assert-Equal $finalOrder.status "Cancelled" "Order status should be Cancelled"
  Assert-Equal $finalOrder.status "Cancelled" "Order status should be Cancelled"
}

function Test-PaymentFailure($token) {
  Write-Log "HEADER" "Test 2b: Payment Failure (Chaos)"

  Write-Log "STEP" "Creating Cursed Product (Price 6666)"
  $prod = (Request "POST" "http://localhost:8081/products" @( @{ name = "Cursed Product"; price = 6666; stock = 10 } ))[0]
    
  Request "POST" "http://localhost:8082/cart" @{ product_id = $prod.id; quantity = 1 } $token | Out-Null
    
  Write-Log "STEP" "Checking Out (Should Fail Payment)"
  $order = Request "POST" "http://localhost:8083/orders" @{} $token
  Write-Log "INFO" "Order Created: $($order.uuid)"
    
  Write-Log "STEP" "Waiting for Saga..."
  Start-Sleep -Seconds 5
    
  $finalOrder = Request "GET" "http://localhost:8083/orders/$($order.id)" $null $token
  Assert-Equal $finalOrder.status "Cancelled" "Order status should be Cancelled"
}

function Test-Rollback($token) {
  Write-Log "HEADER" "Test 3: Rollback"

  Write-Log "STEP" "Creating Mixed Cart (1 Valid, 1 Sabotaged)"
  $prodA = (Request "POST" "http://localhost:8081/products" @( @{ name = "Safe Product"; price = 100; stock = 50 } ))[0]
  $prodB = (Request "POST" "http://localhost:8081/products" @( @{ name = "Rollback Target"; price = 100; stock = 50 } ))[0]

  try { Invoke-RestMethod -Method DELETE -Uri "http://localhost:8082/cart" -Headers @{ "Authorization" = "Bearer $token" } -ErrorAction SilentlyContinue | Out-Null } catch {}

  Request "POST" "http://localhost:8082/cart" @{ product_id = $prodA.id; quantity = 1 } $token | Out-Null
  Request "POST" "http://localhost:8082/cart" @{ product_id = $prodB.id; quantity = 1 } $token | Out-Null

  Request "PUT" "http://localhost:8081/products/$($prodB.id)" @{ name = $prodB.name; price = $prodB.price; stock = 0 } | Out-Null
  Write-Log "INFO" "Sabotaged Item B (Stock=0)."

  Write-Log "STEP" "Checking Out"
  $order = Request "POST" "http://localhost:8083/orders" @{} $token

  Write-Log "STEP" "Waiting for Saga..."
  Start-Sleep -Seconds 5
    
  $finalOrder = Request "GET" "http://localhost:8083/orders/$($order.id)" $null $token
  Assert-Equal $finalOrder.status "Cancelled" "Order should be Cancelled"

  $checkA = Request "GET" "http://localhost:8081/products/$($prodA.id)"
  Assert-Equal $checkA.stock 50 "Safe Product stock should remain 50 (Rollback)"

  $checkB = Request "GET" "http://localhost:8081/products/$($prodB.id)"
  Assert-Equal $checkB.stock 0 "Sabotaged Product stock should remain 0"
}

try {
  $token = Initialize-Environment
    
  Test-HappyPath $token
  Test-StockFailure $token
  Test-PaymentFailure $token
  Test-Rollback $token

  Write-Log "HEADER" "Test Summary"
  Write-Host " Total Tests : $Global:TotalTests"
  Write-Host " Passed      : $Global:PassedTests" -ForegroundColor Green
  Write-Host " Failed      : $Global:FailedTests" -ForegroundColor Red

  if ($Global:FailedTests -gt 0) { exit 1 }
}
catch {
  Write-Log "FAIL" "Critical Script Error: $_"
  exit 1
}
