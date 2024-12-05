# Skript dlya zapuska mikroslervisov i frontenda vyborochno
$OutputEncoding = [System.Text.Encoding]::UTF8
chcp 65001  # Ustanuvlyaem kodovuyu stranitsu UTF-8 dlya konsoli

# Opredelenie putey k mikroslervisam
$services = @(
    @{ Name = "ApplicationService"; Path = "ApplicationService/cmd/server/main.go" },
    @{ Name = "AuthMock"; Path = "AuthMock/cmd/server/main.go" },
    @{ Name = "EventService"; Path = "event-service/cmd/server/main.go" },
    @{ Name = "JuryAssignmentsService"; Path = "Jure-assignments-service/cmd/server/main.go" },
    @{ Name = "OlympiadReact"; Path = "OlympiadReact" }
)

# Funkciya zapuska mikroslervisa
function Start-Service {
    param (
        [string]$Name,
        [string]$Path
    )

    Write-Host "Zapusk servisa: $Name..." -ForegroundColor Green

    if ($Name -eq "OlympiadReact") {
        # Zapusk frontenda (React)
        Push-Location $Path
        Start-Process -FilePath "deno" -ArgumentList "run", "dev"
        Pop-Location
    } else {
        # Zapusk Go-servisa
        Push-Location (Split-Path $Path)
        Start-Process -FilePath "go" -ArgumentList "run", (Split-Path $Path -Leaf)
        Pop-Location
    }

    Write-Host "Servis $Name zapushchen." -ForegroundColor Yellow
}

# Pechat dostupnyh servisov s nomerami
Write-Host "Dostupnye servisy dlya zapuska:" -ForegroundColor Cyan
for ($i = 0; $i -lt $services.Count; $i++) {
    Write-Host "$($i + 1). $($services[$i].Name)"
}

# Zapros vvoda u pol'zovatelya
$selectedIndexes = Read-Host "Vvedite nomera servisov cherez probel (naprimer: 1 3 5)"
$selectedIndexes = $selectedIndexes -split ' ' | ForEach-Object { [int]$_ - 1 }

# Zapusk vybrannykh servisov
foreach ($index in $selectedIndexes) {
    if ($index -ge 0 -and $index -lt $services.Count) {
        $service = $services[$index]
        Start-Service -Name $service.Name -Path $service.Path
    } else {
        Write-Host "Nevernyy nomer: $($index + 1)." -ForegroundColor Red
    }
}
