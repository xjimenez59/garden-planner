

API :
La racine de tous les points d'entrée est : <monsite.com>/api/1.0/<clientId>/garden/<gardenId>

Puis :
GET /logs
      Renvoie toutes les actions enregistrées
GET /logs?from=<yyyymmdd>
      Renvoie les actions planifiées ou effectuées à partir le la date transmise,
GET /logs?from=<yyyymmdd>&to=<yyyymmdd>
      Renvoie les actions planifiées ou effectuées entre les deux bornes,

POST /log
      Importe une liste d'actions.
      Les actions référencées par un ID sont mises à jour ;
      Referentiels :
      - si une valeur de referentiel n'a pas d'id correspondant, on cree le referentiel
      - 

PUT  /log/<logID>
      Modifie l'action correspondante
GET  /log/<logID>
      Récupère l'action correspondante
DELETE  /log/<logID>
      Supprime l'action correspondante
