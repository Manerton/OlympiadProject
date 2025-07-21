chcp 65001 | Out-Null
$OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Пути к сервисам (проверьте, что event-service путь правильный!)
$services = @(
    @{ Name = "ApplicationService"; Path = ".\ApplicationService\cmd\server\main.go" },
    @{ Name = "AuthMock"; Path = ".\AuthMock\cmd\server\main.go" },
    @{ Name = "EventService"; Path = ".\event-service\cmd\server\main.go" },
    @{ Name = "JuryAssignmentsService"; Path = ".\Jure-assignments-service\cmd\server\main.go" },
    @{ Name = "OlympiadReact"; Path = ".\OlympiadReact" }
)

function Start-Service {
    param (
        [string]$Name,
        [string]$Path
    )

    Write-Host "Zapusk $Name..." -ForegroundColor Cyan

    if ($Name -eq "OlympiadReact") {
        # Запуск фронта на Deno
        Push-Location $Path
        Write-Host "Ispol'zuetsya Deno dlya zapuska React" -ForegroundColor Magenta
        Start-Process -NoNewWindow -FilePath "deno" -ArgumentList "run", "-A", "dev.ts"  # или ваш файл запуска
        Pop-Location
    }
    else {
        # Запуск Go-сервиса с отладкой
        $dir = Split-Path $Path -Parent
        Push-Location $dir

        # Вариант 1: Запуск с видимым окном (если падает)
        Start-Process -NoNewWindow -FilePath "cmd.exe" -ArgumentList "/k", "go", "run", (Split-Path $Path -Leaf)

        Pop-Location
    }
}

# Вывод списка сервисов
Write-Host "Spisok servisov:" -ForegroundColor Green
for ($i = 0; $i -lt $services.Count; $i++) {
    Write-Host "$($i + 1). $($services[$i].Name)"
}

# Выбор сервисов
$selected = Read-Host "Vvedite nomera servisov cherez probel (naprimer: 1 3 5)"
$selected -split ' ' | ForEach-Object {
    $index = [int]$_ - 1
    if ($index -ge 0 -and $index -lt $services.Count) {
        $service = $services[$index]
        Start-Service -Name $service.Name -Path $service.Path
    }
    else {
        Write-Host "Oshibka: nevernyy nomer $_" -ForegroundColor Red
    }
}

Write-Host "Zapusk zavershen." -ForegroundColor Green
