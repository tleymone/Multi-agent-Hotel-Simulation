# Rapport du projet de IA04 : Simulation d'un hôtel

Multi-agent hotel simulation system with API realized with Go, made in 2023.

## Fonctionnement

Ce projet permet de simuler la "vie" d'un hôtel.

Notre hôtel est composé d’un nombre prédéfini de chambres. Chaque chambre est définie par un numéro et par un nombre de personnes qu’elle peut accueillir.

Un client peut demander à réserver une chambre, en indiquant le nombre de personnes pour la réservation, ainsi que les dates. A partir de ces informations, la réservation est faite. Les clients peuvent arriver à partir de 16h et doivent repartir avant 11h le matin.

Il y a plusieurs types d’employés qui travaillent dans l’hôtel. Un réceptionniste est présent 24 heures sur 24 pour accueillir les clients et gérer les réservations. Ainsi, plusieurs réceptionnistes se relaient pour une journée complète.
La journée, une équipe de nettoyage est présente afin de ranger les chambres qui ont été libérées.

Vous pouvez accèder aux différentes informations par le lien `localhost:8080/view` qui imite l'interface de gestion du manager de l'hôtel.

Sur l'interface, on observe à gauche la liste des employés. S'ils sont sur leurs horaires de travail, leur pastille grise devient verte. Dans ce cas, ils peuvent être ou bien en "travail" s'ils sont occupé, "libre" s'ils ont réalisé toutes leurs taches pour l'instant.

En haut s'affiche l'heure de la journée, le nombre de jours depuis le début de la simulation et l'état de la trésorerie.

Le solde initial de l'hotel est de 10.000€ et évolue au fur et à mesure du paiement des chambres et du versement du salaire des employés.

A travers cette simulation, nous allons tenter d’optimiser l'hôtel selon différentes situations, pour que notre petite entreprise fonctionne au mieux. Notre indicateur pour cela est son évolution financière.

## Installation et architecture

Pour installer le projet, clonez le projet depuis le repository https://gitlab.utc.fr/mdelcroi/IA04-hotel.git
Ensuite, lancez la démo en exécutant les commandes

```bash=
cd IA04-hotel/
go run .\bin\demo\launch-demo.go
```

Le front est ensuite accessible à l'adresse : `http://localhost:8080/view`

L'architecture du projet est centrée autour d'un serveur qui contrôle tous les messages d'entrées provenant des différents agents clients (réceptionnistes, agents de nettoyage, clients). Tous les messages passent obligatoirement par le serveur, même si des agents communiquent entre eux, afin que la synchronisation du temps se fasse correctement.

Le serveur a aussi pour rôle d'envoyer des informations importantes pour chaque salarié toutes les 500 millisecondes afin qu'ils sachent ce qu'ils ont à faire.

Le serveur est donc une clé de voûte au projet. C'est lui qui permet de mettre à jour l'état du monde au fur et à mesure que les agents interagissent.

Le Front utilise une librairie JavaScript : **p5.js**.

P5.js permet de dessiner sur une page web.
L'avantage est que la librairie ajoute 2 fonctions principales, une fonction `setup()` et une fonction `draw()`.
A la manière d'un framework de jeu vidéo, la fonction draw est appelée à une fréquence définit par `frameRate()`, pour notre cas à 30Hz.
De plus il permet de dessiner des ronds et du texte beaucoup plus facilement que d'autres librairies.
L'inconvénient est qu'il est beaucoup plus compliqué de faire du responsive puisque tout est dessiné (le texte n'est pas sélectionnable par exemple).

## Analyse et critique

Les emplois du temps des agents employés sont générés aléatoirement. Certaines conditions ont été mises en places :
chaque employé a au minimum 2 jours de congé par semaine
les réceptionnistes ont au maximum un shift par jour, chaque shift faisant 8h.

Le fait que cela soit aléatoire peut donner des emplois du temps avec par exemple plusieurs réceptionnistes travaillant sur certains créneaux, et aucun sur d’autres. Les heures de travail ne sont donc pas optimisées pour chaque employé, alors qu’ils sont payés et considérés de la même manière.

Cependant, en diminuant le nombre de réceptionnistes, moins de réservations peuvent être validées. De la même manière, s’il y a moins d’employés de ménage, ils ne pourront pas nettoyer l’entièreté des chambres. Or, tant qu’une chambre n'est pas nettoyée, elle ne peut pas être réservée. Cela implique donc de grosses pertes d’argent.

D’autre part, pour plus de réalisme, il est possible d’adapter la trésorerie de l'hôtel. Pour l’instant, la situation choisie est telle que :

- Il a une trésorerie de 10.000€ au jour 0.
- Les salaire des employés sont versés à la fin de chaque semaine.
- Les chambres sont payées tous les jours à 18h si elles sont réservées.
- Chaque chambre, occupée ou non, coûte à l'hôtel un prix fixe (charge).

Dans un cas général, tout est facilement adaptable, que cela soit d’un point de vu du nombre d’agents ou de chambres, des prix, des salaires, ou en ajoutant de nouveaux modules.

Notre hôtel peut être considéré comme situé dans un endroit à forte affluence. En effet, il y a toujours de nouveaux clients, il n’y a pas de période creuse ni de concurrence. Pour cette raison, même avec peu d’employés, l'hôtel est toujours attractif et donc très souvent au-delà de 50% de son occupation.

D’un point de vu global sur la situation, il n’y a pas de retour qualité, ni des clients, ni des employés. Pour cette raison, durant nos tests, lorsque nous essayons des cas “classiques”, ils semblent relativement proches de la réalité d’après nos recherches.
C’est lorsque nous faisons des tests dans des cas “extrêmes” que nous voyons des absurdités : dans certains cas, par exemple avec un unique employé de ménage, l'hôtel tourne toujours et gagne de l’argent là où dans la réalité cela n’aurait pas pu tenir.

Enfin, si nous parlons maintenant de l'interface graphique en elle-même, nous pensons qu'elle est assez claire pour comprendre l'état d'avancement des réservations et de l'état des salariés au fil du temps. Elle ne permet cependant pas de voir les réservations futures. Nous aurions donc pu ajouter une partie permettant de visualiser ces dernières.

## Conclusion

Ce projet nous a permis de mettre en oeuvre une simulation faisant communiquer des milliers d'agents entre eux. Malgré les différentes critiques que nous avons remonté, nous sommes heureux du résultat puisque le programme permet de visualiser en temps réel les changements d'état de l'hôtel et qu'il soutient très bien des quantités très importantes de requêtes.
