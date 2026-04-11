#!/bin/sh
set -e

# Symlink DB to persistent volume
ln -sf /app/db/db.sqlite3 /app/db.sqlite3

# Run migrations
python manage.py migrate --noinput

# Enable WAL journal mode so concurrent workers don't pile up on
# "database is locked" during clientSubmitResults. journal_mode is a
# persistent, file-level setting — once set it sticks across reopens.
python -c "
import sqlite3
c = sqlite3.connect('/app/db.sqlite3')
print('journal_mode =', c.execute('PRAGMA journal_mode=WAL').fetchone()[0])
c.close()
"

# Create admin user if it doesn't exist
if [ -n "$OPENBENCH_ADMIN_USER" ] && [ -n "$OPENBENCH_ADMIN_PASS" ]; then
  python manage.py shell <<PYEOF
from django.contrib.auth.models import User
from OpenBench.models import Profile

if not User.objects.filter(username='${OPENBENCH_ADMIN_USER}').exists():
    user = User.objects.create_superuser(
        username='${OPENBENCH_ADMIN_USER}',
        password='${OPENBENCH_ADMIN_PASS}',
        email='admin@localhost',
    )
    Profile.objects.create(user=user, enabled=True, approver=True)
    print(f'Created admin user: ${OPENBENCH_ADMIN_USER}')
else:
    print('Admin user already exists, skipping')
PYEOF
fi

exec gunicorn OpenSite.wsgi:application --bind 0.0.0.0:8000 --workers 2
