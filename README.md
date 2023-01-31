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
  - [Tests et mode debug  🔧](#tests-et-mode-debug--)
    - [Tests automatisés ](#tests-automatisés-)
    - [Mode debug ](#mode-debug-)

## Introduction 
Ce laboratoire a pour but d'implémenter l'algorithme ondulatoire et l'algorithme sondes et echos en go. Les communications client-serveur sont réalisées avec le protocole UDP. 
La concurrence d'accès aux variables est gérée avec des goroutines et des channels. Cette partie est dédiée à l'algorithme ondulatoire. L'algorithme sondes et echos est implémenté dans une autre branche.

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
* _send_ est envoyé à tout les serveurs du réseau, il permet de compter le nombre de lettre dans le message envoyé.
  * Le message est précisé avec l'argument -word suivi du message à envoyer. Par défaut, le message est "Barack Obama".

Exemple de commande pour demander au serveur 2 qui est le processus élu :
```
go run . -server 2  -command leader
```


## Tests et mode debug <a name="tests"/> 🔧
### Tests automatisés <a name="automated-test"/> 
Pour lancer les tests automatisés, il faut lancer tous les serveurs :
* ``` go run . -id 0 ```
* ``` go run . -id 1 ```
* ``` go run . -id 2 ```
* ``` go run . -id 3 ```  

Puis, dans src/main/test, lancez ```go run .``` pour lancer les tests automatisés. 
Le résultat devrait être le suivant :
![](automated-tests.jpg)
### Mode debug <a name="debug-mode"/>
Les serveurs peuvent être lancés en mode debug, ce qui aura pour effet de ralentir le
traitement des messages de 1 seconde. Pour lancer un serveur en mode debug, ajoutez
l'argument -debug à la commande de lancement du serveur.
Pour tester le cas où plusieurs élections sont demandées simultanément depuis plusieurs
clients, on peut par exemple procéder comme suit :
* Lancer 4 serveurs en mode debug
  * ``` go run . -id 0 -debug ```
  * ``` go run . -id 1 -debug ```
  * ``` go run . -id 2 -debug ```
  * ``` go run . -id 3 -debug ```
* Ajouter des charges sur les serveurs 0 et 2 :
  * ``` go run . -server 0 -command charge 10 ```
  * ``` go run . -server 2 -command charge 5 ```
* Lancer 2 clients qui demandent deux élections en envoyant rapidment les 2 requêtes :
  * ```go run . -server 0  -command elect```
  * ```go run . -server 2  -command elect```    
  
On peut observer les échanges entre les serveurs et constater qu'ils se mettent en 
accord sur le même serveur élu.  
On peut également varier les tests, par exemple en utilisant la commande qui stoppe un serveur
depuis un client avec la commande ```-command stop```, ou augmenter la charge d'un serveur
  avec ```-command charge [amount]``` et observer les échanges entre serveur ainsi que
le résultat du serveur élu. On peut aussi lancer un nouveau serveur pendant l'élection.

Par exemple :
* Lancer 3 serveurs en mode debug
  * ``` go run . -id 0 -debug ```
  * ``` go run . -id 1 -debug ```
  * ``` go run . -id 2 -debug ```
* Ajouter des charges sur les serveurs 0 et 2 :
  * ``` go run . -server 0 -command charge 10 ```
  * ``` go run . -server 2 -command charge 5 ```
* Lancer 1 client qui demande un election
  * ```go run . -server 1  -command elect```
* Lancer un nouveau serveur
  * ``` go run . -id 3 -debug ```
* Redemander une élection
  * ```go run . -server 0  -command elect```


Ou bien :
* Lancer 4 serveurs en mode debug
  * ``` go run . -id 0 -debug ```
  * ``` go run . -id 1 -debug ```
  * ``` go run . -id 2 -debug ```
  * ``` go run . -id 3 -debug ```	
* Ajouter des charges sur les serveurs 0 et 2 :
  * ``` go run . -server 0 -command charge 10 ```
  * ``` go run . -server 2 -command charge 5 ```
* Lancer une élection
  * ```go run . -server 1  -command elect```
* Rapidement tuer le serveur élu
  * ```go run . -server 0  -command stop```