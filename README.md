# Projet groupie tracker

Ce projet a été entièrement réalisé par [Valentin LAMINE](https://github.com/valentinlamine) dans le cadre d'un devoir scolaire visant à faire découvrir l'utilisation des API en Golang sur un site web. Le site utilise une API externe publique et y effectue différentes requêtes.

## Informations sur le projet

Ce projet est donc conformément aux consignes un site Web en Golang qui gère l'API de iTunes, il permet de rechercher des musiques, des films, des épisodes de série TV et des livres. Il permet aussi de consulter la preview de se contenu, une partie de 30 seconde pour la musique et la bande annonce pour les films et série.

## Utilisation du projet

Le projet comprenant donc un serveur, c'est ce dernier que nous devons lancer, il n'y a cependant aucun hébergement pour le projet, vous devez donc tout d'abord cloner le projet :

```bash
git clone https://github.com/valentinlamine/GroupieTracker_Lamine_Valentin.git
```

Une fois cloner, il est nécessaire de posséder Golang sur son ordinateur, si tel est le cas, il suffit de se mettre dans le dossier src du projet et d'exécuter :

```bash
go run main.go
```

Une fois fait le projet devrait être accessible à l'adresse [localhost:80](localhost:80)

## Spécifications techniques

* Front-end : HTML-CSS-JS
* Back-end : Golang
* API-Utilisé [iTunes API](https://performance-partners.apple.com/search-api)

Compatibilité : Le site a été testé sur Chrome, Brave, Firefox ainsi que sur leur versions mobiles.
