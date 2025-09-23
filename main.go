package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// modèle de Contact 
type Contact struct {
	ID    int
	Name  string
	Email string
}

// base de données
var (
	contacts = make(map[int]*Contact) // map[int]*Contact
	nextID   = 1
)

//Constructeur : NewContact 
func NewContact(name, email string) (*Contact, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return nil, fmt.Errorf("nom vide")
	}
	if !isValidEmail(email) {
		return nil, fmt.Errorf("email invalide")
	}

	return &Contact{
		Name:  name,
		Email: email,
	}, nil
}

var emailRe = regexp.MustCompile(`^[\w.+-]+@[\w.-]+\.[A-Za-z]{2,}$`)

func isValidEmail(s string) bool {
	return emailRe.MatchString(s)
}

//Méthodes sur Contact 

// Add 
func (c *Contact) Add() (int, error) {
	if c.ID != 0 {
		return -1, fmt.Errorf("ce contact a déjà un ID (%d), utilisez Update", c.ID)
	}
	c.ID = nextID
	nextID++
	if _, exists := contacts[c.ID]; exists {
		return -1, fmt.Errorf("l'ID %d existe déjà", c.ID)
	}
	contacts[c.ID] = c
	return c.ID, nil
}

// Delete 
func (c *Contact) Delete() bool {
	if c.ID == 0 {
		return false
	}
	if _, ok := contacts[c.ID]; ok {
		delete(contacts, c.ID)
		return true
	}
	return false
}

// Update 
func (c *Contact) Update(newName, newEmail string) error {
	if s := strings.TrimSpace(newName); s != "" {
		c.Name = s
	}
	if s := strings.TrimSpace(newEmail); s != "" {
		if !isValidEmail(s) {
			return fmt.Errorf("email invalide")
		}
		c.Email = s
	}
	contacts[c.ID] = c
	return nil
}

func main() {
	// add contact grace a flag 
	addFlag := flag.Bool("add", false, "Ajouter un contact via des flags")
	nameFlag := flag.String("name", "", "Nom du contact à ajouter")
	emailFlag := flag.String("email", "", "Email du contact à ajouter")
	flag.Parse()

	if *addFlag {
		c, err := NewContact(*nameFlag, *emailFlag)
		if err != nil {
			fmt.Println("Erreur:", err)
			fmt.Println(`Exemple: go run . -add -name "Anna" -email "anna@gmail.com"`)
		} else {
			id, err := c.Add()
			if err != nil {
				fmt.Println("Erreur:", err)
			} else {
				fmt.Printf("Contact '%s' ajouté avec l'ID %d.\n", c.Name, id)
			}
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
		case 1: // ajouter
			name, _ := readLine(reader, "Entrez le nom du contact : ")
			email, _ := readLine(reader, "Entrez l'email du contact : ")
			c, err := NewContact(name, email)
			if err != nil {
				fmt.Println("Erreur:", err)
				fmt.Println()
				continue
			}
			id, err := c.Add()
			if err != nil {
				fmt.Println("Erreur:", err)
			} else {
				fmt.Printf("Contact '%s' ajouté avec l'ID %d.\n\n", c.Name, id)
			}

		case 2: // lister
			listContacts()

		case 3: // supprimer
			idStr, _ := readLine(reader, "Entrez l'ID du contact à supprimer : ")
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("ID invalide.")
				continue
			}
			c := contacts[id]
			if c == nil {
				fmt.Println("Aucun contact avec cet ID.\n")
				continue
			}
			if c.Delete() {
				fmt.Println("Contact supprimé.\n")
			} else {
				fmt.Println("Suppression échouée.\n")
			}

		case 4: // mettre à jour
			idStr, _ := readLine(reader, "Entrez l'ID du contact à mettre à jour : ")
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("ID invalide.")
				continue
			}
			c := contacts[id]
			if c == nil {
				fmt.Println("Aucun contact avec cet ID.\n")
				continue
			}
			name, _ := readLine(reader, "Nouveau nom (laisser vide pour ne pas changer) : ")
			email, _ := readLine(reader, "Nouvel email (laisser vide pour ne pas changer) : ")
			if err := c.Update(name, email); err != nil {
				fmt.Println("Erreur:", err)
			} else {
				fmt.Println("Contact mis à jour.\n")
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

func listContacts() {
	fmt.Println("\n--- Liste des Contacts ---")
	if len(contacts) == 0 {
		fmt.Println("(aucun contact)\n")
		return
	}

	ids := make([]int, 0, len(contacts))
	for id := range contacts {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		c := contacts[id]
		fmt.Printf("ID: %d, Nom: %s, Email: %s\n", c.ID, c.Name, c.Email)
	}
	fmt.Println()
}
