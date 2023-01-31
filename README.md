# SDR_L4
## Table des matières
- [SDR\_L4](#sdr_l4)
    - [Table des matières](#table-des-matières)
    - [Introduction](#introduction)
        - [Auteurs 🧍️‍♂️🧍‍♂](#auteurs-️️)
    - [Guide d'utilisation  📚](#guide-dutilisation--)
        - [Installation des ressources  💾](#installation-des-ressources--)
        - [Lancement serveur](#lancement-serveur)
        - [Lancement d'un client](#lancement-dun-client)
    - [Aspects réseau  📶](#aspects-réseau--)
        - [Format du fichier de configuration ](#format-du-fichier-de-configuration-)
            - [Exemple de fichier de configuration:](#exemple-de-fichier-de-configuration)
    - [Application cliente  👥](#application-cliente--)
    - [Tests 🔧](#tests-)

## Introduction
Ce laboratoire a pour but d'implémenter l'algorithme ondulatoire et l'algorithme sondes et echos en go. Les communications client-serveur sont réalisées avec le protocole UDP.
La concurrence d'accès aux variables est gérée avec des goroutines et des channels. Les deux parties ont été séparées sur deux branches différentes afin de faciliter leur utilisation.
Nous allons préciser dans le guide d'utilisation comment gérer ces deux parties. Les deux branches ont été mises en place de façon similaire, il y a donc beaucoup de point communs et
les commandes pour lancer les serveurs/clients sont les mêmes.

### Auteurs <a name="auteurs"/>🧍️‍♂️🧍‍♂
* Nelson Jeanrenaud
* Vincent Peer

## Guide d'utilisation <a name="utilisation" /> 📚
### Installation des ressources <a name="installation"/> 💾
Commencez par cloner notre repository dans le dossier de votre choix, la commande
git est la suivante :
```
git clone https://github.com/VincentPeer/SDR_labo04.git
```
Changer de branche pour passer à l'algorithme ondulatoire :
```
git checkout partie1
```

Changer de branche pour passer à l'algorithme sondes et échos :
```
git checkout partie2
```

### Lancement serveur
Le lancement d'un serveur requiert l'id du serveur à lancer et le chemin qui mène
au fichier de configuration des serveurs. Le chemin possède une valeur par défaut avec un [fichier
de configuration](#file-config) proposé comme exemple et illustré dans les aspects réseau. L'id du premier serveur est
0 et le dernier a pour id max_serveur -1, un id entré hors de ces bornes entraîne une erreur.
Une fois dans le dossier src/main/server, le format de l'entrée à saisir est le suivant :
>go run . -id [ID] -path [PATH]

Voici un exemple de commande pour lancer un serveur :
```
go run . -id 1 -path MaConfigPerso.json
```

Lancement d'un serveur avec la config par défaut :
```
go run . -id 1
```
### Lancement d'un client
Le lancement d'un client requiert le port avec lequel le client communique, l'id du serveur choisi, une commande pour l'action désirée. Encore une fois les valeurs par défaut :
* port : 8079
* config : ../data/config.json
* server : 1
* command : send


Une fois dans le dossier src/main/client, le format de l'entrée à saisir est le suivant :
>go run . -server [ID] -port [PORT] -path [PATH] -command [COMMAND]

Voici un exemple de commande pour lancer un client depuis l'emplacement src/main/client :
```
go run . -server 2 -port 8082 -path config.json -command send
```
Pour les commandes disponibles, voir la section [Application cliente](#client).
## Aspects réseau <a name="reseau" /> 📶
### Format du fichier de configuration <a name="file-config"/>
La configuration réseau est définie dans un fichier de configuration au format JSON. Ce fichier est passé en paramètre au lancement d'un serveur. Il contient les informations suivantes:
* `servers` : liste des serveurs de l'application. Chaque serveur est identifié par un nom unique et possède une adresse IP et un port d'écoute:
    * `id` : nom du serveur
    * `address` : adresse IP du serveur
    * `port` : port d'écoute du serveur
    * `neighbors` : liste des voisins du serveur
    * `letter` : lettre que le serveur doit compter
* `maxServers` : nombre maximal de serveurs dans le réseau
* `timeout` : le délai maximum à attendre pour une réponse en millisecondes

#### Exemple de fichier de configuration:
```json
{
    "servers" : [
        {
            "id" : "server_0",
            "port" : 8080,
            "address" : "127.0.0.1",
            "neighbors" : ["server_1", "server_2", "server_3"],
            "letter" : "L"
        },
        {
            "id" : "server_1",
            "port" : 8081,
            "address" : "127.0.0.1",
            "neighbors" : ["server_0", "server_2"],
            "letter" : "O"
        },
        {
            "id" : "server_2",
            "port" : 8082,
            "address" : "127.0.0.1",
            "neighbors" : ["server_0", "server_1"],
            "letter" : "V"
        },
        {
            "id" : "server_3",
            "port" : 8083,
            "address" : "127.0.0.1",
            "neighbors" : ["server_0"],
            "letter" : "E"
        }
    ],
    "maxServers" : 4,
    "timeout" : 2000
}

```
## Application cliente <a name="client" /> 👥
Le client propose plusieurs commandes que l'on peut soumettre sur n’importe quel
serveur dont on précise le numéro N en paramètre.
Voici les commandes disponibles à ajouter avec l'argument -command :
* _send_ est envoyé à tout les serveurs du réseau, il permet de compter le nombre de lettre dans le message envoyé et va lancer l'algorithme.
    * Le message est précisé avec l'argument -word suivi du message à envoyer. Par défaut, le message est "BarackObama".
* _result_ est envoyé à un serveur spécifique, il permet de récupérer le résultat du comptage de lettre.

Exemple de commande pour demander au serveur 2 qui est le processus élu :
```
go run . -server 2  -command leader
```

## Tests<a name="tests"/> 🔧

Déplacer les 4 terminaux serveurs dans le dossier src/main/server et le terminal client dans le dossier src/main/client.

Lancer les serveurs avec la commande suivante :
```
go run . -id i
```
où i est l'id du serveur (0, 1, 2 ou 3).

Lancer le client avec la commande suivante :
```
go run . -command send -word "VotreMot"
```

Les serveurs vont se partager le travail de comptage de lettre et le client peut demander le résultat du comptage à n'importe quel serveur avec la commande suivante :
```
go run . -command result
```
Ici le client demande le résultat au serveur 1.
