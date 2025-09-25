# Mini-CRM CLI

Mini-CRM est un gestionnaire de contacts en ligne de commande écrit en Go. Il illustre une architecture modulaire reposant sur Cobra pour la CLI, Viper pour la configuration et plusieurs backends de persistance (mémoire, JSON et SQLite via GORM).

## Installation

```bash
go build -o bin/crm .
```

Le binaire est placé dans `bin/crm`. Vous pouvez aussi exécuter les commandes à la volée avec `go run . --config config.yaml …` pendant le développement.

## Configuration

Le fichier `config.yaml` pilote le stockage utilisé :

```yaml
storage:
  type: "json"   # valeurs possibles: memory, json, gorm
  path: "contacts.json"  # requis pour json, optionnel pour gorm (contacts.db par défaut)
```

- `memory` : stockage éphémère pratique pour les tests.
- `json` : persistance dans un fichier JSON lisible.
- `gorm` : persistance dans une base SQLite (`contacts.db` par défaut).

## Utilisation

### Exécution des commandes

```bash
./bin/crm --config config.yaml <commande>
```

Si `config.yaml` est dans le dossier courant, l’option `--config` peut être omise.

Commandes disponibles :

- `add` : ajoute un contact. Peut être utilisé en mode flag (`--name`, `--email`) ou en mode interactif (lancez simplement `./bin/crm add` et suivez les invites).
- `list` : affiche tous les contacts avec un tableau formaté.
- `update` : met à jour un contact. Les flags (`--id`, `--name`, `--email`) sont facultatifs ; sans eux l’outil demande l’ID et les champs à modifier.
- `delete` : supprime un contact ; `--id` est optionnel et une invite interactive est proposée sinon.

Toutes les commandes écrivent leur résultat sur `stdout` et renvoient les erreurs sur `stderr`; le code de sortie est `0` en cas de succès.

### Choisir le backend de stockage

- `memory`: idéal pour des tests rapides, les données disparaissent quand le processus se termine.
- `json`: stocke les contacts dans le fichier indiqué par `path`. Le fichier est créé si besoin.
- `gorm`: initialise (ou migre) une base SQLite via GORM. Changez `path` pour choisir le nom/fichier `.db`.

## Architecture

- `cmd/` : commandes Cobra regroupant la logique CLI.
- `internal/contacts` : logique métier (validation, orchestrations) autour des contacts.
- `internal/model` : entités du domaine.
- `internal/storage` : implémentations des backends de persistance.

## Dépendances

- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [GORM](https://gorm.io) + driver SQLite

Tous les fichiers Go sont formatés via `gofmt`.
