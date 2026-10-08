package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// main lit les arguments, lit le fichier d'entrée,
// applique les transformations et écrit le résultat.
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Erreur lecture fichier:", err)
		os.Exit(1)
	}

	result := ProcessText(string(content))

	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Println("Erreur écriture fichier:", err)
		os.Exit(1)
	}
}

// ProcessText applique toutes les transformations dans le bon ordre :
// nombres -> casse -> ponctuation -> a/an.
func ProcessText(text string) string {
	text = processNumbers(text)
	text = processCases(text)
	text = processPunctuation(text)
	text = processAToAn(text)
	return text
}

// processNumbers convertit "<mot> (hex)" et "<mot> (bin)" en décimal.
func processNumbers(text string) string {
	reHex := regexp.MustCompile(`(\S+)\s*\(hex\)`)
	reBin := regexp.MustCompile(`(\S+)\s*\(bin\)`)

	// (hex) -> décimal (base 16)
	text = reHex.ReplaceAllStringFunc(text, func(match string) string {
		sub := reHex.FindStringSubmatch(match)
		val, err := strconv.ParseInt(sub[1], 16, 64)
		if err != nil {
			return match
		}
		return strconv.FormatInt(val, 10)
	})

	// (bin) -> décimal (base 2)
	text = reBin.ReplaceAllStringFunc(text, func(match string) string {
		sub := reBin.FindStringSubmatch(match)
		val, err := strconv.ParseInt(sub[1], 2, 64)
		if err != nil {
			return match
		}
		return strconv.FormatInt(val, 10)
	})

	return text
}

// processCases gère (up), (low), (cap) et leurs variantes avec nombre.
// On parcourt de droite à gauche pour ne pas décaler les index.
func processCases(text string) string {
	reCase := regexp.MustCompile(`\((?i)(up|low|cap)(?:,\s*(\d+))?\)`)
	matches := reCase.FindAllStringSubmatchIndex(text, -1)

	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		fullStart, fullEnd := m[0], m[1]
		mode := strings.ToLower(text[m[2]:m[3]])

		count := 1
		if m[4] != -1 {
			if n, err := strconv.Atoi(text[m[4]:m[5]]); err == nil {
				count = n
			}
		}

		before := text[:fullStart]
		wordsStart := findNthWordStart(before, count)
		if wordsStart == -1 {
			continue
		}

		prefix := text[:wordsStart]
		wordsPart := text[wordsStart:fullStart]
		suffix := text[fullEnd:]

		transformed := applyCaseToWords(wordsPart, mode)
		text = prefix + transformed + suffix
	}

	return text
}

// findNthWordStart retourne l'index de début des N derniers mots de s.
func findNthWordStart(s string, n int) int {
	count := 0
	i := len(s) - 1

	// Sauter les espaces finaux.
	for i >= 0 && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n') {
		i--
	}

	inWord := true
	for i >= 0 {
		isSpace := s[i] == ' ' || s[i] == '\t' || s[i] == '\n'
		if isSpace && inWord {
			count++
			if count == n {
				return i + 1
			}
			inWord = false
		} else if !isSpace {
			inWord = true
		}
		i--
	}
	return 0
}

// applyCaseToWords applique la transformation sur chaque mot de s.
func applyCaseToWords(s, mode string) string {
	var result strings.Builder
	var word strings.Builder

	flushWord := func() {
		if word.Len() > 0 {
			result.WriteString(transformWord(word.String(), mode))
			word.Reset()
		}
	}

	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			flushWord()
			result.WriteRune(r)
		} else {
			word.WriteRune(r)
		}
	}
	flushWord()
	return result.String()
}

// transformWord applique up/low/cap à un seul mot.
func transformWord(w, mode string) string {
	switch mode {
	case "up":
		return strings.ToUpper(w)
	case "low":
		return strings.ToLower(w)
	case "cap":
		return capitalize(w)
	}
	return w
}

// capitalize met la première lettre en majuscule, le reste en minuscules.
func capitalize(w string) string {
	if len(w) == 0 {
		return w
	}
	return strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
}

// processPunctuation colle les signes aux mots et gère ..., !?, '...'.
func processPunctuation(text string) string {
	// Groupes de ponctuation : pas d'espace interne.
	for _, g := range []string{"...", "!?", "?!", "!!", "??"} {
		text = strings.ReplaceAll(text, " "+g, g)
	}

	// Apostrophes.
	text = processApostrophes(text)

	// Pas d'espace AVANT un signe.
	rePunct := regexp.MustCompile(`\s+([.,!?:;])`)
	text = rePunct.ReplaceAllString(text, "$1")

	// Un espace APRÈS un signe (sauf si suivi d'un autre signe).
	reAfter := regexp.MustCompile(`([.,!?:;])([^\s.,!?:;])`)
	text = reAfter.ReplaceAllString(text, "$1 $2")

	return text
}

// processApostrophes colle les apostrophes au contenu, sans espace interne.
func processApostrophes(text string) string {
	re := regexp.MustCompile(`'\s*([^']*?)\s*'`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		sub := re.FindStringSubmatch(match)
		return "'" + strings.TrimSpace(sub[1]) + "'"
	})
}

// processAToAn transforme "a" en "an" si le mot suivant commence
// par une voyelle ou par "h".
func processAToAn(text string) string {
	words := tokenize(text)
	for i := 0; i < len(words)-1; i++ {
		w := words[i]
		if w == "a" || w == "A" {
			next := words[i+1]
			if len(next) > 0 && isVowelOrH(next[0]) {
				if w == "a" {
					words[i] = "an"
				} else {
					words[i] = "An"
				}
			}
		}
	}
	return strings.Join(words, "")
}

// tokenize découpe en mots et espaces séparés (pour reconstruire à l'identique).
func tokenize(s string) []string {
	var tokens []string
	var current strings.Builder

	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

// isVowelOrH teste si c est une voyelle ou 'h' (insensible à la casse).
func isVowelOrH(c byte) bool {
	c = toLowerByte(c)
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' || c == 'h'
}

// toLowerByte convertit un byte majuscule ASCII en minuscule.
func toLowerByte(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}
