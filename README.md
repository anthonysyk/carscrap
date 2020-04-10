## Car Scrap Test

1) Lancer MongoDB et Redis
```
cd deploy
docker-compose up -d
```

2) Lancer le consumer
```
go run carscrap_consumer/consumer.go
```

3) Lancer le producer
```
go run carscrap_producer/producer.go
```

- Heure de début: 14h00

- Heure de fin: 17h00

## Done

- Initialisation d'un Redis et d'un MongoDB
- Création de la collection `Cars` dans MongoDB
- Définition du modèle de données `Car` dans MongoDB
- Producer qui récupère les pages standalone
- Consumer qui consomme les pages standalone

## Reste à faire

- Extraction des données depuis la page standalone
- Stockage dans MongoDB

## A améliorer:

- Extraire plus de données depuis les pages standalone
- Ajouter des paramètres de configuration pour le producer et consumer
    - Nombre de threads
