# Favourites migration

Bus Eta Bot used to store user favourites in GCP's Cloud Datastore, but has
since moved to Postgres. This script migrates existing favourites from Cloud
Datastore into the favourites table in Postgres.

## Usage

```
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python migrate.py gcp_project_id datastore_namespace database_url
```