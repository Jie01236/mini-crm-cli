# Mini CRM en Go

Un petit projet en Go (Golang) qui implémente un **mini système de gestion de contacts** en ligne de commande (CLI).  
Ce projet a été réalisé pour pratiquer les notions de base de Go : `for {}`, `switch`, `map`, "comma ok idiom", `if err != nil`, `strconv`, `os.Stdin`, `bufio`, `flag`.

---

## 🎯 Fonctionnalités

1. **Afficher un menu principal en boucle**  
   - Le programme affiche un menu et attend l'entrée de l'utilisateur jusqu’à ce que l’option *Quitter* soit choisie.  

2. **Ajouter un contact**  
   - Possibilité d’ajouter un contact en saisissant son nom et son email.  

3. **Lister les contacts existants**  
   - Affiche tous les contacts stockés dans la mémoire.  

4. **Supprimer un contact (par ID)**  
   - On peut supprimer un contact en entrant son identifiant unique.  

5. **Mettre à jour un contact (par ID)**  
   - On peut modifier le nom et/ou l’email d’un contact existant.  

6. **Quitter le programme**  
   - Termine l’exécution.  

7. **Ajouter un contact grâce à des flags**  
   - Exemple :  
     ```bash
     go run . -add -name "Anna" -email "anna@gmail.com"
     ```  
   - Cette commande ajoute directement un contact sans passer par le menu, puis lance le programme avec ce contact déjà en mémoire.  

---

## 🛠 Technologies utilisées

- **Langage** : Go (Golang)  
- **Librairies standard** :  
  - `fmt` → formatage et affichage  
  - `bufio` → lecture en ligne de commande (entrée standard)  
  - `os` → gestion de l’entrée standard  
  - `strconv` → conversion string ↔ int  
  - `strings` → manipulation de chaînes  
  - `sort` → tri des identifiants avant affichage  
  - `flag` → gestion des arguments de ligne de commande (flags)  

---

## 🚀 Installation et exécution

1. Cloner le projet ou copier le fichier `main.go` :
   ```bash
   mkdir mini-crm-cli && cd mini-crm-cli
   go mod init mini-crm-cli
   ```

2. Lancer le programme en mode **menu interactif** :
   ```bash
   go run .
   ```

3. Ajouter un contact **avec des flags** :
   ```bash
   go run . -add -name "Ema" -email "ema@gmail.com"
   ```

---

## 📖 Exemple d’exécution

```bash
--- Mini CRM ---
1. Ajouter un contact
2. Lister les contacts
3. Supprimer un contact
4. Mettre à jour un contact
5. Quitter
Votre choix : 1
Entrez le nom du contact : Paul
Entrez l'email du contact : paul@example.com
Contact 'Paul' ajouté avec l'ID 1.

--- Mini CRM ---
1. Ajouter un contact
2. Lister les contacts
3. Supprimer un contact
4. Mettre à jour un contact
5. Quitter
Votre choix : 2

--- Liste des Contacts ---
ID: 1, Nom: Paul, Email: paul@example.com
```

---
