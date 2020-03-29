import psycopg2
import sys
from google.cloud import datastore

project = sys.argv[1]
namespace = sys.argv[2]
dsn = sys.argv[3]

conn = psycopg2.connect(dsn=dsn)
cur = conn.cursor()

client = datastore.Client(project=project, namespace=namespace)
query = client.query(kind='Favourites')
results = query.fetch()
count = 0
for entity in results:
    if 'Favourites' in entity:
        user_id = entity.id
        favourites = entity['Favourites']
        print(f'{user_id}:')
        for f in favourites:
            print(f'  - {f}')
            cur.execute(
                '''insert into favourites (user_id, name, query)
values (%s, %s, %s)
on conflict (user_id, name) do update set query = excluded.query''',
                (user_id, f, f))
            count += 1

print(f'{count} favourites migrated')

conn.commit()
cur.close()
conn.close()
