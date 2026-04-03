#!/bin/sh
set -e

# Symlink DB to persistent volume
ln -sf /app/db/db.sqlite3 /app/db.sqlite3

# Apply custom config overlay
cp /app/openbench-config/config.json /app/Config/config.json
cp /app/openbench-config/Tofiks.json /app/Engines/Tofiks.json

# Run migrations
python manage.py migrate --noinput

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
