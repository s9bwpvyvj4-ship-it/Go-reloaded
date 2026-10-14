# Go-Reloaded

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Status](https://img.shields.io/badge/status-active-success.svg)

> Un outil en ligne de commande écrit en Go qui complète, édite et autocorrige un texte selon une série de règles définies.

---

## Description

**Go-Reloaded** est un utilitaire qui prend un fichier texte en entrée, applique plusieurs transformations typographiques et de formatage, puis écrit le résultat dans un fichier de sortie.

Il gère notamment :

- La **conversion de nombres** (hexadécimal, binaire → décimal)
- La **modification de casse** (majuscules, minuscules, capitalisation)
- La **correction de la ponctuation** (espaces, apostrophes, groupes de signes)
- La **correction grammaticale** (`a` → `an` devant une voyelle ou un `h`)

Ce projet a été réalisé dans le cadre du cursus **Zone01 / 01 Edu**.

---

## Badges

| Badge | Signification |
|---|---|
| ![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go) | Version de Go requise |
| ![License](https://img.shields.io/badge/license-MIT-blue.svg) | Licence du projet |
| ![Status](https://img.shields.io/badge/status-active-success.svg) | Projet maintenu |

---

## Visuals

Exemple d'utilisation en ligne de commande :

```bash
$ cat sample.txt
it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it

$ go run . sample.txt result.txt

$ cat result.txt
It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age