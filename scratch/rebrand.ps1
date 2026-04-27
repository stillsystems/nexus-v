$targets = Get-ChildItem -Path . -Recurse -File -Exclude "*.exe", "*.png", "*.svg", "*.gif", "*.jpg", "*.ico"

foreach ($file in $targets) {
    $content = Get-Content -Path $file.FullName -Raw
    if ($null -eq $content) { continue }
    
    $newContent = $content -replace "SailorOps", "Still Systems" `
                           -replace "sailorops", "stillsystems" `
                           -replace "⚓", "🧱" `
                           -replace "The Anchor Style", "The Foundation Style" `
                           -replace "anchor logo", "foundation logo"
    
    if ($content -ne $newContent) {
        Write-Host "Updating $($file.FullName)"
        Set-Content -Path $file.FullName -Value $newContent -Encoding UTF8
    }
}
