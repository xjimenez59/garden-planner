# deploy-nas.ps1 — Déploiement de Garden Planner sur le NAS QNAP
#
# Prérequis (une seule fois) : configurer la clé SSH pour éviter les saisies de mot de passe
#   ssh-keygen -t ed25519 -C "garden-planner-deploy"   # si pas encore de clé SSH
#   type "$env:USERPROFILE\.ssh\id_ed25519.pub" | ssh xavier@192.168.1.3 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
#
# Note : après le premier déploiement via ce script, Container Station marque l'application
# comme "gérée en externe" — l'édition du compose depuis l'UI est désactivée, mais les
# conteneurs restent visibles et gérables (start/stop/restart) depuis Container Station.

param(
    [string]$NasHost    = "192.168.1.3",
    [string]$NasUser    = "xavier",
    [string]$NasAppName = "garden-planner",  # nom de l'application dans Container Station
    [switch]$NoBuild                         # sauter le build (utiliser les images déjà construites)
)

$NasSsh      = "${NasUser}@${NasHost}"
$NasAppDir   = "/share/local/container-station-data/application/${NasAppName}"
$NasDocker   = "/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker"
$Images      = @("garden-planner-api", "garden-planner-meteo", "garden-planner-webapp")

function Step($n, $total, $msg) {
    Write-Host "`n[$n/$total] $msg" -ForegroundColor Cyan
}

function Die($msg) {
    Write-Error $msg
    exit 1
}

# 1 — Build local (cible nas)
if ($NoBuild) {
    Write-Host "`n[1/4] Build ignoré (--NoBuild)." -ForegroundColor DarkGray
} else {
    Step 1 4 "Build des images (cible nas)..."
    docker compose build
    if ($LASTEXITCODE -ne 0) { Die "Build échoué." }
}

# 2 — Transfert des images vers le NAS
# Note : cmd /c est requis pour le pipe binaire — PowerShell corrompt les flux binaires
$i = 1
foreach ($img in $Images) {
    Step 2 4 "Transfert $img vers le NAS ($i/$($Images.Count))..."
    cmd /c "docker save ${img}:latest | ssh $NasSsh $NasDocker load"
    if ($LASTEXITCODE -ne 0) { Die "Transfert de $img échoué." }
    $i++
}

# 3 — Mise à jour du compose file dans le répertoire de l'app Container Station
Step 3 4 "Mise à jour du compose file sur le NAS..."
scp docker-compose-nas.yaml "${NasSsh}:${NasAppDir}/docker-compose.yml"
if ($LASTEXITCODE -ne 0) { Die "Copie du compose file échouée. Vérifier que '$NasAppDir' existe sur le NAS." }

# 4 — Redémarrage des conteneurs + nettoyage
Step 4 4 "Redémarrage des conteneurs + nettoyage des images obsolètes..."
ssh $NasSsh @"
cd '$NasAppDir' && \
$NasDocker compose up -d --remove-orphans && \
$NasDocker image prune -f
"@
if ($LASTEXITCODE -ne 0) { Die "Redémarrage échoué." }

Write-Host "`nDéploiement terminé avec succès !" -ForegroundColor Green
