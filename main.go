package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// model de Contact 
type Contact struct {
	ID    int
	Name  string
	Email string
}

// base de donnee
var (
	contacts = make(map[int]Contact) 
	nextID   = 1
)

func main() {
	// add contact grace a flag 
	addFlag := flag.Bool("add", false, "Ajouter un contact via des flags")
	nameFlag := flag.String("name", "", "Nom du contact à ajouter")
	emailFlag := flag.String("email", "", "Email du contact à ajouter")
	flag.Parse()

	if *addFlag {
    if strings.TrimSpace(*nameFlag) == "" || strings.TrimSpace(*emailFlag) == "" {
        fmt.Println("Erreur: utilisez -name et -email avec -add, ex.:")
        fmt.Println(`  go run . -add -name "Anna" -email "anna@gmail.com"`)
    } else {
        id := addContact(strings.TrimSpace(*nameFlag), strings.TrimSpace(*emailFlag))
        fmt.Printf("Contact '%s' ajouté avec l'ID %d.\n", *nameFlag, id)
    }
}

	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu()
		choiceStr, _ := readLine(reader, "Votre choix : ")
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
		if err != nil {
			fmt.Println("Choix invalide, veuillez entrer un nombre.")
			continue
		}

		switch choice {
		case 1: // add
			name, _ := readLine(reader, "Entrez le nom du contact : ")
			email, _ := readLine(reader, "Entrez l'email du contact : ")
			id := addContact(strings.TrimSpace(name), strings.TrimSpace(email))
			fmt.Printf("Contact '%s' ajouté avec l'ID %d.\n\n", strings.TrimSpace(name), id)

		case 2: // liste
			listContacts()

		case 3: // supprimer
			idStr, _ := readLine(reader, "Entrez l'ID du contact à supprimer : ")
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("ID invalide.")
				continue
			}
			if deleteContact(id) {
				fmt.Println("Contact supprimé.\n")
			} else {
				fmt.Println("Aucun contact avec cet ID.\n")
			}

		case 4: // update
			idStr, _ := readLine(reader, "Entrez l'ID du contact à mettre à jour : ")
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("ID invalide.")
				continue
			}
			name, _ := readLine(reader, "Nouveau nom (laisser vide pour ne pas changer) : ")
			email, _ := readLine(reader, "Nouvel email (laisser vide pour ne pas changer) : ")
			if updateContact(id, strings.TrimSpace(name), strings.TrimSpace(email)) {
				fmt.Println("Contact mis à jour.\n")
			} else {
				fmt.Println("Aucun contact avec cet ID.\n")
			}

		case 5: // quitter
			fmt.Println("Au revoir.")
			return

		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// les fonctions

func printMenu() {
	fmt.Println("--- Mini CRM ---")
	fmt.Println("1. Ajouter un contact")
	fmt.Println("2. Lister les contacts")
	fmt.Println("3. Supprimer un contact")
	fmt.Println("4. Mettre à jour un contact")
	fmt.Println("5. Quitter")
}

func readLine(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(text, "\r\n"), nil
}

func addContact(name, email string) int {
	if name == "" || email == "" {
		fmt.Println("Nom et email ne doivent pas être vides.")
		return -1
	}
	id := nextID
	contacts[id] = Contact{
		ID:    id,
		Name:  name,
		Email: email,
	}
	nextID++
	return id
}

func listContacts() {
	if len(contacts) == 0 {
		fmt.Println("\n--- Liste des Contacts ---")
		fmt.Println("(aucun contact)\n")
		return
	}

	ids := make([]int, 0, len(contacts))
	for id := range contacts {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	fmt.Println("\n--- Liste des Contacts ---")
	for _, id := range ids {
		c := contacts[id]
		fmt.Printf("ID: %d, Nom: %s, Email: %s\n", c.ID, c.Name, c.Email)
	}
	fmt.Println()
}

func deleteContact(id int) bool {
	if _, ok := contacts[id]; ok {
		delete(contacts, id)
		return true
	}
	return false
}

func updateContact(id int, newName, newEmail string) bool {
	c, ok := contacts[id]
	if !ok {
		return false
	}
	if newName != "" {
		c.Name = newName
	}
	if newEmail != "" {
		c.Email = newEmail
	}
	contacts[id] = c
	return true
}
