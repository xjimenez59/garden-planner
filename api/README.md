

API :
La racine de tous les points d'entrée est : <monsite.com>/api/1.0/<clientId>/garden/<gardenId>

Puis :
GET /logs
      Renvoie toutes les actions enregistrées,
      dans l'ordre chronologique inverse (les plus récentes en premier)
GET /logs?from=<yyyymmdd>
      Renvoie les actions planifiées ou effectuées à partir le la date transmise,
GET /logs?from=<yyyymmdd>&to=<yyyymmdd>
      Renvoie les actions planifiées ou effectuées entre les deux bornes,

POST /logs
      Importe une liste d'actions.
      Les actions référencées par un ID sont mises à jour ;
      Referentiels : pour l'instant, on stocke les libellés
       

PUT  /log/<logID>
      Modifie l'action correspondante
GET  /log/<logID>
      Récupère l'action correspondante
DELETE  /log/<logID>
      Supprime l'action correspondante


GET /legumes
   Renvoie la liste de tous les légumes référencés, dans l'ordre alphabétique
