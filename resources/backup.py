
# crontab -e

# chmod +x /var/data_mount/da_team/backups/backup.py

# python3 /var/data_mount/da_team/backups/backup.py

# 0 17 * * 5 python3 /var/data_mount/da_team/backups/backup.py

import os
import datetime
import glob

# Database credentials
db_user = "lemma_rw"
db_password = "Lemm@r0cks!"
db_host = "23.108.100.104"
db_name = "lm_teda_crawler"

# Backup directory
backup_dir = "/var/data_mount/da_team/backups"

# Create a backup filename with timestamp
backup_file = f"{backup_dir}/{db_name}_backup_{datetime.datetime.now().strftime('%Y%m%d_%H%M%S')}.sql"

# Create the backup command
backup_command = f"mysqldump -u {db_user} -p'{db_password}' -h {db_host} {db_name} > {backup_file}"

# Execute the backup command
os.system(backup_command)

# Retain only the last 4 backups
backups = sorted(glob.glob(f"{backup_dir}/{db_name}_backup_*.sql"), reverse=True)

# Remove older backups
if len(backups) > 4:
    for old_backup in backups[4:]:
        os.remove(old_backup)
        print(f"Deleted old backup: {old_backup}")

print(f"Backup completed: {backup_file}")




